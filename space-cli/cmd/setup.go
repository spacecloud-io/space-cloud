package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

func getSpaceCloudHostsFilePath() string {
	return fmt.Sprintf("%s/hosts", getSpaceCloudDirectory())
}

func getSpaceCloudRoutingConfigPath() string {
	return fmt.Sprintf("%s/routing-config.json", getSpaceCloudDirectory())
}

func getSecretsDir() string {
	return fmt.Sprintf("%s/secrets", getSpaceCloudDirectory())
}

func getTempSecretsDir() string {
	return fmt.Sprintf("%s/secrets/temp-secrets", getSpaceCloudDirectory())
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "ExcbsVQs"
}

// CodeSetup initializes development environment
func CodeSetup(id, username, key, secret string, dev bool, portHTTP, portHTTPS int64, volumes, environmentVariables []string) error {
	// TODO: old keys always remain in accounts.yaml file
	const ContainerGateway string = "space-cloud-gateway"
	const ContainerRunner string = "space-cloud-runner"

	_ = createDirIfNotExist(getSpaceCloudDirectory())
	_ = createDirIfNotExist(getSecretsDir())
	_ = createDirIfNotExist(getTempSecretsDir())

	_ = createFileIfNotExist(getSpaceCloudRoutingConfigPath(), "{}")

	logrus.Infoln("Setting up Space Cloud on docker on your command...")

	if username == "" {
		username = "local-admin"
	}
	if id == "" {
		id = username
	}
	if key == "" {
		key = generateRandomString(12)
	}
	if secret == "" {
		secret = generateRandomString(24)
	}

	selectedAccount := model.Account{
		ID:        id,
		UserName:  username,
		Key:       key,
		ServerURL: "http://localhost:4122",
	}

	if err := checkCred(&selectedAccount); err != nil {
		logrus.Errorf("error in setup unable to check credentials - %v", err)
		return err
	}

	devMode := "false"
	if dev {
		devMode = "true" // todo: even the flag set true in dev of container sc didn't start in prod mode
	}

	portHTTPValue := strconv.FormatInt(portHTTP, 10)
	portHTTPSValue := strconv.FormatInt(portHTTPS, 10)

	envs := []string{
		"ARTIFACT_ADDR=store.space-cloud.svc.cluster.local:4122",
		"RUNNER_ADDR=runner.space-cloud.svc.cluster.local:4050",
		"ADMIN_USER=" + username,
		"ADMIN_PASS=" + key,
		"ADMIN_SECRET=" + secret,
		"DEV=" + devMode,
	}

	envs = append(envs, environmentVariables...)

	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: getSpaceCloudHostsFilePath(),
			Target: "/etc/hosts",
		},
	}

	for _, volume := range volumes {
		temp := strings.Split(volume, ":")
		if len(temp) != 2 {
			logrus.Errorf("Error in volume flag (%s) - incorrect format", volume)
			return errors.New("incorrect format for volume flag")
		}

		mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: temp[0], Target: temp[1]})
	}

	containersToCreate := []struct {
		dnsName        string
		containerImage string
		containerName  string
		name           string
		envs           []string
		mount          []mount.Mount
		exposedPorts   nat.PortSet
		portMapping    nat.PortMap
	}{
		{
			name:           "gateway",
			containerImage: "spaceuptech/gateway",
			containerName:  ContainerGateway,
			dnsName:        "gateway.space-cloud.svc.cluster.local",
			envs:           envs,
			exposedPorts: nat.PortSet{
				"4122": struct{}{},
				"4126": struct{}{},
			},
			portMapping: nat.PortMap{
				"4122": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portHTTPValue}},
				"4126": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portHTTPSValue}},
			},
			mount: mounts,
		},

		{
			// runner
			name:           "runner",
			containerImage: "spaceuptech/runner",
			containerName:  ContainerRunner,
			dnsName:        "runner.space-cloud.svc.cluster.local",
			envs: []string{
				"DEV=" + devMode,
				"ARTIFACT_ADDR=store.space-cloud.svc.cluster.local:4122", // TODO Change the default value in runner it starts with http
				"DRIVER=docker",
				"JWT_SECRET=" + secret,
				"JWT_PROXY_SECRET=" + generateRandomString(24),
				"SECRETS_PATH=/secrets",
				"HOME_SECRETS_PATH=" + getTempSecretsDir(),
				"HOSTS_FILE_PATH=" + getSpaceCloudHostsFilePath(),
				"ROUTING_FILE_PATH=" + "/routing-config.json",
			},
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getSecretsDir(),
					Target: "/secrets",
				},
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudHostsFilePath(),
					Target: "/etc/hosts",
				},
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudRoutingConfigPath(),
					Target: "/routing-config.json",
				},
			},
		},
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("Unable to initialize docker client - %s", err)
		return err
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("Unable to load host file with suitable default - %s", err)
		return err
	}
	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
	hosts.WriteFilePath = getSpaceCloudHostsFilePath()
	if err := hosts.SaveAs(getSpaceCloudHostsFilePath()); err != nil {
		logrus.Errorf("Unable to save as host file to specified path (%s) - %s", getSpaceCloudHostsFilePath(), err)
		return err
	}

	// First we create a network for space cloud
	if _, err := cli.NetworkCreate(ctx, "space-cloud", types.NetworkCreate{Driver: "bridge"}); err != nil {
		return utils.LogError("Unable to create a network named space-cloud", "operations", "setup", err)
	}

	for _, c := range containersToCreate {
		logrus.Infof("Starting container %s...", c.containerName)
		// check if image already exists
		if err := pullImageIfNotExist(ctx, cli, c.containerImage); err != nil {
			logrus.Errorf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", c.containerImage)
			return err
		}

		// check if container is already running
		args := filters.Arg("name", c.containerName)
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
		if err != nil {
			logrus.Errorf("error deleting service in docker unable to list containers - %s", err)
			return err
		}
		if len(containers) != 0 {
			logrus.Errorf("Container (%s) already exists", c.containerName)
			return fmt.Errorf("container (%s) already exists", c.containerName)
		}

		// create container with specified defaults
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Labels:       map[string]string{"app": "space-cloud", "service": c.name},
			Image:        c.containerImage,
			ExposedPorts: c.exposedPorts,
			Env:          c.envs,
		}, &container.HostConfig{
			Mounts:       c.mount,
			PortBindings: c.portMapping,
			NetworkMode:  "space-cloud",
		}, nil, c.containerName)
		if err != nil {
			logrus.Errorf("Unable to create container (%s) - %s", c.containerName, err)
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			logrus.Errorf("Unable to start container (%s) - %s", c.containerName, err.Error())
			return err
		}

		// get the ip address assigned to container
		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			logrus.Errorf("Unable to inspect container (%s) - %s", c.containerName, err)
		}

		ip := data.NetworkSettings.Networks["space-cloud"].IPAddress
		utils.LogDebug(fmt.Sprintf("Adding entry (%s - %s) to hosts file", c.dnsName, ip), "operations", "setup", nil)
		hosts.AddHost(ip, c.dnsName)
	}

	if err := hosts.Save(); err != nil {
		logrus.Errorf("Unable to save host file - %s", err.Error())
		return err
	}

	fmt.Println()
	logrus.Infof("Space Cloud (id: \"%s\") has been successfully setup! 👍", selectedAccount.ID)
	logrus.Infof("You can visit mission control at %s/mission-control 💻", selectedAccount.ServerURL)
	logrus.Infof("Your login credentials: [username: \"%s\"; key: \"%s\"] 🤫", selectedAccount.UserName, selectedAccount.Key)
	return nil
}

func pullImageIfNotExist(ctx context.Context, dockerClient *client.Client, image string) error {
	_, _, err := dockerClient.ImageInspectWithRaw(ctx, image)
	if err != nil {
		// pull image from public repository
		logrus.Infof("Image %s does not exist. Need to pull it from Docker Hub. This may take some time.", image)
		out, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			logrus.Errorf("Unable to pull public image with id (%s) - %s", image, err.Error())
			return err
		}
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		s.Suffix = "    Downloading image..."
		_ = s.Color("green")
		s.Start()
		time.Sleep(4 * time.Second) // Run for some time to simulate work// Start the spinner
		_, _ = io.Copy(ioutil.Discard, out)
		s.Stop()
	}
	logrus.Infof("Image %s already exists. No need to pull it again", image)
	return nil
}
