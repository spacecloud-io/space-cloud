package addons

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/txn2/txeh"
)

func addDatabase(dbtype, username, password, alias, version string) error {
	ctx := context.Background()
	dockerImage := strings.Join([]string{dbtype, version}, ":")

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", "add", "database", err)
	}

	// Check if a database container already exist
	filterArgs := filters.Arg("label", "app=space-cloud")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if database already exists", "add", "database", err)
	}
	if len(containers) == 0 {
		utils.LogInfo("No space-cloud instance found. Run 'space-cli setup' first", "add", "database")
		return nil
	}

	// Pull image if it doesn't already exist
	if err := utils.PullImageIfNotExist(ctx, docker, dockerImage); err != nil {
		return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", dockerImage), "add", "database", err)
	}

	// Create the database
	containerRes, err := docker.ContainerCreate(ctx, &container.Config{
		Labels: map[string]string{"app": "addon", "service": dbtype, "name": alias},
		Image:  dockerImage,
	}, &container.HostConfig{
		NetworkMode: "space-cloud",
	}, nil, strings.Join([]string{"space-cloud--addon", alias}, "--"))
	if err != nil {
		return utils.LogError("Unable to create local docker database", "add", "database", err)
	}

	// Start the database
	if err := docker.ContainerStart(ctx, containerRes.ID, types.ContainerStartOptions{}); err != nil {
		return utils.LogError("Unable to start local docker database", "add", "database", err)
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to open hosts file", "add", "database", err)
	}

	for _, container := range containers {
		// First start the container
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%s)", container.ID), "add", "database", err)
		}

		// Get the container's info
		info, err := docker.ContainerInspect(ctx, container.ID)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inspect container (%s)", container.ID), "add", "database", err)
		}

		hostName := utils.GetServiceDomain(info.Config.Labels["service"], info.Config.Labels["name"])

		// Remove the domain from the hosts file
		hosts.RemoveHost(hostName)

		// Add it back with the new ip address
		hosts.AddHost(info.NetworkSettings.Networks["space-cloud"].IPAddress, hostName)
	}

	// Save the hosts file
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", "add", "database", err)
	}

	return nil
}

func removeDatabase(alias string) error {
	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", "remove", "database", err)
	}

	// Check if a database container already exist
	filterArgs := filters.Arg("label", "app=space-cloud")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if database already exists", "remove", "database", err)
	}
	if len(containers) == 0 {
		utils.LogInfo("No space-cloud instance found. Run 'space-cli setup' first", "remove", "database")
		return nil
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to open hosts file", "remove", "database", err)
	}

	for _, container := range containers {
		// First start the container
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%s)", container.ID), "remove", "database", err)
		}

		// Get the container's info
		info, err := docker.ContainerInspect(ctx, container.ID)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inspect container (%s)", container.ID), "remove", "database", err)
		}

		hostName := utils.GetServiceDomain(info.Config.Labels["service"], alias)

		// Remove the domain from the hosts file
		hosts.RemoveHost(hostName)
	}

	// Save the hosts file
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", "add", "database", err)
	}

	return nil
}
