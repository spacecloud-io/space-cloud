package cmd

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Destroy cleans the environment which has been setup by the SETUP command
// it does the above by removing container, secrets & host file
func Destroy() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error cli setup unable to initialize docker client got error message - %v", err)
		return err
	}

	// get all containers containing < space--cloud > in their name
	args := filters.Arg("name", "space--cloud")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("error deleting service in docker unable to list containers got error message - %v", err)
		return err
	}

	// remove all container
	for _, containerInfo := range containers {
		// remove the container from host machine
		if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			logrus.Errorf("error deleting service in docker unable to remove container %s got error message - %v", containerInfo.ID, err)
			return err
		}
	}

	// remove secrets directory
	if err := os.RemoveAll(getSecretsDir()); err != nil {
		logrus.Errorf("error in destroy unable to remove secrets directory - %s", err.Error())
		return err
	}

	// remove host file
	if err := os.RemoveAll(getSpaceCloudHostsFilePath()); err != nil {
		logrus.Errorf("error in destroy unable to remove host file - %s", err.Error())
		return err
	}
	return nil
}
