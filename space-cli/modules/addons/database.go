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
)

func addDatabase(dbtype, username, password, alias, version string) error {
	ctx := context.Background()
	dockerImage := strings.Join([]string{dbtype, version}, ":")

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", "add", "database", err)
	}

	// Check if a registry container already exist
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

	return nil
}
