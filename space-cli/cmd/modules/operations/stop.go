package operations

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// DockerStop stops the services which have been started
func DockerStop(clusterName string) error {

	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	containers, err := utils.GetContainers(ctx, docker, clusterName, model.ServiceContainers)
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to list containers - %s", err.Error()), nil)
		return err
	}

	for _, container := range containers {
		// First stop the container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to stop container (%s)", container.ID), err)
		}
	}

	argsSC := filters.Arg("label", "app=space-cloud")
	argsNetwork := filters.Arg("network", utils.GetNetworkName(clusterName))
	scContainers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsNetwork, argsSC), All: true})
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
	argsNetwork = filters.Arg("network", utils.GetNetworkName(clusterName))
	addOnContainers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsNetwork, argsAddOns), All: true})
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
