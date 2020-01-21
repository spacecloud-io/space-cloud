package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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

// ApplyService creates containers for specified service
func (d *docker) ApplyService(service *model.Service, token string) error {
	ctx := context.Background()
	// remove containers if already exits
	if err := d.DeleteService(service.ID); err != nil {
		logrus.Errorf("error applying service in docker unable delete existing containers - %v", err)
		return err
	}

	// todo check host file overiding concern
	// client for CRUD operation on host file
	// default location /etc/hosts
	hostFile, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("error applying service in docker unable to load host file with suitable default - %v", err)
		return err
	}

	for _, task := range service.Tasks {
		// todo get image
		// for now images have been created locally but not uploaded on docker hub
		out, err := d.client.ImagePull(ctx, task.Docker.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, out)

		// TODO check docker runtime field in host config
		// expose ports of docker container as specified in task
		exposedPorts := map[nat.Port]struct{}{}
		for _, port := range task.Ports {
			portString := strconv.Itoa(int(port.Port))
			portWithProtocol := nat.Port(fmt.Sprintf("%s/%s", portString, port.Protocol))
			exposedPorts[portWithProtocol] = struct{}{}
		}

		// set environment variables of docker container
		task.Env["URL"] = fmt.Sprintf("http://service.artifact.svc.cluster.local:4122")
		task.Env["TOKEN"] = fmt.Sprintf("%s", token) // todo token
		envs := []string{}
		for envName, envValue := range task.Env {
			envs = append(envs, fmt.Sprintf("%s=%s", envName, envValue))
		}

		// store docker image name in labels so that we get it back for constructing service struct in get service
		service.Labels[fmt.Sprintf("dockerImage-%s", task.ID)] = task.Docker.Image

		containerName := fmt.Sprintf("%s-%s-%s-%s", service.ID, service.ProjectID, task.ID, service.Version)
		resp, err := d.client.ContainerCreate(ctx, &container.Config{
			Image:        task.Docker.Image,
			Env:          envs,
			ExposedPorts: exposedPorts,
			Labels:       service.Labels,
		}, &container.HostConfig{
			// receiving memory in mega bytes converting into bytes
			// convert received mill cpus to cpus by diving by 1000 then multiply with 100000 to get cpu quota
			Resources: container.Resources{Memory: task.Resources.Memory * 1024 * 1024, CPUQuota: task.Resources.CPU * 100},
		}, nil, containerName)
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

	// TODO CHECK IF THE BELOW FUNCTION IS NEEDED WHILE TESTING
	// if err := hostFile.Save(); err != nil {
	// 	logrus.Errorf("error applying service in docker unable to save host file - %v", err)
	// 	return err
	// }
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
		// there will be only 1 container name set
		if len(containerInfo.Names) != 1 {
			logrus.Errorf("error deleting service in docker containers length not equal to one")
			return fmt.Errorf("error deleting service in docker containers length not equal to one")
		}
		// container name < serviceId-projectId-taskId-version >
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
func (d *docker) GetService(serviceId, projectId, version string) (*model.Service, error) {
	ctx := context.Background()
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		logrus.Errorf("error getting service in docker unable to list containers got error message - %v", err)
		return nil, err
	}

	service := new(model.Service)
	service.ID = serviceId
	service.ProjectID = projectId
	service.Version = version
	service.Whitelist = []string{fmt.Sprintf("%s:*", projectId)}
	service.Upstreams = []model.Upstream{{ProjectID: projectId, Service: "*"}}
	tasks := []model.Task{}
	for _, containerInfo := range containers {
		// while creating containers we assign single name to the container that is in the form < serviceId-taskId >
		// so if length not equal to 1 throw error
		if len(containerInfo.Names) != 1 {
			logrus.Errorf("error getting service in docker containers length not equal to one")
			return nil, err
		}
		containerName := containerInfo.Names[0]
		// container name < serviceId-projectId-taskId-version >
		if strings.Split(containerName, "-")[0] == serviceId {
			containerInspect, err := d.client.ContainerInspect(ctx, containerName)
			if err != nil {
				logrus.Errorf("error getting service in docker unable to inspect container - %v", err)
				return nil, err
			}

			service.Labels = containerInfo.Labels

			// set ports of task
			task := model.Task{}
			ports := []model.Port{}
			for portWithProtocol := range containerInspect.Config.ExposedPorts {
				arr := strings.Split(string(portWithProtocol), "/")
				portNumber, err := strconv.Atoi(arr[0])
				if err != nil {
					logrus.Errorf("error getting service in docker unable to convert string to int got error message - %v", err)
					return nil, err
				}
				ports = append(ports, model.Port{Protocol: model.Protocol(arr[1]), Port: int32(portNumber)}) // port name remaining
			}
			task.Ports = ports

			// set environment variable of task
			envs := map[string]string{}
			for _, value := range containerInspect.Config.Env {
				env := strings.Split(value, "=")
				envs[env[0]] = env[1]
			}
			delete(envs, "URL")
			delete(envs, "TOKEN")
			task.Env = envs

			// container name < serviceId-projectId-taskId-version >
			task.ID = strings.Split(containerName, "-")[3]

			// set task resource struct
			task.Resources.CPU = containerInspect.HostConfig.CPUShares
			task.Resources.Memory = containerInspect.HostConfig.Memory / (1024 * 1024)

			// set docker struct values
			task.Docker.Cmd = []string{envs["CMD"]}
			task.Docker.Image = containerInfo.Labels[fmt.Sprintf("dockerImage-%s", task.ID)]

			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		// through error as containerInfo doesn't exits
		logrus.Errorf("error getting service in docker specified service not found")
		return nil, fmt.Errorf("error getting service in docker specified service not found")
	}
	return service, nil
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
