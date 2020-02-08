package cmd

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cli/model"
)

func getSpaceCloudHostsFilePath() string {
	return fmt.Sprintf("%s/hosts", getSpaceCloudDirectory())
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
func CodeSetup(id, username, key, secret string, dev bool) error {
	// todo old keys always remain in accounts.yaml file
	const ContainerGateway string = "space--cloud--gateway"
	const ContainerRunner string = "space--cloud--runner"

	_ = createDirIfNotExist(getSpaceCloudDirectory())
	_ = createDirIfNotExist(getSecretsDir())
	_ = createDirIfNotExist(getTempSecretsDir())

	logrus.Infoln("Setting up space cloud on docker")

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
		ServerUrl: "http://localhost:4122",
	}

	if err := checkCred(&selectedAccount); err != nil {
		logrus.Errorf("error in setup unable to check credentials - %v", err)
		return err
	}

	devMode := "false"
	if dev {
		devMode = "true" // todo: even the flag set true in dev of container sc didn't start in prod mode
	}

	containersToCreate := []struct {
		dnsName        string
		containerImage string
		containerName  string
		envs           []string
		mount          []mount.Mount
		exposedPorts   nat.PortSet
		portMapping    nat.PortMap
	}{
		{
			containerImage: "spaceuptech/gateway",
			containerName:  ContainerGateway,
			dnsName:        "gateway.space-cloud.svc.cluster.local",
			envs: []string{
				"ARTIFACT_ADDR=store.space-cloud.svc.cluster.local:4122",
				"RUNNER_ADDR=runner.space-cloud.svc.cluster.local:4050",
				"ADMIN_USER=" + username,
				"ADMIN_PASS=" + key,
				"ADMIN_SECRET=" + secret,
				"DEV=" + devMode,
			},
			exposedPorts: nat.PortSet{
				"4122": struct{}{},
			},
			portMapping: nat.PortMap{
				"4122": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "4122"}},
			},
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudHostsFilePath(),
					Target: "/etc/hosts",
				},
			},
		},
		{
			// runner
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
			},
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind, // TODO CHECK THIS
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
			},
		},
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error cli setup unable to initialize docker client got error message - %v", err)
		return err
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("error cli setup unable to load host file with suitable default got error message - %v", err)
		return err
	}
	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
	hosts.WriteFilePath = getSpaceCloudHostsFilePath()
	if err := hosts.SaveAs(getSpaceCloudHostsFilePath()); err != nil {
		logrus.Errorf("error cli setup unable to save as host file to specified path (%s) got error message - %v", getSpaceCloudHostsFilePath(), err)
		return err
	}

	for _, c := range containersToCreate {
		logrus.Infof("Starting container %s...", c.containerName)
		// check if image already exists
		if err := pullImageIfNotExist(ctx, cli, c.containerImage); err == nil {
			logrus.Infof("Image %s already exists", c.containerImage)
		}

		// check if container is already running
		args := filters.Arg("name", c.containerName)
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
		if err != nil {
			logrus.Errorf("error deleting service in docker unable to list containers got error message - %v", err)
			return err
		}
		if len(containers) != 0 {
			logrus.Errorf("error in cli setup container already running with name %s", c.containerName)
			return fmt.Errorf("container already running with name %s", c.containerName)
		}

		// create container with specified defaults
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        c.containerImage,
			ExposedPorts: c.exposedPorts,
			Env:          c.envs,
		}, &container.HostConfig{
			Mounts:       c.mount,
			PortBindings: c.portMapping,
		}, nil, c.containerName)
		if err != nil {
			logrus.Errorf("error cli setup unable to create container %s got error message  - %v", c.containerName, err)
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			logrus.Errorf("error cli setup unable to start container %s got error message - %v", c.containerName, err)
			return err
		}

		// get the ip address assigned to container
		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			logrus.Errorf("error cli setup unable to inspect container %s got error message - %v", c.containerName, err)
		}
		hosts.AddHost(data.NetworkSettings.IPAddress, c.dnsName)
	}

	if err := hosts.Save(); err != nil {
		logrus.Errorf("error cli setup unable to save host file got error message - %v", err)
		return err
	}
	logrus.Infof("Space Cloud (id: \"%s\") has been successfully setup! :D", selectedAccount.ID)
	logrus.Infof("You can visit mission control at %s/mission-control", selectedAccount.ServerUrl)
	logrus.Infof("Your login credentials: [username: \"%s\"; key: \"%s\"]", selectedAccount.UserName, selectedAccount.Key)
	return nil
}

func pullImageIfNotExist(ctx context.Context, dockerClient *client.Client, image string) error {
	_, _, err := dockerClient.ImageInspectWithRaw(ctx, image)
	if err != nil {
		// pull image from public repository
		out, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			logrus.Errorf("error cli setup unable to pull public image with id (%s) - %s", image, err.Error())
			return err
		}
		io.Copy(os.Stdout, out)
		logrus.Infof("Image %s already exists", image)
	}
	return nil
}
