package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
)

// Upgrade upgrades the environment which has been setup
func Upgrade() error {
	// TODO: old keys always remain in accounts.yaml file
	const ContainerGateway string = "space-cloud-gateway"
	const ContainerRunner string = "space-cloud-runner"

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

	//parameters for gateway
	gatewayMounts := []mount.Mount{}
	gatewayPorts := nat.PortMap{}
	gatewayEnvs := []string{}
	gatewayVersion := ""

	//parameters for runner
	runnerEnvs := []string{}
	runnerVersion := ""

	// Remove all container
	for _, containerInfo := range containers {
		containerInspect, err := cli.ContainerInspect(ctx, containerInfo.ID)
		if err != nil {
			logrus.Errorf("error getting service in docker unable to inspect container - %v", err)
			return err
		}

		imageName := strings.Split(containerInfo.Image, ":")
		if imageName[0] == "spaceuptech/gateway" {
			b, _ := json.MarshalIndent(containerInspect.NetworkSettings.Ports, "", " ")
			fmt.Println("container: ", string(b))
		}
		switch imageName[0] {
		case "spaceuptech/gateway":

			if err := json.Unmarshal([]byte(containerInspect.Config.Labels["internalPorts"]), &gatewayPorts); err != nil {
				logrus.Errorf("error getting service in docker unable to unmarshal ports - %v", err)
				return err
			}

			gatewayEnvs = containerInspect.Config.Env

			result := make(map[string]interface{})
			if err := Get(http.MethodGet, "/v1/config/env", map[string]string{}, &result); err != nil {
				return err
			}
			gatewayVersion = gatewayVersion + result["version"].(string)

			gatewayMounts = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudHostsFilePath(),
					Target: "/etc/hosts",
				},
			}

			gatewayPorts = containerInspect.NetworkSettings.Ports

			if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				logrus.Errorf("Unable to remove container %s - %s", containerInfo.ID, err.Error())
				return err
			}

		case "spaceuptech/runner":
			runnerEnvs = containerInspect.Config.Env

			result := make(map[string]interface{})
			if err := Get(http.MethodGet, "/v1/config/env", map[string]string{}, &result); err != nil {
				return err
			}
			runnerVersion = runnerVersion + result["version"].(string)

			if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				logrus.Errorf("Unable to remove container %s - %s", containerInfo.ID, err.Error())
				return err
			}
		}

	}
	containersToCreate := []struct {
		dnsName        string
		containerImage string
		containerName  string
		envs           []string
		mount          []mount.Mount
		exposedPorts   nat.PortSet
		portMapping    nat.PortMap
	}{
		{
			containerImage: "spaceuptech/gateway",
			containerName:  ContainerGateway,
			dnsName:        "gateway.space-cloud.svc.cluster.local",
			envs:           gatewayEnvs,
			exposedPorts: nat.PortSet{
				"4122": struct{}{},
				"4126": struct{}{},
			},
			portMapping: gatewayPorts,
			mount:       gatewayMounts,
		},

		{
			// runner
			containerImage: "spaceuptech/runner",
			containerName:  ContainerRunner,
			dnsName:        "runner.space-cloud.svc.cluster.local",
			envs:           runnerEnvs,
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind, // TODO CHECK THIS
					Source: getSecretsDir(),
					Target: "/secrets",
				},
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudHostsFilePath(),
					Target: "/etc/hosts",
				},
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudRoutingConfigPath(),
					Target: "/routing-config.json",
				},
			},
		},
	}

	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("Unable to initialize docker client - %s", err)
		return err
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("Unable to load host file with suitable default - %s", err)
		return err
	}
	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
	hosts.WriteFilePath = getSpaceCloudHostsFilePath()
	if err := hosts.SaveAs(getSpaceCloudHostsFilePath()); err != nil {
		logrus.Errorf("Unable to save as host file to specified path (%s) - %s", getSpaceCloudHostsFilePath(), err)
		return err
	}

	for _, c := range containersToCreate {
		logrus.Infof("Starting container %s...", c.containerName)
		// check if image already exists
		if err := pullImageIfNotExist(ctx, cli, c.containerImage); err != nil {
			logrus.Errorf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", c.containerImage)
			return err
		}

		// check if container is already running
		args := filters.Arg("name", c.containerName)
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
		if err != nil {
			logrus.Errorf("error deleting service in docker unable to list containers - %s", err)
			return err
		}
		if len(containers) != 0 {
			logrus.Errorf("Container (%s) already exists", c.containerName)
			return fmt.Errorf("container (%s) already exists", c.containerName)
		}

		// create container with specified defaults
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        c.containerImage,
			ExposedPorts: c.exposedPorts,
			Env:          c.envs,
		}, &container.HostConfig{
			Mounts:       c.mount,
			PortBindings: c.portMapping,
		}, nil, c.containerName)
		if err != nil {
			logrus.Errorf("Unable to create container (%s) - %s", c.containerName, err)
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			logrus.Errorf("Unable to start container (%s) - %s", c.containerName, err.Error())
			return err
		}

		// get the ip address assigned to container
		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			logrus.Errorf("Unable to inspect container (%s) - %s", c.containerName, err)
		}
		hosts.AddHost(data.NetworkSettings.IPAddress, c.dnsName)
	}

	if err := hosts.Save(); err != nil {
		logrus.Errorf("Unable to save host file - %s", err.Error())
		return err
	}
	return nil
}
