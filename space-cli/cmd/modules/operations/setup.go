package operations

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

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

// Setup initializes development environment
func Setup(username, key, config, version, secret, clusterName string, dev bool, portHTTP, portHTTPS int64, volumes, environmentVariables []string) error {
	// TODO: old keys always remain in accounts.yaml file

	_ = utils.CreateDirIfNotExist(utils.GetSpaceCloudDirectory())
	_ = utils.CreateDirIfNotExist(utils.GetSecretsDir(clusterName))
	_ = utils.CreateDirIfNotExist(utils.GetTempSecretsDir(clusterName))

	_ = utils.CreateFileIfNotExist(utils.GetSpaceCloudRoutingConfigPath(clusterName), "{}")
	_ = utils.CreateConfigFile(utils.GetSpaceCloudConfigFilePath(clusterName))

	utils.LogInfo("Setting up Space Cloud on docker.")

	if username == "" {
		username = "local-admin"
	}

	if key == "" {
		key = generateRandomString(32)
	}
	if config == "" {
		config = utils.GetSpaceCloudConfigFilePath(clusterName)
	}
	if !strings.Contains(config, ".yaml") {
		return fmt.Errorf("full path not provided for config file")
	}
	if version == "" {
		utils.LogInfo("Fetching latest Space Cloud Version")

		var err error
		version, err = utils.GetLatestVersion("")
		if err != nil {
			_ = utils.LogError("Unable to fetch the latest Space Cloud version. Sticking to tag latest", err)
			version = "latest"
		}
	}

	if secret == "" {
		secret = generateRandomString(24)
	}

	portHTTPValue := strconv.FormatInt(portHTTP, 10)
	portHTTPSValue := strconv.FormatInt(portHTTPS, 10)

	// set a docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client ", err)
	}

	// check if network already exists
	args := filters.Arg("name", utils.GetNetworkName(clusterName))
	nws, err := cli.NetworkList(ctx, types.NetworkListOptions{Filters: filters.NewArgs(args)})
	if err != nil {
		return utils.LogError("Unable to list networks", err)
	}
	for _, nw := range nws {
		if nw.Name == utils.GetNetworkName(clusterName) {
			return utils.LogError(fmt.Sprintf("Network (%s) already exists, try using different cluster", utils.GetNetworkName(clusterName)), errors.New(""))
		}
	}

	portHTTPValue, err = utils.CheckPortAvailability(portHTTPValue, "HTTP")
	if err != nil {
		return err
	}
	portHTTPSValue, err = utils.CheckPortAvailability(portHTTPSValue, "HTTPS")
	if err != nil {
		return err
	}

	clusterID := fmt.Sprintf("%s--%s", clusterName, randomdata.SillyName())

	selectedAccount := model.Account{
		ID:        clusterID,
		UserName:  username,
		Key:       key,
		ServerURL: "http://localhost:" + portHTTPValue,
	}

	if err := utils.StoreCredentials(&selectedAccount); err != nil {
		return utils.LogError("Unable to store credentials", err)
	}

	devMode := "false"
	if dev {
		devMode = "true" // todo: even the flag set true in dev of container sc didn't start in prod mode
	}

	envs := []string{
		"RUNNER_ADDR=runner.space-cloud.svc.cluster.local:4050",
		"ADMIN_USER=" + username,
		"ADMIN_PASS=" + key,
		"ADMIN_SECRET=" + secret,
		"DEV=" + devMode,
		"GOOGLE_APPLICATION_CREDENTIALS=/root/.gcp/credentials.json",
		"CLUSTER_ID=" + clusterID,
		"PORT=" + portHTTPValue,
	}

	envs = append(envs, environmentVariables...)

	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: utils.GetMountHostsFilePath(clusterName),
			Target: "/etc/hosts",
		},
		{
			Type:   mount.TypeBind,
			Source: utils.GetMountConfigFilePath(clusterName),
			Target: "/app/config.yaml",
		},
	}

	for _, volume := range volumes {
		temp := strings.Split(volume, ":")
		if len(temp) != 2 {
			return utils.LogError(fmt.Sprintf("Error in volume flag (%s) - incorrect format", volume), errors.New(""))
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
			containerImage: fmt.Sprintf("%s:%s", "spaceuptech/gateway", version),
			containerName:  utils.GetScContainers(clusterName, "gateway"),
			dnsName:        "gateway.space-cloud.svc.cluster.local",
			envs:           envs,
			exposedPorts: nat.PortSet{
				nat.Port(portHTTPValue):  struct{}{},
				nat.Port(portHTTPSValue): struct{}{},
			},
			portMapping: nat.PortMap{
				nat.Port(portHTTPValue):  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portHTTPValue}},
				nat.Port(portHTTPSValue): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portHTTPSValue}},
			},
			mount: mounts,
		},

		{
			// runner
			name:           "runner",
			containerImage: fmt.Sprintf("%s:%s", "spaceuptech/runner", version),
			containerName:  utils.GetScContainers(clusterName, "runner"),
			dnsName:        "runner.space-cloud.svc.cluster.local",
			envs: []string{
				"DEV=" + devMode,
				"ARTIFACT_ADDR=store.space-cloud.svc.cluster.local:" + portHTTPValue, // TODO Change the default value in runner it starts with http
				"DRIVER=docker",
				"JWT_SECRET=" + secret,
				"JWT_PROXY_SECRET=" + generateRandomString(24),
				"SECRETS_PATH=/secrets",
				"HOME_SECRETS_PATH=" + utils.GetMountTempSecretsDir(clusterName),
				"HOSTS_FILE_PATH=" + utils.GetMountHostsFilePath(clusterName),
				"ROUTING_FILE_PATH=" + "/routing-config.json",
				"CLUSTER_ID=" + clusterID,
				"PORT=4050",
			},
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: utils.GetMountSecretsDir(clusterName),
					Target: "/secrets",
				},
				{
					Type:   mount.TypeBind,
					Source: utils.GetMountHostsFilePath(clusterName),
					Target: "/etc/hosts",
				},
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
				{
					Type:   mount.TypeBind,
					Source: utils.GetMountRoutingConfigPath(clusterName),
					Target: "/routing-config.json",
				},
			},
		},
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return utils.LogError("Unable to load host file with suitable default", err)
	}
	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
	hosts.WriteFilePath = utils.GetSpaceCloudHostsFilePath(clusterName)
	if err := hosts.SaveAs(utils.GetSpaceCloudHostsFilePath(clusterName)); err != nil {
		return utils.LogError(fmt.Sprintf("Unable to save as host file to specified path (%s)", utils.GetSpaceCloudHostsFilePath(clusterName)), errors.New(""))
	}

	// First we create a network for space cloud
	if _, err := cli.NetworkCreate(ctx, utils.GetNetworkName(clusterName), types.NetworkCreate{Driver: "bridge"}); err != nil {
		return utils.LogError("Unable to create a network named space-cloud", err)
	}

	for _, c := range containersToCreate {
		utils.LogInfo(fmt.Sprintf("Starting container %s...", c.containerName))
		// check if image already exists
		if err := utils.PullImageIfNotExist(ctx, cli, c.containerImage); err != nil {
			return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", c.containerImage), errors.New(""))
		}

		// check if container is already running
		args := filters.Arg("name", c.containerName)
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
		if err != nil {
			return utils.LogError("error deleting service in docker unable to list containers", err)
		}
		if len(containers) != 0 {
			return utils.LogError(fmt.Sprintf("Container (%s) already exists", c.containerName), errors.New(""))
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
			NetworkMode:  container.NetworkMode(utils.GetNetworkName(clusterName)),
		}, nil, c.containerName)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to create container (%v)", c.containerName), err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%v)", c.containerName), err)
		}

		// get the ip address assigned to container
		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inspect container (%v)", c.containerName), err)
		}

		ip := data.NetworkSettings.Networks[utils.GetNetworkName(clusterName)].IPAddress
		utils.LogDebug(fmt.Sprintf("Adding entry (%s - %s) to hosts file", c.dnsName, ip), nil)
		hosts.AddHost(ip, c.dnsName)
	}

	if err := hosts.SaveAs(utils.GetSpaceCloudHostsFilePath(clusterName)); err != nil {
		return utils.LogError("Unable to save host file - %s", err)
	}

	fmt.Println()
	utils.LogInfo(fmt.Sprintf("Space Cloud (cluster id: \"%s\") has been successfully setup! 👍", selectedAccount.ID))
	utils.LogInfo(fmt.Sprintf("You can visit mission control at %s/mission-control 💻", selectedAccount.ServerURL))
	utils.LogInfo(fmt.Sprintf("Your login credentials: [username: \"%s\"; key: \"%s\"] 🤫", selectedAccount.UserName, selectedAccount.Key))

	return nil
}
