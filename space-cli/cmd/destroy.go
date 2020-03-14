package cmd

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Destroy cleans the environment which has been setup. It removes the containers, secrets & host file
func Destroy() error {
	logrus.Infoln("Destroying the cluster...")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("Unable to initialize docker client - %s", err.Error())
		return err
	}

	// get all containers containing < space-cloud > in their name
	args := filters.Arg("name", "space-cloud")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("Unable to list containers - %s", err.Error())
		return err
	}

	// Remove all container
	for _, containerInfo := range containers {
		// remove the container from host machine
		if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			logrus.Errorf("Unable to remove container %s - %s", containerInfo.ID, err.Error())
			return err
		}
	}

	// Remove secrets directory
	if err := os.RemoveAll(getSecretsDir()); err != nil {
		logrus.Errorf("Unable to remove secrets directory - %s", err.Error())
		return err
	}

	// Remove host file
	if err := os.RemoveAll(getSpaceCloudHostsFilePath()); err != nil {
		logrus.Errorf("Unable to remove host file - %s", err.Error())
		return err
	}

	// Remove the service routing file
	if err := os.RemoveAll(getSpaceCloudRoutingConfigPath()); err != nil {
		logrus.Errorf("Unable to remove service routing file file - %s", err.Error())
		return err
	}

	logrus.Infoln("Space cloud cluster has been destroyed successfully 😢")
	logrus.Infoln("Looking forward to seeing you again! 😊")
	return nil
}
