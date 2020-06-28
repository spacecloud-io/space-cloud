package addons

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/spaceuptech/space-cli/cmd/modules/operations"
	"github.com/spaceuptech/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func addRegistry(projectID string) error {
	ctx := context.Background()
	dockerImage := "registry:2"

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Check if a registry container already exist
	filterArgs := filters.Arg("label", "service=registry")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if registry already exists", err)
	}
	if len(containers) != 0 {
		utils.LogInfo("Registry already exists. Do you want to remove it?")
		return nil
	}

	// Pull image if it doesn't already exist
	if err := utils.PullImageIfNotExist(ctx, docker, dockerImage); err != nil {
		return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", dockerImage), err)
	}

	// Check if projectID id is valid. If no projectID was provided
	if projectID == "" {
		utils.LogInfo("Project id not provided. Fetching projects from Space Cloud...")

		// Get projectID list from space cloud
		projects, err := utils.GetProjectsFromSC()
		if err != nil {
			return utils.LogError("Could not fetch list of projects from Space Cloud. Did you run `space-cli setup` once?", err)
		}

		// Throw error if no project has been created
		if len(projects) == 0 {
			return utils.LogError("No projects found. Run this command after creating a project", err)
		}

		projectID = projects[0].ID
		if len(projects) > 1 {
			var projectIDOptions []string
			for _, projectInfo := range projects {
				projectIDOptions = append(projectIDOptions, projectInfo.ID)
			}

			if err := survey.AskOne(&survey.Select{Message: "Select project ID", Options: projectIDOptions}, &projectID); err != nil {
				return err
			}
		}

		utils.LogInfo(fmt.Sprintf("Adding registry to project - %s with ID - %s", projects[0].Name, projectID))
	}

	// Set registry config in SpaceCloud. We will first get the projectID config, then apply the registry url to it
	specObj, err := project.GetProjectConfig(projectID, "project", nil)
	if err != nil {
		return utils.LogError(fmt.Sprintf("Unable to fetch project config of project (%s)", projectID), err)
	}
	if len(specObj) == 0 {
		return utils.LogError(fmt.Sprintf("No project found with id (%s)", projectID), err)
	}
	specObj[0].Spec.(map[string]interface{})["dockerRegistry"] = "localhost:5000"

	account, token, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return utils.LogError("Couldn't get account details or login token", err)
	}

	if err := operations.ApplySpec(token, account, specObj[0]); err != nil {
		return utils.LogError(fmt.Sprintf("Unable to update project (%s) with docker registry url by spec object with id (%v) type (%v)", projectID, specObj[0].Meta["id"], specObj[0].Type), err)
	}

	// Create the registry
	containerRes, err := docker.ContainerCreate(ctx, &container.Config{
		Labels:       map[string]string{"app": "addon", "service": "registry", "name": "registry"},
		Image:        dockerImage,
		ExposedPorts: nat.PortSet{"5000": struct{}{}},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{"5000": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5000"}}},
		NetworkMode:  "space-cloud",
	}, nil, "space-cloud--addon--registry")
	if err != nil {
		return utils.LogError("Unable to create local docker registry", err)
	}

	// Start the registry
	if err := docker.ContainerStart(ctx, containerRes.ID, types.ContainerStartOptions{}); err != nil {
		return utils.LogError("Unable to start local docker registry", err)
	}

	return nil
}

func removeRegistry(projectID string) error {
	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Check if a registry container already exist
	filterArgs := filters.Arg("label", "service=registry")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if registry already exists", err)
	}
	if len(containers) == 0 {
		utils.LogInfo("No registry exists. Do you want to add one?")
		return nil
	}

	// Remove registry config in SpaceCloud. We will first get the projectID config, then apply the registry url to it
	specObj, err := project.GetProjectConfig(projectID, "project", nil)
	if err != nil {
		return utils.LogError(fmt.Sprintf("Unable to fetch project config of project (%s)", projectID), err)
	}
	if len(specObj) == 0 {
		return utils.LogError(fmt.Sprintf("No project found with id (%s)", projectID), err)
	}
	specObj[0].Spec.(map[string]interface{})["dockerRegistry"] = ""

	account, token, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return utils.LogError("Couldn't get account details or login token", err)
	}

	if err := operations.ApplySpec(token, account, specObj[0]); err != nil {
		return utils.LogError(fmt.Sprintf("Unable to remove project (%s) with docker registry url by spec object with id (%v) type (%v)", projectID, specObj[0].Meta["id"], specObj[0].Type), err)
	}

	// Remove all container
	for _, containerInfo := range containers {
		// remove the container from host machine
		if err := docker.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to remove container %s", containerInfo.ID), err)
		}
	}

	return nil
}
