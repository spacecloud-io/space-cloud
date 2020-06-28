package operations

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// DockerStop stops the services which have been started
func DockerStop() error {

	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	argsServices := filters.Arg("label", "app=service")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsServices), All: true})
	if err != nil {
		return utils.LogError("Unable to list space-cloud services containers", err)
	}

	for _, container := range containers {
		// First stop the container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to stop container (%s)", container.ID), err)
		}
	}

	argsSC := filters.Arg("label", "app=space-cloud")
	scContainers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsSC), All: true})
	if err != nil {
		return utils.LogError("Unable to list space-cloud core containers", err)
	}

	for _, container := range scContainers {
		// First stop the container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to stop container (%s)", container.ID), err)
		}
	}

	argsAddOns := filters.Arg("label", "app=addon")
	addOnContainers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsAddOns), All: true})
	if err != nil {
		return utils.LogError("Unable to list space-cloud core containers", err)
	}

	for _, container := range addOnContainers {
		// First stop the container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to stop container (%s)", container.ID), err)
		}
	}
	return nil
}
