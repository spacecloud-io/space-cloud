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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
)

type docker struct {
	client *client.Client
	auth   *auth.Module
}

func NewDockerDriver(auth *auth.Module) (*docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error creating docker module instance in docker in docker unable to initialize docker client - %v", err)
		return nil, err
	}
	return &docker{client: cli, auth: auth}, nil
}

// ApplyService creates containers for specified service
func (d *docker) ApplyService(ctx context.Context, service *model.Service) error {
	// remove containers if already exits
	if err := d.DeleteService(ctx, service.ID, service.ProjectID, service.Version); err != nil {
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

	service.Labels["internalRuntime"] = string(service.Runtime)

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
		token, err := d.auth.GenerateHS256Token(service.ID, service.ProjectID, service.Version)
		if err != nil {
			logrus.Errorf("error applying service in docker unable generate token - %v", err)
			return err
		}
		task.Env["TOKEN"] = token
		envs := []string{}
		for envName, envValue := range task.Env {
			envs = append(envs, fmt.Sprintf("%s=%s", envName, envValue))
		}
		//
		// service.Labels["internalServiceId"] = service.ID
		// service.Labels["internalProjectId"] = service.ProjectID

		containerName := fmt.Sprintf("%s--%s--%s--%s", service.ProjectID, service.ID, task.ID, service.Version)
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
func (d *docker) DeleteService(ctx context.Context, serviceId, projectId, version string) error {
	args := filters.Arg("name", fmt.Sprintf("%s-%s", serviceId, projectId))
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("error deleting service in docker unable to list containers got error message - %v", err)
		return err
	}

	for _, containerInfo := range containers {
		containerInspect, err := d.client.ContainerInspect(ctx, containerInfo.ID)
		if err != nil {
			logrus.Errorf("error getting service in docker unable to inspect container - %v", err)
			return err
		}
		containerName := strings.Split(strings.TrimPrefix(containerInspect.Name, "/"), "--")
		pId := containerName[0]
		sId := containerName[1]
		if sId == serviceId && pId == projectId {
			// remove the container from host machine
			if err := d.client.ContainerRemove(ctx, containerInspect.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				logrus.Errorf("error deleting service in docker unable to remove container %s got error message - %v", containerName, err)
				return err
			}
		}
	}
	// handle gracefully if no containers found for specified serviceId
	return nil
}

// GetServices gets the specified service info from docker container
func (d *docker) GetService(serviceId, projectId, version string) (*model.Service, error) {
	ctx := context.Background()
	args := filters.Arg("name", fmt.Sprintf("%s-%s", serviceId, projectId))
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
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
		containerInspect, err := d.client.ContainerInspect(ctx, containerInfo.ID)
		if err != nil {
			logrus.Errorf("error getting service in docker unable to inspect container - %v", err)
			return nil, err
		}
		containerName := strings.Split(strings.TrimPrefix(containerInspect.Name, "/"), "--")
		pId := containerName[0]
		sId := containerName[1]
		taskId := containerName[2]
		if sId == serviceId && pId == projectId {
			runtime := containerInspect.Config.Labels["internalRuntime"]
			service.Labels = containerInspect.Config.Labels
			delete(service.Labels, "internalRuntime")

			// set ports of task
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

			// set environment variable of task
			envs := map[string]string{}
			for _, value := range containerInspect.Config.Env {
				env := strings.Split(value, "=")
				envs[env[0]] = env[1]
			}
			if runtime == "code" {
				delete(envs, "URL")
				delete(envs, "TOKEN")
			}

			tasks = append(tasks, model.Task{
				ID: taskId,
				Docker: model.Docker{
					Image: containerInspect.Config.Image,
					Cmd:   []string{envs["CMD"]},
				},
				Resources: model.Resources{
					CPU:    containerInspect.HostConfig.Memory / (1024 * 1024),
					Memory: containerInspect.HostConfig.CPUQuota / 100,
				},
				Env:   envs,
				Ports: ports,
			})
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
	logrus.Debug("create project not implemented for docker")
	return nil
}

func (d *docker) AdjustScale(service *model.Service, activeReqs int32) error {
	logrus.Debug("adjust scale not implemented for docker")
	return nil
}

func (d *docker) WaitForService(service *model.Service) error {
	logrus.Debug("wait for service not implemented for docker")
	return nil
}

func (d *docker) Type() model.DriverType {
	return model.TypeDocker
}
