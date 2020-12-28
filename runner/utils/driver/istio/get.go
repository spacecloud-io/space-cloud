package istio

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// GetServices gets the services for istio
func (i *Istio) GetServices(ctx context.Context, projectID string) ([]*model.Service, error) {
	// Get all deployments in project
	deploymentList, err := i.kube.AppsV1().Deployments(projectID).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to find deployments in project", err, nil)
	}

	// Get all keda trigger authentication in project
	triggerAuthList, err := i.keda.KedaV1alpha1().TriggerAuthentications(projectID).List(ctx, metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=space-cloud"})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to find keda trigger auths in project", err, nil)
	}

	// Get all the keda scaled objects in projects
	scaledObjectList, err := i.keda.KedaV1alpha1().ScaledObjects(projectID).List(ctx, metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=space-cloud"})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to find keda scaled object in project", err, nil)
	}

	services := []*model.Service{}
	for _, deployment := range deploymentList.Items {
		service := new(model.Service)
		service.ProjectID = projectID
		service.ID = deployment.Labels["app"]
		service.Version = deployment.Labels["version"]
		service.Affinity = make([]model.Affinity, 0)
		service.StatsInclusionPrefixes = deployment.Spec.Template.Annotations["sidecar.istio.io/statsInclusionPrefixes"]

		// Extract affinities
		if deployment.Spec.Template.Spec.Affinity != nil {

			// node affinity preferred
			if deployment.Spec.Template.Spec.Affinity.NodeAffinity != nil {

				if deployment.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
					affinities := extractPreferredNodeAffinityObject(deployment.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution)
					if len(affinities) > 0 {
						service.Affinity = append(service.Affinity, affinities...)
					}
				}

				// node affinity required
				if deployment.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
					affinities := extractRequiredNodeAffinityObject(deployment.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms)
					if len(affinities) > 0 {
						service.Affinity = append(service.Affinity, affinities...)
					}
				}
			}

			// service affinity
			if deployment.Spec.Template.Spec.Affinity.PodAffinity != nil {
				affinities := extractPreferredServiceAffinityObject(deployment.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution, 1)
				if len(affinities) > 0 {
					service.Affinity = append(service.Affinity, affinities...)
				}
				affinities = extractRequiredServiceAffinityObject(deployment.Spec.Template.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution, 1)
				if len(affinities) > 0 {
					service.Affinity = append(service.Affinity, affinities...)
				}
			}

			// service anti affinity
			if deployment.Spec.Template.Spec.Affinity.PodAntiAffinity != nil {
				affinities := extractPreferredServiceAffinityObject(deployment.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution, -1)
				if len(affinities) > 0 {
					service.Affinity = append(service.Affinity, affinities...)
				}
				affinities = extractRequiredServiceAffinityObject(deployment.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, -1)
				if len(affinities) > 0 {
					service.Affinity = append(service.Affinity, affinities...)
				}
			}
		}

		// service labels
		service.Labels = deployment.Spec.Template.Labels

		// Get scale config
		service.AutoScale = getScaleConfigFromKedaConfig(service.ID, service.Version, scaledObjectList.Items, triggerAuthList.Items)
		if service.AutoScale == nil {
			service.AutoScale = getScaleConfigFromDeployment(deployment)
		}

		for _, containerInfo := range deployment.Spec.Template.Spec.Containers {
			if containerInfo.Name == "metric-proxy" || containerInfo.Name == "istio-proxy" {
				continue
			}
			// get ports
			ports := make([]model.Port, len(containerInfo.Ports))
			for i, port := range containerInfo.Ports {
				proto := strings.Split(port.Name, "-")[0]
				ports[i] = model.Port{Name: port.Name, Protocol: model.Protocol(proto), Port: port.ContainerPort}
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
			// if runtime == model.Code {
			// 	delete(envs, model.ArtifactURL)
			// 	delete(envs, model.ArtifactToken)
			// 	delete(envs, model.ArtifactProject)
			// 	delete(envs, model.ArtifactService)
			// 	delete(envs, model.ArtifactVersion)
			// }

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
		authPolicy, err := i.istio.SecurityV1beta1().AuthorizationPolicies(projectID).Get(ctx, getAuthorizationPolicyName(service.ProjectID, service.ID, service.Version), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
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
						_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error getting service in istio length of whitelist array is not equal to 5", nil, nil)
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

		// todo serviceName, runtime
		services = append(services, service)
	}

	return services, nil
}

// GetServiceStatus gets the services status for istio
func (i *Istio) GetServiceStatus(ctx context.Context, projectID string) ([]*model.ServiceStatus, error) {
	deploymentList, err := i.kube.AppsV1().Deployments(projectID).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Error getting service in istio - unable to find deployment", err, nil)
	}
	result := make([]*model.ServiceStatus, 0)
	for _, deployment := range deploymentList.Items {
		serviceID := deployment.Labels["app.kubernetes.io/name"]
		serviceVersion := deployment.Labels["app.kubernetes.io/version"]

		podlist, err := i.kube.CoreV1().Pods(deployment.Namespace).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s,version=%s", serviceID, serviceVersion)})
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Error getting service in istio - unable to find pods", err, nil)
		}
		replicas := make([]*model.ReplicaInfo, 0)
		for _, p := range podlist.Items {
			replicas = append(replicas, &model.ReplicaInfo{ID: p.Name, Status: strings.ToUpper(string(p.Status.Phase))})
		}
		result = append(result, &model.ServiceStatus{
			ServiceID:       serviceID,
			Version:         serviceVersion,
			DesiredReplicas: deployment.Spec.Replicas,
			Replicas:        replicas,
		})
	}
	return result, nil
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
		routes := make(model.Routes, len(service.Spec.Http)+len(service.Spec.Tcp))

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
			routes[i] = &model.Route{ID: serviceID, RequestTimeout: route.Retries.PerTryTimeout.Seconds, RequestRetries: route.Retries.Attempts, Source: model.RouteSource{Port: int32(route.Match[0].Port), Protocol: model.HTTP}, Targets: targets}
		}

		for i, route := range service.Spec.Tcp {

			// Generate the targets
			targets := make([]model.RouteTarget, len(route.Route))
			for j, destination := range route.Route {
				target := model.RouteTarget{Weight: destination.Weight}

				// Figure out the route type
				target.Type = model.RouteTargetExternal
				if checkIfInternalServiceDomain(projectID, serviceID, destination.Destination.Host) {
					target.Type = model.RouteTargetVersion
				}

				switch target.Type {
				case model.RouteTargetVersion:
					// Set the version field if target type was version
					_, _, version := splitInternalServiceDomain(destination.Destination.Host)
					target.Version = version

					// Set the port
					target.Port = int32(destination.Destination.Port.Number)

				case model.RouteTargetExternal:
					// Set the host field if target type was external
					target.Host = destination.Destination.Host

					// Set the port
					target.Port = int32(destination.Destination.Port.Number)
				}

				targets[j] = target
			}

			// Set the route
			routes[i] = &model.Route{ID: serviceID, Source: model.RouteSource{Port: int32(route.Match[0].Port), Protocol: model.TCP}, Targets: targets}
		}

		// Set the routes of a service
		serviceRoutes[serviceID] = routes
	}

	return serviceRoutes, nil
}

