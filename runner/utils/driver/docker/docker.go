package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	proxy_manager "github.com/spaceuptech/space-cloud/runner/utils/driver/docker/proxy-manager"
)

// Docker defines the type for docker instance
type Docker struct {
	client       *client.Client
	auth         *auth.Module
	artifactAddr string
	secretPath   string
	hostFilePath string
	manager      *proxy_manager.Manager
}

// NewDockerDriver returns a new docker instance
func NewDockerDriver(auth *auth.Module, artifactAddr string) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error creating docker module instance in docker in docker unable to initialize docker client - %v", err)
		return nil, err
	}

	secretPath := os.Getenv("SECRETS_PATH")
	if secretPath == "" {
		secretPath = "."
	}

	hostFilePath := os.Getenv("HOSTS_FILE_PATH")
	if hostFilePath == "" {
		logrus.Fatal("Failed to create docker driver: HOSTS_FILE_PATH environment variable not provided")
	}

	manager, err := proxy_manager.New(os.Getenv("ROUTING_FILE_PATH"))
	if err != nil {
		return nil, err
	}

	return &Docker{client: cli, auth: auth, artifactAddr: artifactAddr, secretPath: secretPath, hostFilePath: hostFilePath, manager: manager}, nil
}

// ApplyServiceRoutes sets the traffic splitting logic of each service
func (d *Docker) ApplyServiceRoutes(_ context.Context, projectID, serviceID string, routes model.Routes) error {
	return d.manager.SetServiceRoutes(projectID, serviceID, routes)
}

// GetServiceRoutes gets the routing rules of each service
func (d *Docker) GetServiceRoutes(_ context.Context, projectID string) (map[string]model.Routes, error) {
	return d.manager.GetServiceRoutes(projectID)
}

// ApplyService creates containers for specified service
func (d *Docker) ApplyService(ctx context.Context, service *model.Service) error {
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
		ports = append(ports, task.Ports...)
	}

	var containerName, containerIP string
	for index, task := range service.Tasks {
		if index == 0 {
			var err error
			containerName, containerIP, err = d.createContainer(ctx, index, task, service, ports, "")
			if err != nil {
				return err
			}
			hostFile.AddHost(containerIP, utils.GetInternalServiceDomain(service.ProjectID, service.ID, service.Version))
			continue
		}
		_, _, err := d.createContainer(ctx, index, task, service, []model.Port{}, containerName)
		return err
	}

	// Don't forget to set the service routing initially
	if err := d.manager.SetServiceRouteIfNotExists(service.ProjectID, service.ID, service.Version, ports); err != nil {
		logrus.Errorf("Could not create initial service routing for service (%s:%s)", service.ProjectID, service.ID)
		return err
	}

	// Point runner to Proxy (it's own IP address!)
	p, proxyIP, _ := hostFile.HostAddressLookup("runner.space-cloud.svc.cluster.local")
	if !p {
		return errors.New("no hosts entry found for runner domain")
	}

	hostFile.AddHost(proxyIP, utils.GetServiceDomain(service.ProjectID, service.ID))
	return hostFile.Save()
}

