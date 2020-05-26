package operations

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

// Destroy cleans the environment which has been setup. It removes the containers, secrets & host file
func Destroy(clusterID string) error {
	utils.LogInfo("Destroying the cluster...")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to initialize docker client - %s", err.Error()), nil)
		return err
	}

	// get all containers containing < space-cloud > in their name
	args := filters.Arg("name", getNetworkName(clusterID))
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to list containers - %s", err.Error()), nil)
		return err
	}

	gatewayName := "/" + getNetworkName(clusterID) + "-gateway"
	runnerName := "/" + getNetworkName(clusterID) + "-runner"

	// Remove all container
	for _, containerInfo := range containers {
		if containerInfo.Names[0] == gatewayName || containerInfo.Names[0] == runnerName {
			// remove the container from host machine
			if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				_ = utils.LogError(fmt.Sprintf("Unable to remove container %s - %s", containerInfo.ID, err.Error()), nil)
				return err
			}
		}
	}

	// Remove the space-cloud network
	nws, err := cli.NetworkList(ctx, types.NetworkListOptions{Filters: filters.NewArgs(args)})
	if err != nil {
		return utils.LogError("Unable to list networks", err)
	}
	for _, nw := range nws {
		if nw.Name == getNetworkName(clusterID) {
			_ = cli.NetworkRemove(ctx, nw.ID)
		}
	}

	// Remove secrets directory
	if err := os.RemoveAll(utils.GetSecretsDir(clusterID)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove secrets directory - %s", err.Error()), nil)
		return err
	}

	// Remove host file
	if err := os.RemoveAll(utils.GetSpaceCloudHostsFilePath(clusterID)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove host file - %s", err.Error()), nil)
		return err
	}

	// Remove the service routing file
	if err := os.RemoveAll(utils.GetSpaceCloudRoutingConfigPath()); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove service routing file file - %s", err.Error()), nil)
		return err
	}

	// Remove the config file
	if err := os.RemoveAll(utils.GetSpaceCloudConfigFilePath(clusterID)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove config file file - %s", err.Error()), nil)
		return err
	}

	if clusterID != "default" {
		// Remove the directory
		if err := os.RemoveAll(utils.GetSpaceCloudDirectory() + "/" + clusterID); err != nil {
			_ = utils.LogError(fmt.Sprintf("Unable to remove cluster dir- %s", err.Error()), nil)
			return err
		}
	}

	utils.LogInfo("Space cloud cluster has been destroyed successfully 😢")
	utils.LogInfo("Looking forward to seeing you again! 😊")
	return nil
}
