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
	client       *client.Client
	auth         *auth.Module
	artifactAddr string
}

func NewDockerDriver(auth *auth.Module, artifactAddr string) (*docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error creating docker module instance in docker in docker unable to initialize docker client - %v", err)
		return nil, err
	}
	return &docker{client: cli, auth: auth, artifactAddr: artifactAddr}, nil
}

// ApplyService creates containers for specified service
func (d *docker) ApplyService(ctx context.Context, service *model.Service) error {
	service.Version = "v1"
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

	var containerName, containerIp string
	for index, task := range service.Tasks {
		if index == 0 {
			containerName, containerIp, err = d.createContainer(ctx, task, service, "")
			if err != nil {
				return err
			}
			hostFile.AddHost(containerIp, fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID))
		}
		_, _, err = d.createContainer(ctx, task, service, containerName)
		return err
	}

	// TODO CHECK IF THE BELOW FUNCTION IS NEEDED WHILE TESTING
	// if err := hostFile.Save(); err != nil {
	// 	logrus.Errorf("error applying service in docker unable to save host file - %v", err)
	// 	return err
	// }
	return nil
}

func (d *docker) createContainer(ctx context.Context, task model.Task, service *model.Service, cName string) (string, string, error) {
	out, err := d.client.ImagePull(ctx, task.Docker.Image, types.ImagePullOptions{})
	if err != nil {
		return "", "", err
	}
	io.Copy(os.Stdout, out)
	service.Labels["internalRuntime"] = string(task.Runtime)

	// expose ports of docker container as specified in task
	exposedPorts := map[nat.Port]struct{}{}
	for _, port := range task.Ports {
		portString := strconv.Itoa(int(port.Port))
		portWithProtocol := nat.Port(fmt.Sprintf("%s/%s", portString, port.Protocol))
		exposedPorts[portWithProtocol] = struct{}{}
	}

	if task.Runtime == model.Code {
		token, err := d.auth.GenerateTokenForArtifactStore(service.ID, service.ProjectID, service.Version)
		if err != nil {
			logrus.Errorf("error applying service in docker unable generate token - %v", err)
			return "", "", err
		}
		task.Env["ARTIFACT_URL"] = d.artifactAddr
		task.Env["ARTIFACT_TOKEN"] = token
		task.Env["ARTIFACT_PROJECT"] = service.ProjectID
		task.Env["ARTIFACT_SERVICE"] = service.ID
		task.Env["ARTIFACT_VERSION"] = service.Version
	}

	envs := []string{}
	for envName, envValue := range task.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", envName, envValue))
	}

	service.Labels["internalServiceId"] = service.ID
	service.Labels["internalProjectId"] = service.ProjectID

	hostConfig := &container.HostConfig{
		// receiving memory in mega bytes converting into bytes
		// convert received mill cpus to cpus by diving by 1000 then multiply with 100000 to get cpu quota
		Resources: container.Resources{Memory: task.Resources.Memory * 1024 * 1024, CPUQuota: task.Resources.CPU * 100},
	}
	if cName != "" {
		hostConfig.NetworkMode = container.NetworkMode("container:" + cName)
	}

	containerName := fmt.Sprintf("%s--%s--%s--%s", service.ProjectID, service.ID, task.ID, service.Version)
	resp, err := d.client.ContainerCreate(ctx, &container.Config{
		Image:        task.Docker.Image,
		Env:          envs,
		Cmd:          task.Docker.Cmd,
		ExposedPorts: exposedPorts,
		Labels:       service.Labels,
	}, hostConfig, nil, containerName)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to create container %s got error message - %v", containerName, err)
		return "", "", err
	}

	if err := d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		logrus.Errorf("error applying service in docker unable to start container %s got error message - %v", containerName, err)
		return "", "", err
	}

	// get ip address of service & store it in host file
	data, err := d.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to inspect container %s got error message  -%v", containerName, err)
		return "", "", err
	}
	return containerName, data.NetworkSettings.IPAddress, nil
}

// DeleteService removes every docker container related to specified service id
func (d *docker) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	args := filters.Arg("name", fmt.Sprintf("%s--%s", serviceID, projectID))
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
		if sId == serviceID && pId == projectID {
			// remove the container from host machine
			if err := d.client.ContainerRemove(ctx, containerInspect.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				logrus.Errorf("error deleting service in docker unable to remove container %s got error message - %v", containerName, err)
				return err
			}
		}
	}
	// handle gracefully if no containers found for specified serviceID
	return nil
}

// GetServices gets the specified service info from docker container
func (d *docker) GetServices(ctx context.Context, projectId string) ([]*model.Service, error) {
	args := filters.Arg("name", fmt.Sprintf("%s", projectId))
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("error getting service in docker unable to list containers got error message - %v", err)
		return nil, err
	}

	services := map[string]*model.Service{}
	for _, containerInfo := range containers {
		service := new(model.Service)

		containerInspect, err := d.client.ContainerInspect(ctx, containerInfo.ID)
		if err != nil {
			logrus.Errorf("error getting service in docker unable to inspect container - %v", err)
			return nil, err
		}
		containerName := strings.Split(strings.TrimPrefix(containerInspect.Name, "/"), "--")
		taskId := containerName[2]
		service.Version = containerName[3]
		service.ID = containerName[1]

		service.ProjectID = projectId
		service.Whitelist = []model.Whitelist{{ProjectID: projectId, Service: "*"}}
		service.Upstreams = []model.Upstream{{ProjectID: projectId, Service: "*"}}
		tasks := []model.Task{}
		existingService, ok := services[fmt.Sprintf("%s-%s", service.ID, service.Version)]
		if ok {
			tasks = existingService.Tasks
		}

		runtime := model.Runtime(containerInspect.Config.Labels["internalRuntime"])
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
		if runtime == model.Code {
			delete(envs, model.ArtifactURL)
			delete(envs, model.ArtifactToken)
			delete(envs, model.ArtifactProject)
			delete(envs, model.ArtifactService)
			delete(envs, model.ArtifactVersion)
		}

		tasks = append(tasks, model.Task{
			ID: taskId,
			Docker: model.Docker{
				Image: containerInspect.Config.Image,
				Cmd:   containerInspect.Config.Cmd,
			},
			Resources: model.Resources{
				CPU:    containerInspect.HostConfig.Memory / (1024 * 1024),
				Memory: containerInspect.HostConfig.CPUQuota / 100,
			},
			Env:   envs,
			Ports: ports,
		})
		service.Tasks = tasks
		services[fmt.Sprintf("%s-%s", service.ID, service.Version)] = service
	}

	serviceArr := []*model.Service{}
	for _, service := range services {
		serviceArr = append(serviceArr, service)
	}

	return serviceArr, nil
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
