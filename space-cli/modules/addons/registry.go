package addons

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/modules/project"
	"github.com/spaceuptech/space-cli/utils"
)

func addRegistry(projectID string) error {
	ctx := context.Background()
	dockerImage := "registry:2"

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", "add", "registry", err)
	}

	// Check if a registry container already exist
	filterArgs := filters.Arg("label", "service=registry")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if registry already exists", "add", "registry", err)
	}
	if len(containers) != 0 {
		utils.LogInfo("Registry already exists. Do you want to remove it?", "add", "registry")
		return nil
	}

	// Pull image if it doesn't already exist
	if err := utils.PullImageIfNotExist(ctx, docker, dockerImage); err != nil {
		return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", dockerImage), "add", "registry", err)
	}

	// Check if projectID id is valid. If no projectID was provided
	if projectID == "" {
		utils.LogInfo("Project id not provided. Fetching projects from Space Cloud...", "add", "registry")

		// Get projectID list from space cloud
		projects, err := utils.GetProjectsFromSC()
		if err != nil {
			return utils.LogError("Could not fetch list of projects from Space Cloud. Did you run `space-cli setup` once?", "add", "registry", err)
		}

		// Throw error if no project has been created
		if len(projects) == 0 {
			return utils.LogError("No projects found. Run this command after creating a project", "add", "registry", err)
		}

		// TODO: Ask the user to select a projectID
		projectID = projects[0].ID
		utils.LogInfo(fmt.Sprintf("Adding registry to project - %s", projects[0].Name), "add", "registry")
	}

	// Set registry config in SpaceCloud. We will first get the projectID config, then apply the registry url to it
	specObj, err := project.GetProjectConfig(projectID, "project", nil)
	if err != nil {
		return utils.LogError(fmt.Sprintf("Unable to fetch project config of project (%s)", projectID), "add", "registry", err)
	}
	specObj.Spec.(map[string]interface{})["dockerRegistry"] = "localhost:5000"

	account, err := utils.GetSelectedAccount()
	if err != nil {
		return err
	}
	login, err := utils.Login(account)
	if err != nil {
		return err
	}

	if err := cmd.ApplySpec(login.Token, account, specObj); err != nil {
		return utils.LogError(fmt.Sprintf("Unable to update project (%s) with docker registry url", projectID), "add", "registry", err)
	}

	// Create the registry
	containerRes, err := docker.ContainerCreate(ctx, &container.Config{
		Labels:       map[string]string{"app": "addon", "service": "registry", "name": "registry"},
		Image:        dockerImage,
		ExposedPorts: nat.PortSet{"5000": struct{}{}},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{"5000": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5000"}}},
		NetworkMode:  "space-cloud",
	}, nil, "space-cloud--addon-registry")
	if err != nil {
		return utils.LogError("Unable to create local docker registry", "add", "registry", err)
	}

	// Start the registry
	if err := docker.ContainerStart(ctx, containerRes.ID, types.ContainerStartOptions{}); err != nil {
		return utils.LogError("Unable to start local docker registry", "add", "registry", err)
	}

	return nil
}