// GetServiceRole gets the service role rules of each service
func (i *Istio) GetServiceRole(ctx context.Context, projectID string) ([]*model.Role, error) {
	ns := projectID

	rolelist, err := i.kube.RbacV1().Roles(ns).List(ctx, metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=space-cloud"})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list roles in project (%s)", projectID), err, nil)
	}

	clusterRoleList, err := i.kube.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=space-cloud"})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list cluster roles in project (%s)", projectID), err, nil)
	}
	serviceRole := make([]*model.Role, len(rolelist.Items)+len(clusterRoleList.Items))

	for _, role := range rolelist.Items {
		serviceID := role.Labels["app"]
		Role := new(model.Role)
		Role.ID = role.Name
		Role.Project = role.Namespace
		Role.Service = serviceID
		Role.Type = model.ServiceRoleProject
		Rules := make([]model.Rule, 0)
		for _, rule := range role.Rules {
			Rules = append(Rules, model.Rule{APIGroups: rule.APIGroups, Verbs: rule.Verbs, Resources: rule.Resources})
		}
		Role.Rules = Rules
		serviceRole = append(serviceRole, Role)
	}

	for _, role := range clusterRoleList.Items {
		serviceID := role.Labels["app"]
		Role := new(model.Role)
		Role.ID = role.Name
		Role.Project = projectID
		Role.Service = serviceID
		Role.Type = model.ServiceRoleCluster
		Rules := make([]model.Rule, 0)
		for _, rule := range role.Rules {
			Rules = append(Rules, model.Rule{APIGroups: rule.APIGroups, Verbs: rule.Verbs, Resources: rule.Resources})
		}
		Role.Rules = Rules
		serviceRole = append(serviceRole, Role)
	}

	return serviceRole, nil
}