func (d *Docker) createContainer(ctx context.Context, index int, task model.Task, service *model.Service, overridePorts []model.Port, cName string) (string, string, error) {
	tempSecretPath := fmt.Sprintf("%s/temp-secrets/%s/%s", os.Getenv("SECRETS_PATH"), service.ProjectID, fmt.Sprintf("%s--%s", service.ID, service.Version))

	if err := d.pullImageByPolicy(ctx, service.ProjectID, task.Docker); err != nil {
		logrus.Error("error in docker unable to pull image ", err)
		return "", "", err
	}

	// Create empty labels if not exists
	if service.Labels == nil {
		service.Labels = map[string]string{}
	}

	// Overwrite important labels
	service.Labels["app"] = "service"
	service.Labels["project"] = service.ProjectID
	service.Labels["service"] = service.ID
	service.Labels["version"] = service.Version
	service.Labels["task"] = task.ID

	service.Labels["internalRuntime"] = string(task.Runtime)
	portsJSONString, err := json.Marshal(&task.Ports)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalPorts"] = string(portsJSONString)
	scaleJSONString, err := json.Marshal(&service.Scale)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalScale"] = string(scaleJSONString)

	affinityJSONString, err := json.Marshal(&service.Affinity)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalAffinity"] = string(affinityJSONString)

	whitelistJSONString, err := json.Marshal(&service.Whitelist)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalWhitelist"] = string(whitelistJSONString)

	upstreamJSONString, err := json.Marshal(&service.Upstreams)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalUpstream"] = string(upstreamJSONString)

	secretsJSONString, err := json.Marshal(&task.Secrets)
	if err != nil {
		logrus.Errorf("error applying service in docker unable to marshal ports - %v", err)
		return "", "", err
	}
	service.Labels["internalSecrets"] = string(secretsJSONString)

	service.Labels["internalDockerSecrets"] = task.Docker.Secret

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

	// set secrets
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: d.hostFilePath,
			Target: "/etc/hosts",
		},
	}
	for _, secretName := range task.Secrets {

		// check if file exists
		filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, service.ProjectID, secretName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", "", fmt.Errorf("secret cannot be read - %v", err)
		}

		// file already exists read it's content
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return "", "", err
		}
		fileContent := new(model.Secret)
		if err := json.Unmarshal(data, fileContent); err != nil {
			return "", "", err
		}
		switch fileContent.Type {
		case model.FileType:
			uniqueDirName := fmt.Sprintf("%s--%s", task.ID, secretName)
			path := fmt.Sprintf("%s/%s", tempSecretPath, uniqueDirName)
			if err := os.MkdirAll(path, 0777); err != nil {
				return "", "", err
			}
			for key, value := range fileContent.Data {
				if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", path, key), []byte(value), 0777); err != nil {
					return "", "", err
				}
			}
			mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: fmt.Sprintf("%s/%s/%s/%s", os.Getenv("HOME_SECRETS_PATH"), service.ProjectID, fmt.Sprintf("%s--%s", service.ID, service.Version), uniqueDirName), Target: fileContent.RootPath})
		case model.EnvType:
			for key, value := range fileContent.Data {
				envs = append(envs, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	service.Labels["internalServiceId"] = service.ID
	service.Labels["internalProjectId"] = service.ProjectID
	service.Labels["internalPullPolicy"] = string(task.Docker.ImagePullPolicy)

	hostConfig := &container.HostConfig{
		// receiving memory in mega bytes converting into bytes
		// convert received mill cpus to cpus by diving by 1000 then multiply with 100000 to get cpu quota
		Resources: container.Resources{Memory: task.Resources.Memory * 1024 * 1024, NanoCPUs: task.Resources.CPU * 1000000},
		Mounts:    mounts,
	}

	exposedPorts := map[nat.Port]struct{}{}
	if cName != "" {
		hostConfig.NetworkMode = container.NetworkMode("container:" + cName)
	} else {
		hostConfig.NetworkMode = "space-cloud"
		// expose ports of docker container as specified for 1st task
		task.Ports = overridePorts // override all ports while creating container for 1st task
		for _, port := range task.Ports {
			portString := strconv.Itoa(int(port.Port))
			exposedPorts[nat.Port(portString)] = struct{}{}
		}
	}

	containerName := fmt.Sprintf("space-cloud-%s--%s--%s--%d--%s", service.ProjectID, service.ID, service.Version, index, task.ID)
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
	return containerName, data.NetworkSettings.Networks["space-cloud"].IPAddress, nil
}

// DeleteService removes every docker container related to specified service id
func (d *Docker) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	args := filters.Arg("name", fmt.Sprintf("space-cloud-%s--%s--%s", projectID, serviceID, version))
	if serviceID == "" || version == "" {
		args = filters.Arg("name", projectID)
	}

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("error deleting service in docker unable to list containers got error message - %v", err)
		return err
	}
	// todo remove secret directory

	for _, containerInfo := range containers {
		// remove the container from host machine
		if err := d.client.ContainerRemove(ctx, containerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			logrus.Errorf("error deleting service in docker unable to remove container %s got error message - %v", containerInfo.ID, err)
			return err
		}
	}

	tempSecretPath := fmt.Sprintf("%s/temp-secrets/%s/%s", os.Getenv("SECRETS_PATH"), projectID, fmt.Sprintf("%s--%s", serviceID, version))
	if err := os.RemoveAll(tempSecretPath); err != nil {
		return err
	}

	// Remove host from hosts file
	hostFile, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("Could not load host file with suitable default - %v", err)
		return err
	}

	hostFile.RemoveHost(utils.GetInternalServiceDomain(projectID, serviceID, version))

	// Get if current version was the last remaining service
	isLastContainer, err := d.checkIfLastService(ctx, projectID, serviceID)
	if err != nil {
		return err
	}

	// Remove the general service along with the service routes if this was the last version of the service
	if isLastContainer {
		hostFile.RemoveHost(utils.GetServiceDomain(projectID, serviceID))
		if err := d.manager.DeleteServiceRoutes(projectID, serviceID); err != nil {
			logrus.Errorf("Could not remove service routing for service (%s:%s)", projectID, serviceID)
		}
	}

	return hostFile.Save()
}

