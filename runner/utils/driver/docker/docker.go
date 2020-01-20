package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cloud/runner/model"
)

type docker struct {
	client *client.Client
}

func NewDockerDriver() (*docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error creating docker module instance in docker in docker unable to initialize docker client - %v", err)
		return nil, err
	}
	return &docker{client: cli}, nil
}

// ApplyService creates container for specified service
func (d *docker) ApplyService(service *model.Service, token string) error {
	ctx := context.Background()
	// remove containers if already exits
	if err := d.DeleteService(service.ID); err != nil {
		logrus.Errorf("error applying service in docker unable delete existing containers - %v", err)
		return err
	}

	// todo check host file overiding concern
	// client for CRUD operation on host file mounted on docker container at directory < /space-cloud-ip-table/.space-cloud-ip_table.hostFile >
	hostFile, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("error applying service in docker unable to load host file with suitable default - %v", err)
		return err
	}

	isHostAddr, hostAddr, _ := hostFile.HostAddressLookup("service.artifact-task.svc.cluster.local")
	if !isHostAddr {
		logrus.Errorf("error applying serivce in docker unable to load artifact store address from host file")
		return fmt.Errorf("error applying serivce in docker unable to load artifact store address from host file")
	}

	serviceJsonString, err := json.Marshal(service)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal service - %v", err)
		return err
	}

	for _, task := range service.Tasks {
		// for now images have been created locally but not uploaded on docker hub
		out, err := d.client.ImagePull(ctx, task.Docker.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, out)

		containerName := fmt.Sprintf("%s-%s", service.ID, task.ID)
		resp, err := d.client.ContainerCreate(ctx, &container.Config{
			Image: task.Docker.Image,
			Env: []string{
				fmt.Sprintf("URL=%s", fmt.Sprintf("http://%s:4122", hostAddr)),
				fmt.Sprintf("TOKEN=%s", token), // todo token
				fmt.Sprintf("CMD=%s", task.Env["CMD"]),
			},
			Labels: map[string]string{"service": string(serviceJsonString)},
		}, nil, nil, containerName)
		if err != nil {
			logrus.Errorf("error applying service in docker unable to create container %s got error message - %v", containerName, err)
			return err
		}

		if err := d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			logrus.Errorf("error applying service in docker unable to start container %s got error message - %v", containerName, err)
			return err
		}

		// get ip address of service & store it in host file
		data, err := d.client.ContainerInspect(ctx, resp.ID)
		if err != nil {
			logrus.Errorf("error applying service in docker unable to inspect container %s got error message  -%v", containerName, err)
			return err
		}
		hostFile.AddHost(data.NetworkSettings.IPAddress, fmt.Sprintf("%s.%s-%s.svc.cluster.local", service.ID, service.ProjectID, task.ID))
	}

	if err := hostFile.Save(); err != nil {
		logrus.Errorf("error applying service in docker unable to save host file - %v", err)
		return err
	}

	return nil
}

// DeleteService removes every docker container related to specified service id
func (d *docker) DeleteService(serviceId string) error {
	ctx := context.Background()
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		logrus.Errorf("error deleting service in docker unable to list containers got error message - %v", err)
		return err
	}

	for _, containerInfo := range containers {
		if len(containerInfo.Names) != 1 {
			logrus.Errorf("error deleting service in docker containers length not equal to one")
			return fmt.Errorf("error deleting service in docker containers length not equal to one")

		}
		containerName := containerInfo.Names[0]
		if strings.Split(containerName, "-")[0] == serviceId {
			// stop the container forcefully if status is running
			if containerInfo.Status == "running" { // todo check the status
				if err := d.client.ContainerKill(ctx, containerName, "SIGKILL"); err != nil {
					logrus.Errorf("error deleting service in docker unable to kill container %s got error message - %v", containerName, err)
					return err
				}
			}

			// remove the container from host machine
			if err := d.client.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{}); err != nil {
				logrus.Errorf("error deleting service in docker unable to remove container %s got error message - %v", containerName, err)
				return err
			}
		}
	}

	// handle gracefully if no containers found for specified serviceId
	return nil
}

// GetService gets the specified service info from docker container
func (d *docker) GetService(serviceId string) (*model.Service, error) {
	ctx := context.Background()
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		logrus.Errorf("error getting service in docker unable to list containers got error message - %v", err)
		return nil, err
	}

	for _, containerInfo := range containers {
		// while creating containers we assign single name to the container that is in the form < serviceId-taskId >
		// so if length not equal to 1 throw error
		if len(containerInfo.Names) != 1 {
			logrus.Errorf("error getting service in docker containers length not equal to one")
			return nil, err
		}
		containerName := containerInfo.Names[0]
		// container name < serviceId-taskId >
		if strings.Split(containerName, "-")[0] == serviceId {
			serviceInfo, ok := containerInfo.Labels["service"]
			if !ok {
				logrus.Errorf("error getting service in docker container does not have label with key named service")
				return nil, err
			}
			service := new(model.Service)
			if err := json.Unmarshal([]byte(serviceInfo), service); err != nil {
				logrus.Errorf("error getting service in docker unable to unmarshal serviceInfo got from container labels got error message - %v", err)
				return nil, err
			}
			return service, nil
		}
	}

	// through error as containerInfo doesn't exits
	logrus.Errorf("error getting service in docker specified service not found")
	return nil, fmt.Errorf("error getting service in docker specified service not found")
}

func (d *docker) CreateProject(project *model.Project) error {
	return nil
}

func (d *docker) AdjustScale(service *model.Service, activeReqs int32) error {
	return nil
}

func (d *docker) WaitForService(service *model.Service) error {
	return nil
}

func (d *docker) Type() model.DriverType {
	return model.TypeDocker
}
