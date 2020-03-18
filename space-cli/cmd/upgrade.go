package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	scClient "github.com/spaceuptech/space-api-go"
	spaceApiTypes "github.com/spaceuptech/space-api-go/types"
	"github.com/txn2/txeh"
)

type resultData struct {
	Docs []*doc `mapstructure:"result"`
}

type doc struct {
	ID                string `mapstructure:"_id" json:"id"`
	VersionNo         string `mapstructure:"version_no" json:"versionNo"`
	CompatibleVersion string `mapstructure:"compatible_version" json:"compatibleVersion"`
}

// Upgrade upgrades the environment which has been setup
func Upgrade() error {
	const ContainerGateway string = "space-cloud-gateway"
	const ContainerRunner string = "space-cloud-runner"

	result := make(map[string]interface{})
	if err := Get(http.MethodGet, "/v1/config/env", map[string]string{}, &result); err != nil {
		return err
	}

	currentVersion := result["version"].(string)
	latestVersion, err := getLatestVersion(currentVersion)
	if err != nil {
		return err
	}
	if currentVersion == latestVersion {
		return fmt.Errorf("current verion (%s) is up to date", currentVersion)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("Unable to initialize docker client - %s", err.Error())
		return err
	}

	// get all containers containing < space-cloud > in their name
	args := filters.Arg("label", "app=space-cloud")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("Unable to list containers - %s", err.Error())
		return err
	}

	//parameters for gateway
	var gatewayMounts []mount.Mount
	var gatewayPorts nat.PortMap
	var gatewayEnvs []string

	//parameters for runner
	var runnerEnvs []string
	var runnerMounts []mount.Mount

	// Remove all container
	for _, containerInfo := range containers {
		containerInspect, err := cli.ContainerInspect(ctx, containerInfo.ID)
		if err != nil {
			logrus.Errorf("error getting service in docker unable to inspect container - %v", err)
			return err
		}

		switch containerInspect.Config.Labels["service"] {
		case "gateway":

			gatewayEnvs = containerInspect.Config.Env

			gatewayMounts = containerInspect.HostConfig.Mounts

			gatewayPorts = containerInspect.HostConfig.PortBindings

			if err := cli.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				logrus.Errorf("Unable to remove container %s - %s", containerInfo.ID, err.Error())
				return err
			}

		case "runner":
			runnerEnvs = containerInspect.Config.Env

			runnerMounts = containerInspect.HostConfig.Mounts

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
			containerImage: fmt.Sprintf("%s:v%s", "spaceuptech/gateway", latestVersion),
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
			containerImage: fmt.Sprintf("%s:v%s", "spaceuptech/runner", latestVersion),
			containerName:  ContainerRunner,
			dnsName:        "runner.space-cloud.svc.cluster.local",
			envs:           runnerEnvs,
			mount:          runnerMounts,
		},
	}

	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("Unable to initialize docker client - %s", err)
		return err
	}

	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: getSpaceCloudHostsFilePath(), WriteFilePath: getSpaceCloudHostsFilePath()})
	if err != nil {
		logrus.Errorf("Unable to load host file with suitable default - %s", err)
		return err
	}

	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
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
		// Remove the domain from the hosts file
		hosts.RemoveHost(c.dnsName)
		// Add it back with the new ip address
		hosts.AddHost(data.NetworkSettings.Networks["space-cloud"].IPAddress, c.dnsName)
	}

	if err := hosts.Save(); err != nil {
		logrus.Errorf("Unable to save host file - %s", err.Error())
		return err
	}

	return nil
}

func getLatestVersion(version string) (string, error) {
	mongoConn := scClient.New("test", "localhost:4122", false).DB("mongo")
	if mongoConn == nil {
		return "", fmt.Errorf("cannot connect to mongo")
	}
	ctx := context.Background()

	result, _ := mongoConn.Get("table").Where(spaceApiTypes.Cond("compatible_version", "==", version)).Apply(ctx)

	r := new(resultData)
	if err := result.Unmarshal(&r); err != nil {
		return "", err
	}
	newVersion := version
	for _, val := range r.Docs {
		if val.VersionNo > newVersion {
			newVersion = val.VersionNo
			fmt.Println("new version: ", newVersion)
		}
	}
	return newVersion, nil
}
