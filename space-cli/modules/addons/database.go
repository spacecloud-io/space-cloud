package addons

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/txn2/txeh"
)

func addDatabase(dbtype, username, password, alias, version string) error {
	ctx := context.Background()

	// Prepare the docker image name name
	dockerImage := fmt.Sprintf("%s:%s", dbtype, version)

	// Set alias if not provided
	if alias == "" {
		alias = dbtype
	}

	// Set the environment variables
	var env []string
	switch dbtype {
	case "mysql":
		if username != "root" {
			return utils.LogError("Only the username root is allowed for MySQL", nil)
		}
		env = []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password)}
	case "postgres":
		env = []string{fmt.Sprintf("POSTGRES_USER=%s", username), fmt.Sprintf("POSTGRES_PASSWORD=%s", password)}
	case "mongo":
		if username != "" || password != "" {
			return utils.LogError("Cannot set username or password with Mongo", nil)
		}
		env = []string{}
	default:
		return utils.LogError(fmt.Sprintf("Invalid database type (%s) provided", dbtype), nil)
	}

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Check if a database container already exist
	filterArgs := filters.Arg("label", "app=space-cloud")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if database already exists", err)
	}
	if len(containers) == 0 {
		utils.LogInfo("No space-cloud instance found. Run 'space-cli setup' first")
		return nil
	}

	// Pull image if it doesn't already exist
	if err := utils.PullImageIfNotExist(ctx, docker, dockerImage); err != nil {
		return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", dockerImage), err)
	}

	// Create the database
	containerRes, err := docker.ContainerCreate(ctx, &container.Config{
		Labels: map[string]string{"app": "addon", "service": "db", "name": alias},
		Image:  dockerImage,
		Env:    env,
	}, &container.HostConfig{
		NetworkMode: "space-cloud",
	}, nil, fmt.Sprintf("space-cloud--addon--db--%s", alias))
	if err != nil {
		return utils.LogError("Unable to create local docker database", err)
	}

	// Start the database
	if err := docker.ContainerStart(ctx, containerRes.ID, types.ContainerStartOptions{}); err != nil {
		return utils.LogError("Unable to start local docker database", err)
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to open hosts file", err)
	}

	// Get the container's info
	info, err := docker.ContainerInspect(ctx, containerRes.ID)
	if err != nil {
		return utils.LogError(fmt.Sprintf("Unable to inspect c (%s)", containerRes.ID), err)
	}

	hostName := utils.GetServiceDomain("db", alias)

	// Remove the domain from the hosts file
	hosts.RemoveHost(hostName)

	// Add it back with the new ip address
	hosts.AddHost(info.NetworkSettings.Networks["space-cloud"].IPAddress, hostName)

	// Save the hosts file
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", err)
	}

	utils.LogInfo(fmt.Sprintf("Started database (%s) with alias (%s)", dbtype, alias))
	return nil
}

func removeDatabase(alias string) error {
	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Check if a database container already exist
	filterArgs := filters.Arg("name", fmt.Sprintf("space-cloud--addon--db--%s", alias))
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if database already exists", err)
	}
	if len(containers) == 0 {
		utils.LogInfo(fmt.Sprintf("Database (%s) not found.", alias))
		return nil
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to open hosts file", err)
	}

	for _, c := range containers {
		hostName := utils.GetServiceDomain("db", alias)

		// Remove the domain from the hosts file
		hosts.RemoveHost(hostName)

		// remove the container from host machine
		if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to remove container %s", c.ID), err)
		}
	}

	// Save the hosts file
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", err)
	}

	utils.LogInfo(fmt.Sprintf("Removed database (%s)", alias))

	return nil
}
