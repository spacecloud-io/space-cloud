package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
)

// CodeSetup initializes development environment
func CodeSetup() error {
	const dockerImageSpaceCloud string = "spaceuptech/space-cloud"
	containersToCreate := []struct {
		containerName string
		commands      []string
		exposedPorts  nat.PortSet
		portMapping   nat.PortMap
	}{
		{
			containerName: "space-cloud-gateway",
			commands:      []string{"./space-cloud", "run", "-dev"},
			exposedPorts: nat.PortSet{
				"4122": struct{}{},
			},
			portMapping: nat.PortMap{
				"4122": []nat.PortBinding{{HostIP: "localhost", HostPort: "4122"}},
			},
		},
		{
			containerName: "space-cloud-runner",
			commands:      []string{"./space-cloud", "run", "-dev"},
		},

		{
			containerName: "space-cloud-artifact-store",
			commands:      []string{"./space-cloud", "run", "-dev"},
		},
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error cli setup unable to initialize docker client got error message - %v", err)
		return err
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("error cli setup unable to load host file with suitable default got error message - %v", err)
		return err
	}

	// pull image from docker hub
	out, err := cli.ImagePull(ctx, dockerImageSpaceCloud, types.ImagePullOptions{})
	if err != nil {
		logrus.Errorf("error cli setup unable to pull image from docker hub got error message - %v", err)
		return err
	}
	io.Copy(os.Stdout, out)

	if err := os.MkdirAll(fmt.Sprintf("%s/space-cloud-ip-table", getHomeDirectory()), 0755); err != nil {
		logrus.Errorf("error cli setup unable to create directory for storing host file got error message - %v", err)
		return err
	}

	hostFilePath := fmt.Sprintf("%s/space-cloud-ip-table/.space-cloud-ip_table.hosts", getHomeDirectory())
	hostFileDir := fmt.Sprintf("%s/space-cloud-ip-table", getHomeDirectory())

	if err := hosts.SaveAs(hostFilePath); err != nil {
		logrus.Errorf("error cli setup unable to save host file to specified path (%s) got error message - %v", hostFilePath, err)
		return err
	}

	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
	hosts.WriteFilePath = hostFilePath

	for _, c := range containersToCreate {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        dockerImageSpaceCloud,
			Cmd:          c.commands,
			ExposedPorts: c.exposedPorts,
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: hostFileDir,
					Target: "/space-cloud-ip-table",
				},
			},
			PortBindings: c.portMapping,
		}, nil, c.containerName)
		if err != nil {
			logrus.Errorf("error cli setup unable to create container %s got error message  - %v", c.containerName, err)
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			logrus.Errorf("error cli setup unable to start container %s got error message - %v", c.containerName, err)
			return err
		}

		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			logrus.Errorf("error cli setup unable to inspect container %s got error message - %v", c.containerName, err)
		}
		log.Println("container name:", data.Name)
		hosts.AddHost(data.NetworkSettings.IPAddress, c.containerName)
	}

	if err := hosts.Save(); err != nil {
		logrus.Errorf("error cli setup unable to save host file got error message - %v", err)
		return err
	}
	return nil
}