func (d *Docker) checkIfLastService(ctx context.Context, projectID, serviceID string) (bool, error) {
	args := filters.Arg("name", fmt.Sprintf("space-cloud-%s--%s", projectID, serviceID))
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(args), All: true})
	if err != nil {
		logrus.Errorf("Could not list remaining containers got error message - %v", err)
		return false, err
	}

	return len(containers) == 0, nil
}

// GetServices gets the specified service info from docker container
func (d *Docker) GetServices(ctx context.Context, projectID string) ([]*model.Service, error) {
	args := filters.Arg("name", fmt.Sprintf("space-cloud-%s", projectID))
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
		taskID := containerName[4]
		service.Version = containerName[2]
		service.ID = containerName[1]
		service.Name = service.ID

		service.ProjectID = projectID
		service.Whitelist = []model.Whitelist{{ProjectID: projectID, Service: "*"}}
		service.Upstreams = []model.Upstream{{ProjectID: projectID, Service: "*"}}
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

		secrets := []string{}
		if err := json.Unmarshal([]byte(service.Labels["internalSecrets"]), &secrets); err != nil {
			logrus.Errorf("error getting service in docker unable to unmarshal secrets - %v", err)
			return nil, err
		}
		dockerSecrets := service.Labels["internalDockerSecrets"]

		imagePullPolicy := service.Labels["internalPullPolicy"]
		if imagePullPolicy == "" {
			imagePullPolicy = string(model.PullIfNotExists)
		}

		delete(service.Labels, "internalSecrets")
		delete(service.Labels, "internalDockerSecrets")
		delete(service.Labels, "internalRuntime")
		delete(service.Labels, "internalPorts")
		delete(service.Labels, "internalProjectId")
		delete(service.Labels, "internalServiceId")
		delete(service.Labels, "internalScale")
		delete(service.Labels, "internalWhitelist")
		delete(service.Labels, "internalAffinity")
		delete(service.Labels, "internalUpstream")
		delete(service.Labels, "internalPullPolicy")

		// set environment variable of task
		envs := map[string]string{}
		for _, value := range containerInspect.Config.Env {
			env := strings.Split(value, "=")
			envs[env[0]] = env[1]
		}

		for _, secret := range secrets {
			// check if file exists
			filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, service.ProjectID, secret)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("secret cannot be read - %v", err)
			}

			// file already exists read it's content
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			fileContent := new(model.Secret)
			if err := json.Unmarshal(data, fileContent); err != nil {
				return nil, err
			}

			if fileContent.Type == "env" {
				for key := range fileContent.Data {
					delete(envs, key)
				}
			}
		}

		if runtime == model.Code {
			delete(envs, model.ArtifactURL)
			delete(envs, model.ArtifactToken)
			delete(envs, model.ArtifactProject)
			delete(envs, model.ArtifactService)
			delete(envs, model.ArtifactVersion)
		}

		tasks = append(tasks, model.Task{
			ID:      taskID,
			Name:    taskID,
			Secrets: secrets,
			Docker: model.Docker{
				Image:           containerInspect.Config.Image,
				Cmd:             containerInspect.Config.Cmd,
				Secret:          dockerSecrets,
				ImagePullPolicy: model.ImagePullPolicy(imagePullPolicy),
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

// AdjustScale adjust the scale for docker instance
func (d *Docker) AdjustScale(_ context.Context, service *model.Service, activeReqs int32) error {
	logrus.Debug("adjust scale not implemented for docker")
	return nil
}

// WaitForService waits for the docker service
func (d *Docker) WaitForService(_ context.Context, service *model.Service) error {
	logrus.Debug("wait for service not implemented for docker")
	return nil
}

// Type returns the docker type of model
func (d *Docker) Type() model.DriverType {
	return model.TypeDocker
}
