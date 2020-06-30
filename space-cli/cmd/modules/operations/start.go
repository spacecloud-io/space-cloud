package operations

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/txn2/txeh"
)

// DockerStart restarts the services which might have been stopped for any reason
func DockerStart(clusterName string) error {
	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(clusterName), WriteFilePath: utils.GetSpaceCloudHostsFilePath(clusterName)})
	if err != nil {
		return utils.LogError("Unable to open hosts file", err)
	}

	// Start the add ons first
	argsAddOns := filters.Arg("label", "app=addon")
	argsNetwork := filters.Arg("network", utils.GetNetworkName(clusterName))
	addOnContainers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsNetwork, argsAddOns), All: true})
	if err != nil {
		return utils.LogError("Unable to list space-cloud core containers", err)
	}
	for _, container := range addOnContainers {
		// First start the container
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%s)", container.ID), err)
		}

		// Get the container's info
		info, err := docker.ContainerInspect(ctx, container.ID)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inspect container (%s)", container.ID), err)
		}

		hostName := utils.GetServiceDomain(info.Config.Labels["service"], info.Config.Labels["name"])

		// Remove the domain from the hosts file
		hosts.RemoveHost(hostName)

		// Add it back with the new ip address
		hosts.AddHost(info.NetworkSettings.Networks[utils.GetNetworkName(clusterName)].IPAddress, hostName)
	}

	// Save the hosts file before continuing
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", err)
	}

	// Second step is to start the gateway and runner. We'll need the runner's ip address in the next step
	var runnerIP string
	argsSC := filters.Arg("label", "app=space-cloud")
	argsNetwork = filters.Arg("network", utils.GetNetworkName(clusterName))
	scContainers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsNetwork, argsSC), All: true})
	if err != nil {
		return utils.LogError("Unable to list space-cloud core containers", err)
	}

	for _, container := range scContainers {
		// First start the container
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%s)", container.ID), err)
		}

		// Get the container's info
		info, err := docker.ContainerInspect(ctx, container.ID)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inspect container (%s)", container.ID), err)
		}

		// Set the runner ip with the current service is the runner
		name := info.Config.Labels["service"]
		hostName := fmt.Sprintf("%s.space-cloud.svc.cluster.local", name)

		if name == "runner" {
			runnerIP = info.NetworkSettings.Networks[utils.GetNetworkName(clusterName)].IPAddress
		}

		// Remove the domain from the hosts file
		hosts.RemoveHost(hostName)

		// Add it back with the new ip address
		hosts.AddHost(info.NetworkSettings.Networks[utils.GetNetworkName(clusterName)].IPAddress, hostName)
	}

	// Check if the ip address is set
	if runnerIP == "" {
		return utils.LogError("Unable to set ip address of runner. Did you run space-cli setup once?", nil)
	}

	// Save the hosts file before continuing
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating sc containers", err)
	}

	// Get the list of the other containers we need to start
	argsNetwork = filters.Arg("network", utils.GetNetworkName(clusterName))
	argsServices := filters.Arg("label", "app=service")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsNetwork, argsServices), All: true})
	if err != nil {
		return utils.LogError("Unable to list space-cloud services containers", err)
	}

	// Loop over each container and start them
	for _, container := range containers {
		// First start the container
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to start container (%s)", container.ID), err)
		}

		// Get the container's info
		info, err := docker.ContainerInspect(ctx, container.ID)
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to inspect container (%s)", container.ID), err)
		}

		projectID := info.Config.Labels["project"]
		serviceID := info.Config.Labels["service"]
		version := info.Config.Labels["version"]

		generalDomain := utils.GetServiceDomain(projectID, serviceID)
		internalDomain := utils.GetInternalServiceDomain(projectID, serviceID, version)

		// Remove the two domains from the hosts file
		hosts.RemoveHosts([]string{generalDomain, internalDomain})

		// Add them back. The general domain will point to the runner while the internal service domain will point to the actual container
		hosts.AddHost(runnerIP, generalDomain)
		hosts.AddHost(info.NetworkSettings.Networks[utils.GetNetworkName(clusterName)].IPAddress, internalDomain)
	}

	// Save the hosts file before continuing
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating services containers", err)
	}

	return nil
}
