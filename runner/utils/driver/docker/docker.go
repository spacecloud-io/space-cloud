package docker

import (
	"context"
	"encoding/json"
	"fmt"
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

func (d *docker) CreateSecret(projectID string, secretObj *model.Secret) error {
	logrus.Debug("CreateSecret not implemented for docker")
	return nil
}

func (d *docker) ListSecrets(projectID string) ([]*model.Secret, error) {
	logrus.Debug("ListSecrets not implemented for docker")
	return nil, nil
}

func (d *docker) DeleteSecret(projectID, secretName string) error {
	logrus.Debug("DeleteSecret not implemented for docker")
	return nil
}

func (d *docker) SetKey(projectID, secretName, secretKey string, secretObj *model.SecretValue) error {
	logrus.Debug("SetKey not implemented for docker")
	return nil
}

func (d *docker) DeleteKey(projectID, secretName, secretKey string) error {
	logrus.Debug("DeleteKey not implemented for docker")
	return nil
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

	// Get the hosts file
	hostFile, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("Could not load host file with suitable default - %v", err)
		return err
	}

	// remove containers if already exits
	if err := d.DeleteService(ctx, service.ProjectID, service.ID, service.Version); err != nil {
		logrus.Errorf("error applying service in docker unable delete existing containers - %v", err)
		return err
	}

	// get all the ports to be exposed of all tasks
	ports := []model.Port{}
	for _, task := range service.Tasks {
		for _, port := range task.Ports {
			ports = append(ports, port)
		}
	}

	var containerName, containerIp string
	for index, task := range service.Tasks {
		if index == 0 {
			var err error
			containerName, containerIp, err = d.createContainer(ctx, task, service, ports, "")
			if err != nil {
				return err
			}
			hostFile.AddHost(containerIp, getServiceDomain(service.ProjectID, service.ID))
			continue
		}
		_, _, err := d.createContainer(ctx, task, service, []model.Port{}, containerName)
		return err
	}

	return hostFile.Save()
}

func (d *docker) createContainer(ctx context.Context, task model.Task, service *model.Service, overridePorts []model.Port, cName string) (string, string, error) {
	// TODO: pull the images
	// out, err := d.client.ImagePull(ctx, task.Docker.Image, types.ImagePullOptions{})
	// if err != nil {
	// 	return "", "", err
	// }
	// io.Copy(os.Stdout, out)

	// Create empty labels if not exists
	if service.Labels == nil {
		service.Labels = map[string]string{}
	}
	service.Labels["internalRuntime"] = string(task.Runtime)
	portsJsonString, err := json.Marshal(&task.Ports)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalPorts"] = string(portsJsonString)
	scaleJsonString, err := json.Marshal(&service.Scale)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalScale"] = string(scaleJsonString)

	affinityJsonString, err := json.Marshal(&service.Affinity)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalAffinity"] = string(affinityJsonString)

	whitelistJsonString, err := json.Marshal(&service.Whitelist)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalWhitelist"] = string(whitelistJsonString)

	upstreamJsonString, err := json.Marshal(&service.Upstreams)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalUpstream"] = string(upstreamJsonString)

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
		Resources: container.Resources{Memory: task.Resources.Memory * 1024 * 1024, NanoCPUs: task.Resources.CPU * 1000000},
	}

	exposedPorts := map[nat.Port]struct{}{}
	if cName != "" {
		hostConfig.NetworkMode = container.NetworkMode("container:" + cName)
	} else {
		// expose ports of docker container as specified for 1st task
		task.Ports = overridePorts // override all ports while creating container for 1st task
		for _, port := range task.Ports {
			portString := strconv.Itoa(int(port.Port))
			exposedPorts[nat.Port(portString)] = struct{}{}
		}
	}

	containerName := fmt.Sprintf("%s--%s--%s--%s", service.ProjectID, service.ID, service.Version, task.ID)
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
func (d *docker) DeleteService(ctx context.Context, projectId, serviceId, version string) error {
	args := filters.Arg("name", fmt.Sprintf("%s--%s--%s", projectId, serviceId, version))
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("error deleting service in docker unable to list containers got error message - %v", err)
		return err
	}

	for _, containerInfo := range containers {
		// remove the container from host machine
		if err := d.client.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			logrus.Errorf("error deleting service in docker unable to remove container %s got error message - %v", containerInfo.ID, err)
			return err
		}
	}

	// Remove host from hosts file
	hostFile, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("Could not load host file with suitable default - %v", err)
		return err
	}
	hostFile.RemoveHost(getServiceDomain(projectId, serviceId))
	return hostFile.Save()
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
		taskId := containerName[3]
		service.Version = containerName[2]
		service.ID = containerName[1]
		service.Name = service.ID

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
		ports := []model.Port{}
		if err := json.Unmarshal([]byte(service.Labels["internalPorts"]), &ports); err != nil {
			logrus.Errorf("error getting service in docker unable to unmarshal ports - %v", err)
			return nil, err
		}
		scale := model.ScaleConfig{}
		if err := json.Unmarshal([]byte(service.Labels["internalScale"]), &scale); err != nil {
			logrus.Errorf("error getting service in docker unable to unmarshal scale - %v", err)
			return nil, err
		}
		service.Scale = scale

		// Force scale to 1
		service.Scale.Replicas = 1

		whilteList := []model.Whitelist{}
		if err := json.Unmarshal([]byte(service.Labels["internalWhitelist"]), &whilteList); err != nil {
			logrus.Errorf("error getting service in docker unable to unmarshal whitelist - %v", err)
			return nil, err
		}
		service.Whitelist = whilteList

		upstream := []model.Upstream{}
		if err := json.Unmarshal([]byte(service.Labels["internalUpstream"]), &upstream); err != nil {
			logrus.Errorf("error getting service in docker unable to unmarshal upstream - %v", err)
			return nil, err
		}
		service.Upstreams = upstream

		affinity := []model.Affinity{}
		if err := json.Unmarshal([]byte(service.Labels["internalAffinity"]), &affinity); err != nil {
			logrus.Errorf("error getting service in docker unable to unmarshal affinity - %v", err)
			return nil, err
		}
		service.Affinity = affinity
		delete(service.Labels, "internalRuntime")
		delete(service.Labels, "internalPorts")
		delete(service.Labels, "internalProjectId")
		delete(service.Labels, "internalServiceId")
		delete(service.Labels, "internalScale")
		delete(service.Labels, "internalWhitelist")
		delete(service.Labels, "internalAffinity")
		delete(service.Labels, "internalUpstream")

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
			ID:   taskId,
			Name: taskId,
			Docker: model.Docker{
				Image: containerInspect.Config.Image,
				Cmd:   containerInspect.Config.Cmd,
			},
			Resources: model.Resources{
				Memory: containerInspect.HostConfig.Memory / (1024 * 1024),
				CPU:    containerInspect.HostConfig.NanoCPUs / 1000000,
			},
			Env:     envs,
			Ports:   ports,
			Runtime: runtime,
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
