package operations

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Destroy cleans the environment which has been setup. It removes the containers, secrets & host file
func Destroy(clusterName string) error {
	utils.LogInfo(fmt.Sprintf("Destroying the cluster (%s)...", clusterName))
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to initialize docker client - %s", err.Error()), nil)
		return err
	}

	// get all containers containing < space-cloud > in their name
	containers, err := utils.GetContainers(ctx, cli, clusterName, model.AllContainers)
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to list containers - %s", err.Error()), nil)
		return err
	}

	// Remove all container
	for _, containerInfo := range containers {
		if containerInfo.Labels["service"] == string(model.ImageRunner) {
			go func() {
				// Make sure the container is running before deleting secrets
				if err := cli.ContainerStart(ctx, containerInfo.ID, types.ContainerStartOptions{}); err != nil {
					_ = utils.LogError("Unable to start container to delete secrets", err)
					return
				}
				// NOTE: files are created with root permission in runner. If host system want to delete these files it requires root permissions.
				// so to delete files without root permission we remove the files from container itself
				execProcess, err := cli.ContainerExecCreate(ctx, containerInfo.ID, types.ExecConfig{Cmd: []string{"rm", "-rf", "/secrets"}})
				if err != nil {
					_ = utils.LogError("Unable to create delete secrets execution command", err)
					return
				}
				if err := cli.ContainerExecStart(ctx, execProcess.ID, types.ExecStartCheck{}); err != nil {
					_ = utils.LogError("Unable to execute delete secrets command", err)
					return
				}
			}()
		}
		if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			_ = utils.LogError(fmt.Sprintf("Unable to remove container %s - %s", containerInfo.ID, err.Error()), nil)
			return err
		}
	}

	// Remove the space-cloud network
	args := filters.Arg("name", utils.GetNetworkName(clusterName))
	nws, err := cli.NetworkList(ctx, types.NetworkListOptions{Filters: filters.NewArgs(args)})
	if err != nil {
		return utils.LogError("Unable to list networks", err)
	}
	for _, nw := range nws {
		if nw.Name == utils.GetNetworkName(clusterName) {
			_ = cli.NetworkRemove(ctx, nw.ID)
		}
	}

	// Remove secrets directory
	if err := os.RemoveAll(utils.GetSecretsDir(clusterName)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove secrets directory - %s", err.Error()), nil)
		return err
	}

	// Remove host file
	if err := os.RemoveAll(utils.GetSpaceCloudHostsFilePath(clusterName)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove host file - %s", err.Error()), nil)
		return err
	}

	// Remove the service routing file
	if err := os.RemoveAll(utils.GetSpaceCloudRoutingConfigPath(clusterName)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove service routing file file - %s", err.Error()), nil)
		return err
	}

	// Remove the config file
	if err := os.RemoveAll(utils.GetSpaceCloudConfigFilePath(clusterName)); err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to remove config file file - %s", err.Error()), nil)
		return err
	}

	if clusterName != "default" {
		// Remove the directory
		if err := os.RemoveAll(utils.GetSpaceCloudClusterDirectory(clusterName)); err != nil {
			_ = utils.LogError(fmt.Sprintf("Unable to remove cluster dir- %s", err.Error()), nil)
			return err
		}
	}

	if err := utils.RemoveAccount(clusterName); err != nil {
		return err
	}
	utils.LogInfo("Space cloud cluster has been destroyed successfully 😢")
	utils.LogInfo("Looking forward to seeing you again! 😊")
	return nil
}
