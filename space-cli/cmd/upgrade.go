package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	scClient "github.com/spaceuptech/space-api-go"
	spaceApiTypes "github.com/spaceuptech/space-api-go/types"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/txn2/txeh"
)

type resultData struct {
	Docs []*doc `mapstructure:"result"`
}

type doc struct {
	ID                string `mapstructure:"_id" json:"id"`
	VersionNo         string `mapstructure:"version_no" json:"versionNo"`
	CompatibleVersion string `mapstructure:"compatible_version" json:"compatibleVersion"`
}

// Upgrade upgrades the environment which has been setup
func Upgrade() error {
	const ContainerGateway string = "space-cloud-gateway"
	const ContainerRunner string = "space-cloud-runner"

	// getting current version
	result := make(map[string]interface{})
	if err := utils.Get(http.MethodGet, "/v1/config/env", map[string]string{}, &result); err != nil {
		return err
	}
	currentVersion := result["version"].(string)

	// getting latest version
	latestVersion, err := getLatestVersion(currentVersion)
	if err != nil {
		return err
	}

	if currentVersion == latestVersion {
		return fmt.Errorf("current verion (%s) is up to date", currentVersion)
	}

	// creating docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", "operations", "upgrade", err)
	}

	// get all containers containing < space-cloud > in their name
	args := filters.Arg("label", "app=space-cloud")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		utils.LogInfo("space cloud is not setup. Consider running space-cli setup first.", "operation", "upgrade")
		return err
	}

	// parameters for gateway
	var gatewayMounts []mount.Mount
	var gatewayPorts nat.PortMap
	var gatewayEnvs []string

	// parameters for runner
	var runnerEnvs []string
	var runnerMounts []mount.Mount

	// Remove all container
	for _, containerInfo := range containers {
		containerInspect, err := cli.ContainerInspect(ctx, containerInfo.ID)
		if err != nil {
			return utils.LogError("error getting service in docker unable to inspect container", "operations", "upgrade", err)
		}

		switch containerInspect.Config.Labels["service"] {
		case "gateway":
			gatewayEnvs = containerInspect.Config.Env
			gatewayMounts = containerInspect.HostConfig.Mounts
			gatewayPorts = containerInspect.HostConfig.PortBindings
			if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				return utils.LogError(fmt.Sprintf("Unable to remove container - %s", containerInfo.ID), "operations", "upgrade", err)
			}

		case "runner":
			runnerEnvs = containerInspect.Config.Env
			runnerMounts = containerInspect.HostConfig.Mounts
			if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				return utils.LogError(fmt.Sprintf("Unable to remove container - %s", containerInfo.ID), "operations", "upgrade", err)
			}
		}
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
			containerImage: fmt.Sprintf("%s:%s", "spaceuptech/gateway", latestVersion),
			containerName:  ContainerGateway,
			dnsName:        "gateway.space-cloud.svc.cluster.local",
			envs:           gatewayEnvs,
			exposedPorts: nat.PortSet{
				"4122": struct{}{},
				"4126": struct{}{},
			},
			portMapping: gatewayPorts,
			mount:       gatewayMounts,
		},

		{
			// runner
			containerImage: fmt.Sprintf("%s:v%s", "spaceuptech/runner", latestVersion),
			containerName:  ContainerRunner,
			dnsName:        "runner.space-cloud.svc.cluster.local",
			envs:           runnerEnvs,
			mount:          runnerMounts,
		},
	}

	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", "operations", "upgrade", err)
	}

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to load host file with suitable default", "operations", "upgrade", err)
	}

	for _, c := range containersToCreate {
		utils.LogInfo(fmt.Sprintf("Starting container %s...", c.containerName), "operations", "upgrade")
		// check if image already exists
		if err := utils.PullImageIfNotExist(ctx, cli, c.containerImage); err != nil {
			return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", c.containerImage), "operations", "upgrade", err)
		}

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        c.containerImage,
			ExposedPorts: c.exposedPorts,
			Env:          c.envs,
		}, &container.HostConfig{
			Mounts:       c.mount,
			PortBindings: c.portMapping,
			NetworkMode:  "space-cloud",
		}, nil, c.containerName)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to create container (%v)", c.containerName), "operations", "upgrade", err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%v)", c.containerName), "operations", "upgrade", err)
		}

		// get the ip address assigned to container
		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inpect container (%v)", c.containerName), "operations", "upgrade", err)
		}
		// Remove the domain from the hosts file
		hosts.RemoveHost(c.dnsName)
		// Add it back with the new ip address
		ip := data.NetworkSettings.Networks["space-cloud"].IPAddress

		hosts.AddHost(ip, c.dnsName)
	}

	if err := hosts.Save(); err != nil {
		return utils.LogError("Unable to save host file", "operations", "upgrade", err)
	}

	return nil
}

func getLatestVersion(version string) (string, error) {
	db := scClient.New("space_cloud", "localhost:4122", false).DB("db")
	if db == nil {
		return "", fmt.Errorf("cannot connect to db")
	}
	ctx := context.Background()

	var result *spaceApiTypes.Response
	var err error
	if version == "" {
		result, err = db.GetOne("space_cloud_version").Sort("-version_no").Limit(1).Apply(ctx)
		if err != nil {
			return "", err
		}
	} else {
		result, _ = db.Get("space_cloud_version").Where(spaceApiTypes.Cond("compatible_version", "==", version)).Apply(ctx)
		if err != nil {
			return "", err
		}
	}

	r := new(resultData)
	if err := result.Unmarshal(&r); err != nil {
		return "", err
	}
	newVersion := version
	for _, val := range r.Docs {
		if val.VersionNo > newVersion {
			newVersion = val.VersionNo
			fmt.Println("new version: ", newVersion)
		}
	}
	return newVersion, nil
}
