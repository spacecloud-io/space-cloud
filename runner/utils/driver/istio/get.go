package istio

import (
	"context"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// GetServices gets the services for istio
func (i *Istio) GetServices(ctx context.Context, projectID string) ([]*model.Service, error) {
	deploymentList, err := i.kube.AppsV1().Deployments(projectID).List(ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Error getting service in istio - unable to find deployment - %v", err)
		return nil, err
	}
	services := []*model.Service{}
	for _, deployment := range deploymentList.Items {
		service := new(model.Service)
		service.ProjectID = projectID
		service.ID = deployment.Labels["app"]
		service.Version = deployment.Labels["version"]

		// Get scale config
		scale, err := getScaleConfigFromDeployment(deployment)
		if err != nil {
			return nil, err
		}
		service.Scale = scale

		for _, containerInfo := range deployment.Spec.Template.Spec.Containers {
			if containerInfo.Name == "metric-proxy" || containerInfo.Name == "istio-proxy" {
				continue
			}
			// get ports
			ports := make([]model.Port, len(containerInfo.Ports))
			for i, port := range containerInfo.Ports {
				array := strings.Split(port.Name, "-")
				ports[i] = model.Port{Name: array[0], Protocol: model.Protocol(array[1]), Port: port.ContainerPort}
			}

			var dockerSecret string
			secretsMap := make(map[string]struct{})

			// get environment variables
			envs := map[string]string{}
			for _, env := range containerInfo.Env {
				if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
					secretsMap[env.ValueFrom.SecretKeyRef.LocalObjectReference.Name] = struct{}{}
					continue
				}
				envs[env.Name] = env.Value
			}

			// Range over the file mounts for secrets
			for _, volume := range containerInfo.VolumeMounts {
				if checkIfVolumeIsSecret(volume.Name, deployment.Spec.Template.Spec.Volumes) {
					secretsMap[volume.Name] = struct{}{}
				}
			}

			// Get docker secret
			// TODO: Handle case when different tasks have different secrets
			if len(deployment.Spec.Template.Spec.ImagePullSecrets) > 0 {
				dockerSecret = deployment.Spec.Template.Spec.ImagePullSecrets[0].Name
			}

			// Extract the runtime from the environment variable
			runtime := model.Runtime(envs[runtimeEnvVariable])
			delete(envs, runtimeEnvVariable)

			// Delete internal environment variables if runtime was code
			if runtime == model.Code {
				delete(envs, model.ArtifactURL)
				delete(envs, model.ArtifactToken)
				delete(envs, model.ArtifactProject)
				delete(envs, model.ArtifactService)
				delete(envs, model.ArtifactVersion)
			}

			// Get the image pull policy
			imagePullPolicy := model.PullIfNotExists
			if containerInfo.ImagePullPolicy == v1.PullAlways {
				imagePullPolicy = model.PullAlways
			}

			// Move all secrets from map to array
			var secrets []string
			for k := range secretsMap {
				secrets = append(secrets, k)
			}

			// set tasks
			service.Tasks = append(service.Tasks, model.Task{
				ID:    containerInfo.Name,
				Name:  containerInfo.Name,
				Ports: ports,
				Resources: model.Resources{
					CPU:    containerInfo.Resources.Requests.Cpu().MilliValue(),
					Memory: containerInfo.Resources.Requests.Memory().Value() / (1024 * 1024),
				},
				Docker: model.Docker{
					Image:           containerInfo.Image,
					Cmd:             containerInfo.Command,
					Secret:          dockerSecret,
					ImagePullPolicy: imagePullPolicy,
				},
				Env:     envs,
				Runtime: runtime,
				Secrets: secrets,
			})
		}

		// set whitelist
		authPolicy, _ := i.istio.SecurityV1beta1().AuthorizationPolicies(projectID).Get(ctx, getAuthorizationPolicyName(service.ProjectID, service.ID, service.Version), metav1.GetOptions{})
		if len(authPolicy.Spec.Rules[0].From) != 0 {
			for _, rule := range authPolicy.Spec.Rules[0].From {
				for _, projectID := range rule.Source.Namespaces {
					if projectID == "space-cloud" || projectID == "istio-system" {
						continue
					}
					service.Whitelist = append(service.Whitelist, model.Whitelist{ProjectID: projectID, Service: "*"})
				}
				for _, serv := range rule.Source.Principals {
					whitelistArr := strings.Split(serv, "/")
					if len(whitelistArr) != 5 {
						logrus.Error("error getting service in istio length of whitelist array is not equal to 5")
						continue
					}
					service.Whitelist = append(service.Whitelist, model.Whitelist{ProjectID: whitelistArr[2], Service: whitelistArr[4]})
				}
			}
		}

		// Set upstreams
		sideCar, err := i.istio.NetworkingV1alpha3().Sidecars(projectID).Get(ctx, getSidecarName(service.ID, service.Version), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		for _, value := range sideCar.Spec.Egress[0].Hosts {
			a := strings.Split(value, "/")
			if a[0] == "space-cloud" || a[0] == "istio-system" {
				continue
			}
			service.Upstreams = append(service.Upstreams, model.Upstream{ProjectID: a[0], Service: a[1]})
		}

		// todo labels, serviceName, affinity, runtime
		services = append(services, service)
	}

	return services, nil
}

// GetServiceRoutes gets the routing rules of each service
func (i *Istio) GetServiceRoutes(ctx context.Context, projectID string) (map[string]model.Routes, error) {
	ns := projectID

	// Get all virtual services
	services, err := i.getVirtualServices(ctx, ns)
	if err != nil {
		return nil, err
	}

	serviceRoutes := make(map[string]model.Routes, len(services.Items))

	for _, service := range services.Items {
		serviceID := service.Labels["app"]
		routes := make(model.Routes, len(service.Spec.Http))

		for i, route := range service.Spec.Http {

			// Generate the targets
			targets := make([]model.RouteTarget, len(route.Route))
			for j, destination := range route.Route {
				target := model.RouteTarget{Weight: destination.Weight}

				// Figure out the route type
				target.Type = model.RouteTargetExternal
				if destination.Headers != nil {
					target.Type = model.RouteTargetVersion
				}
				switch target.Type {
				case model.RouteTargetVersion:
					// Set the version field if target type was version
					target.Version = destination.Headers.Request.Set["x-og-version"]

					// Set the port
					port, err := strconv.Atoi(destination.Headers.Request.Set["x-og-port"])
					if err != nil {
						return nil, err
					}
					target.Port = int32(port)

				case model.RouteTargetExternal:
					// Set the host field if target type was external
					target.Host = destination.Destination.Host

					// Set the port
					target.Port = int32(destination.Destination.Port.Number)
				}

				targets[j] = target
			}

			// Set the route
			routes[i] = &model.Route{ID: serviceID, Source: model.RouteSource{Port: int32(route.Match[0].Port)}, Targets: targets}
		}

		// Set the routes of a service
		serviceRoutes[serviceID] = routes
	}

	return serviceRoutes, nil
}
