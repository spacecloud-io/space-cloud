package istio

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/segmentio/ksuid"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	securityv1beta1 "istio.io/api/security/v1beta1"
	v1beta12 "istio.io/api/type/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func (i *Istio) prepareContainers(service *model.Service, token string, listOfSecrets map[string]*v1.Secret) ([]v1.Container, []v1.Volume, []v1.LocalObjectReference) {
	// There will be n + 1 containers in the pod. Each task will have it's own container. Along with that,
	// there will be a metric collection container as well which pushes metric data to the autoscaler.
	tasks := service.Tasks
	containers := make([]v1.Container, len(tasks))
	volume := map[string]v1.Volume{}
	imagePull := map[string]v1.LocalObjectReference{}

	for j, task := range tasks {
		volumeMount := make([]v1.VolumeMount, 0)

		// Prepare env variables
		var envVars []v1.EnvVar
		for k, v := range task.Env {
			envVars = append(envVars, v1.EnvVar{Name: k, Value: v})
		}
		// Add an environment variable to hold the runtime value
		envVars = append(envVars, v1.EnvVar{Name: runtimeEnvVariable, Value: string(task.Runtime)})
		if task.Runtime == model.Code {
			artifactURL := v1.EnvVar{Name: model.ArtifactURL, Value: i.config.ArtifactAddr}
			artifactToken := v1.EnvVar{Name: model.ArtifactToken, Value: token}
			artifactProject := v1.EnvVar{Name: model.ArtifactProject, Value: service.ProjectID}
			artifactService := v1.EnvVar{Name: model.ArtifactService, Value: service.ID}
			artifactVersion := v1.EnvVar{Name: model.ArtifactVersion, Value: service.Version}
			envVars = append(envVars, artifactURL, artifactToken, artifactProject, artifactService, artifactVersion)
		}

		// Prepare ports to be exposed
		ports := prepareContainerPorts(task.Ports)

		// Prepare command and args
		var cmd, args []string
		if task.Docker.Cmd != nil && len(task.Docker.Cmd) > 0 {
			cmd = task.Docker.Cmd[0:1]
			if len(task.Docker.Cmd) > 1 {
				args = task.Docker.Cmd[1:]
			}
		}

		for _, secretName := range task.Secrets {
			// Check type of Secret..based on that..append :O
			mySecret := listOfSecrets[secretName]
			switch mySecret.ObjectMeta.Annotations["secretType"] {
			case "file":
				// append to VolumeMount
				volumeMount = append(volumeMount, v1.VolumeMount{
					Name:      mySecret.ObjectMeta.Name,
					MountPath: mySecret.ObjectMeta.Annotations["rootPath"],
					ReadOnly:  true})
				volume[secretName] = v1.Volume{
					Name: mySecret.ObjectMeta.Name,
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: mySecret.ObjectMeta.Name,
						},
					},
				}

			case "env":
				// append to env
				for k := range mySecret.Data {
					envVars = append(envVars, v1.EnvVar{
						Name: k,
						ValueFrom: &v1.EnvVarSource{
							SecretKeyRef: &v1.SecretKeySelector{
								LocalObjectReference: v1.LocalObjectReference{
									Name: mySecret.ObjectMeta.Name,
								},
								Key: k,
							},
						},
					})
				}
			}
		}

		imagePull[task.Docker.Secret] = v1.LocalObjectReference{Name: task.Docker.Secret}

		pullPolicy := v1.PullAlways
		if task.Docker.ImagePullPolicy == model.PullIfNotExists {
			pullPolicy = v1.PullIfNotPresent
		}

		containers[j] = v1.Container{
			Name: task.ID,
			Env:  envVars,
			// Resource Related
			Ports:     ports,
			Resources: *generateResourceRequirements(&task.Resources),
			// Docker related
			Image:           task.Docker.Image,
			Command:         cmd,
			Args:            args,
			ImagePullPolicy: pullPolicy,
			// Secrets Related
			VolumeMounts: volumeMount,
		}
	}

	// Add metric proxy container service is purely http based
	var isTCP bool
	for _, task := range tasks {
		for _, port := range task.Ports {
			if port.Protocol == model.TCP {
				isTCP = true
				break
			}
		}
	}
	if !isTCP {
		token, _ := i.auth.SignProxyToken(ksuid.New().String(), service.ProjectID, service.ID, service.Version)
		containers = append(containers, v1.Container{
			Name: "metric-proxy",
			Env:  []v1.EnvVar{{Name: "TOKEN", Value: token}, {Name: "MODE", Value: service.Scale.Mode}},

			// Resource Related
			Resources: *generateResourceRequirements(&model.Resources{CPU: 20, Memory: 50}),

			// Docker related
			Image:           "spaceuptech/metric-proxy:0.2.0",
			Command:         []string{"./app"},
			Args:            []string{"start"},
			ImagePullPolicy: v1.PullIfNotPresent,
		})
	}

	// Convert map to array
	arrVolume := make([]v1.Volume, 0)
	for _, v := range volume {
		arrVolume = append(arrVolume, v)
	}
	arrImagePull := make([]v1.LocalObjectReference, 0)
	for _, v := range imagePull {
		arrImagePull = append(arrImagePull, v)
	}
	return containers, arrVolume, arrImagePull
}

func prepareContainerPorts(taskPorts []model.Port) []v1.ContainerPort {
	ports := make([]v1.ContainerPort, len(taskPorts))
	for i, p := range taskPorts {
		ports[i] = v1.ContainerPort{Name: fmt.Sprintf("%s-%s", p.Name, p.Protocol), ContainerPort: p.Port}
	}

	return ports
}

func prepareServicePorts(tasks []model.Task) []v1.ServicePort {
	var ports []v1.ServicePort
	for _, task := range tasks {
		for _, p := range task.Ports {
			ports = append(ports, v1.ServicePort{Name: p.Name, Port: p.Port})
		}
	}

	return ports
}

func prepareVirtualServiceHTTPRoutes(projectID, serviceID string, services map[string]model.ScaleConfig, routes model.Routes, proxyPort uint32) ([]*networkingv1alpha3.HTTPRoute, error) {
	var httpRoutes []*networkingv1alpha3.HTTPRoute

	for _, route := range routes {
		// Check if the port provided is correct
		if route.Source.Port == 0 {
			return nil, errors.New("port cannot be zero")
		}

		// Check if at least one target is provided
		if len(route.Targets) == 0 {
			return nil, errors.New("at least one target needs to be provided")
		}

		// Prepare an array of targets / destinations
		var destinations []*networkingv1alpha3.HTTPRouteDestination
		for _, target := range route.Targets {
			switch target.Type {
			case model.RouteTargetVersion:
				// Check if config for version exists
				versionScaleConfig, p := services[target.Version]
				if !p {
					return nil, fmt.Errorf("version (%s) not found for service (%s)", target.Version, serviceID)
				}

				// Prepare variables
				destHost := getInternalServiceDomain(projectID, serviceID, target.Version)
				destPort := uint32(target.Port)

				// Redirect traffic to runner when no of replicas is equal to zero. The runner proxy will scale up the service to service incoming requests.
				if versionScaleConfig.MinReplicas == 0 {
					destHost = "runner.space-cloud.svc.cluster.local"
					destPort = proxyPort
				}

				destinations = append(destinations, &networkingv1alpha3.HTTPRouteDestination{
					// We will always set the headers since it helps us with the routing rules. Also, we check if the headers is present to determine
					// if the destination is version or external.
					Headers: &networkingv1alpha3.Headers{
						Request: &networkingv1alpha3.Headers_HeaderOperations{
							Set: map[string]string{
								"x-og-project": projectID,
								"x-og-service": serviceID,
								"x-og-host":    getInternalServiceDomain(projectID, serviceID, target.Version),
								"x-og-port":    strconv.Itoa(int(target.Port)),
								"x-og-version": target.Version,
							},
						},
					},
					Destination: &networkingv1alpha3.Destination{
						Host: destHost,
						Port: &networkingv1alpha3.PortSelector{Number: destPort},
					},
					Weight: target.Weight,
				})

			case model.RouteTargetExternal:
				destinations = append(destinations, &networkingv1alpha3.HTTPRouteDestination{
					Destination: &networkingv1alpha3.Destination{
						Host: target.Host,
						Port: &networkingv1alpha3.PortSelector{Number: uint32(target.Port)},
					},
					Weight: target.Weight,
				})
			default:
				return nil, fmt.Errorf("invalid target type (%s) provided", target.Type)
			}
		}

		// Add the http route
		httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
			Name:    fmt.Sprintf("http-%d", route.Source.Port),
			Match:   []*networkingv1alpha3.HTTPMatchRequest{{Port: uint32(route.Source.Port), Gateways: []string{"mesh"}}},
			Retries: &networkingv1alpha3.HTTPRetry{Attempts: 3, PerTryTimeout: &types.Duration{Seconds: 180}},
			Route:   destinations,
		})
	}

	return httpRoutes, nil
}

func updateOrCreateVirtualServiceRoutes(service *model.Service, proxyPort uint32, prevVirtualService *v1alpha3.VirtualService) ([]*networkingv1alpha3.HTTPRoute, []*networkingv1alpha3.TCPRoute) {
	// Update the existing destinations of this version if virtual service already exist. We only need to do this for http services.
	if prevVirtualService != nil {
		for _, httpRoute := range prevVirtualService.Spec.Http {
			for _, dest := range httpRoute.Route {

				// Check if the route was for a service with min scale 0. If the destination has the host of runner, it means it is communicating via the proxy.
				if dest.Destination.Host == "runner.space-cloud.svc.cluster.local" {
					// We are only interested in this case if the new min replica for this version is more than 0. If the min replica was zero there would be no change
					if service.Scale.MinReplicas == 0 {
						continue
					}

					// Update this particular destination if the version matches with ours. We need to make the communication `direct`
					if service.Version == dest.Headers.Request.Set["x-og-version"] {
						// Set the destination host
						dest.Destination.Host = getInternalServiceDomain(service.ProjectID, service.ID, service.Version)

						// Set the destination port
						port, _ := strconv.Atoi(dest.Headers.Request.Set["x-og-port"])
						dest.Destination.Port = &networkingv1alpha3.PortSelector{Number: uint32(port)}
					}
				}

				// Since we are here it means the given destination communicated with the target directly. We don't really care if the min replica is greater
				// than zero because this would mean there is no change.
				if service.Scale.MinReplicas > 0 {
					continue
				}

				// Update the destination to communicate via the proxy if its for our version
				if dest.Destination.Host == getInternalServiceDomain(service.ProjectID, service.ID, service.Version) {
					dest.Destination.Host = "runner.space-cloud.svc.cluster.local"
					dest.Destination.Port = &networkingv1alpha3.PortSelector{Number: proxyPort}
				}
			}
		}
		return prevVirtualService.Spec.Http, prevVirtualService.Spec.Tcp
	}

	// Reaching here means we have to create new rules
	var httpRoutes []*networkingv1alpha3.HTTPRoute
	var tcpRoutes []*networkingv1alpha3.TCPRoute

	for i, task := range service.Tasks {
		for j, port := range task.Ports {
			switch port.Protocol {
			case model.HTTP:
				// Prepare variables
				destHost := getInternalServiceDomain(service.ProjectID, service.ID, service.Version)
				destPort := uint32(port.Port)

				// Redirect traffic to runner when no of replicas is equal to zero. The runner proxy will scale up the service to service incoming requests.
				if service.Scale.MinReplicas == 0 {
					destHost = "runner.space-cloud.svc.cluster.local"
					destPort = proxyPort
				}

				httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
					Name:    fmt.Sprintf("http-%d%d-%d", j, i, port.Port),
					Match:   []*networkingv1alpha3.HTTPMatchRequest{{Port: uint32(port.Port), Gateways: []string{"mesh"}}},
					Retries: &networkingv1alpha3.HTTPRetry{Attempts: 3, PerTryTimeout: &types.Duration{Seconds: 180}},
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						{
							// We will always set the headers since it helps us with the routing rules
							Headers: &networkingv1alpha3.Headers{
								Request: &networkingv1alpha3.Headers_HeaderOperations{
									Set: map[string]string{
										"x-og-project": service.ProjectID,
										"x-og-service": service.ID,
										"x-og-host":    getInternalServiceDomain(service.ProjectID, service.ID, service.Version),
										"x-og-port":    strconv.Itoa(int(port.Port)),
										"x-og-version": service.Version,
									},
								},
							},
							Destination: &networkingv1alpha3.Destination{
								Host: destHost,
								Port: &networkingv1alpha3.PortSelector{Number: destPort},
							},
							Weight: 100,
						},
					},
				})

			case model.TCP:
				tcpRoutes = append(tcpRoutes, &networkingv1alpha3.TCPRoute{
					Match: []*networkingv1alpha3.L4MatchAttributes{{Port: uint32(port.Port)}},
					Route: []*networkingv1alpha3.RouteDestination{
						{
							Destination: &networkingv1alpha3.Destination{
								Host: getInternalServiceDomain(service.ProjectID, service.ID, service.Version),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(port.Port)},
							},
						},
					},
				})
			}
		}
	}

	return httpRoutes, tcpRoutes
}

func prepareVirtualServiceHosts(service *model.Service) []string {
	hosts := []string{fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID)}
	return hosts
}

func prepareAuthPolicyRules(service *model.Service) []*securityv1beta1.Rule {
	var froms []*securityv1beta1.Rule_From
	var namespaces []string
	var principals []string

	service.Whitelist = append(service.Whitelist, model.Whitelist{ProjectID: "space-cloud", Service: "*"}, model.Whitelist{ProjectID: "istio-system", Service: "*"})

	for _, whitelist := range service.Whitelist {

		if whitelist.ProjectID == "*" {
			// This means this is an open service
			return []*securityv1beta1.Rule{{}}
		}

		if whitelist.Service == "*" {
			// This means that the service can be accessed from everyone in the project
			namespaces = append(namespaces, whitelist.ProjectID)
		} else {
			// This means that the service can be accessed only from that service in that project
			principals = append(principals, fmt.Sprintf("cluster.local/ns/%s/sa/%s", whitelist.ProjectID, whitelist.Service))
		}
	}

	if namespaces != nil {
		froms = append(froms, &securityv1beta1.Rule_From{
			Source: &securityv1beta1.Source{Namespaces: namespaces},
		})
	}
	if principals != nil {
		froms = append(froms, &securityv1beta1.Rule_From{
			Source: &securityv1beta1.Source{Principals: principals},
		})
	}

	return []*securityv1beta1.Rule{{From: froms}}
}

func prepareUpstreamHosts(service *model.Service) []string {
	hosts := make([]string, len(service.Upstreams)+1)

	// First entry will always be space-cloud
	hosts[0] = "space-cloud/*"

	for i, upstream := range service.Upstreams {
		projectID := upstream.ProjectID
		serviceID := upstream.Service
		if serviceID != "*" {
			serviceID = getServiceDomainName(projectID, serviceID)
		}
		hosts[i+1] = upstream.ProjectID + "/" + serviceID
	}

	return hosts
}

func generateServiceAccount(service *model.Service) *v1.ServiceAccount {
	return &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
		Name:        getServiceAccountName(service.ID),
		Labels:      map[string]string{"account": service.ID},
		Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
	}}
}

func (i *Istio) generateDeployment(service *model.Service, token string, listOfSecrets map[string]*v1.Secret) *appsv1.Deployment {
	preparedContainer, volumes, imagePull := i.prepareContainers(service, token, listOfSecrets)
	// Make sure the desired replica count doesn't cross the min and max range
	if service.Scale.Replicas < service.Scale.MinReplicas {
		service.Scale.Replicas = service.Scale.MinReplicas
	}
	if service.Scale.Replicas > service.Scale.MaxReplicas {
		service.Scale.Replicas = service.Scale.MaxReplicas
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: getDeploymentName(service.ID, service.Version),
			Labels: map[string]string{
				"app":     service.ID,
				"version": service.Version,
			},
			Annotations: map[string]string{
				"concurrency": strconv.Itoa(int(service.Scale.Concurrency)),
				"minReplicas": strconv.Itoa(int(service.Scale.MinReplicas)),
				"maxReplicas": strconv.Itoa(int(service.Scale.MaxReplicas)),
				"mode":        service.Scale.Mode,
				"generatedBy": getGeneratedByAnnotationName(),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &service.Scale.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": service.ID, "version": service.Version}},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"sidecar.istio.io/statsInclusionPrefixes": "cluster.outbound,listener,http.inbound,cluster_manager,listener_manager,http_mixer_filter,tcp_mixer_filter,server,cluster.xds-grpc"},
					Labels:      map[string]string{"app": service.ID, "version": service.Version},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: getServiceAccountName(service.ID),
					Containers:         preparedContainer,
					Volumes:            volumes,
					ImagePullSecrets:   imagePull,
					// TODO: Add config for affinity
				},
			},
		},
	}
}

func generateGeneralService(service *model.Service) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getServiceName(service.ID),
			Labels:      map[string]string{"app": service.ID, "service": service.ID},
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
		},
		Spec: v1.ServiceSpec{
			Ports:    prepareServicePorts(service.Tasks),
			Selector: map[string]string{"app": service.ID},
			Type:     v1.ServiceTypeClusterIP,
		},
	}
}

func generateInternalService(service *model.Service) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getInternalServiceName(service.ID, service.Version),
			Labels:      map[string]string{"app": service.ID, "service": service.ID, "version": service.Version},
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
		},
		Spec: v1.ServiceSpec{
			Ports:    prepareServicePorts(service.Tasks),
			Selector: map[string]string{"app": service.ID, "version": service.Version},
			Type:     v1.ServiceTypeClusterIP,
		},
	}
}

func (i *Istio) updateVirtualService(service *model.Service, prevVirtualService *v1alpha3.VirtualService) *v1alpha3.VirtualService {
	httpRoutes, tcpRoutes := updateOrCreateVirtualServiceRoutes(service, i.config.ProxyPort, prevVirtualService)
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getVirtualServiceName(service.ID),
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
			Labels:      map[string]string{"app": service.ID}, // We use the app label to retrieve service routing rules
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: prepareVirtualServiceHosts(service),
			// Gateways: prepareVirtualServiceGateways(service),
			Http: httpRoutes,
			Tcp:  tcpRoutes,
		},
	}
}
func (i *Istio) generateVirtualServiceBasedOnRoutes(projectID, serviceID string, scaleConfig map[string]model.ScaleConfig, routes model.Routes, prevVirtualService *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error) {
	// Generate the httpRoutes based on the routes provided
	httpRoutes, err := prepareVirtualServiceHTTPRoutes(projectID, serviceID, scaleConfig, routes, i.config.ProxyPort)
	if err != nil {
		return nil, err
	}

	// Create a prevVirtualService if its nil to prevent panic
	if prevVirtualService == nil {
		prevVirtualService = new(v1alpha3.VirtualService)
	}

	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getVirtualServiceName(serviceID),
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
			Labels:      map[string]string{"app": serviceID}, // We use the app label to retrieve service routing rules
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: []string{getServiceDomainName(projectID, serviceID)},
			Http:  httpRoutes,
			Tcp:   prevVirtualService.Spec.Tcp,
		},
	}, nil
}

func generateGeneralDestinationRule(service *model.Service) *v1alpha3.DestinationRule {
	return &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getGeneralDestRuleName(service.ID),
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID),
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				Tls: &networkingv1alpha3.ClientTLSSettings{Mode: networkingv1alpha3.ClientTLSSettings_ISTIO_MUTUAL},
			},
		},
	}
}

func generateInternalDestinationRule(service *model.Service) *v1alpha3.DestinationRule {
	return &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getInternalDestRuleName(service.ID, service.Version),
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: getInternalServiceDomain(service.ProjectID, service.ID, service.Version),
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				Tls: &networkingv1alpha3.ClientTLSSettings{Mode: networkingv1alpha3.ClientTLSSettings_ISTIO_MUTUAL},
			},
		},
	}
}

func generateAuthPolicy(service *model.Service) *v1beta1.AuthorizationPolicy {
	authPolicy := &v1beta1.AuthorizationPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getAuthorizationPolicyName(service.ProjectID, service.ID, service.Version),
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
		},
		Spec: securityv1beta1.AuthorizationPolicy{
			Selector: &v1beta12.WorkloadSelector{MatchLabels: map[string]string{"app": service.ID, "version": service.Version}},
			Rules:    prepareAuthPolicyRules(service),
		},
	}
	return authPolicy
}

func generateSidecarConfig(service *model.Service) *v1alpha3.Sidecar {
	return &v1alpha3.Sidecar{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getSidecarName(service.ID, service.Version),
			Annotations: map[string]string{"generatedBy": getGeneratedByAnnotationName()},
		},
		Spec: networkingv1alpha3.Sidecar{
			WorkloadSelector:      &networkingv1alpha3.WorkloadSelector{Labels: map[string]string{"app": service.ID, "version": service.Version}},
			Egress:                []*networkingv1alpha3.IstioEgressListener{{Hosts: prepareUpstreamHosts(service)}},
			OutboundTrafficPolicy: &networkingv1alpha3.OutboundTrafficPolicy{Mode: networkingv1alpha3.OutboundTrafficPolicy_ALLOW_ANY},
		},
	}
}

// We don't need gateways anymore since space cloud gateway will take care of it
// func generateGateways(service *model.Service) *v1alpha3.Gateway {
// 	hosts := make([]string, 0)
//
// 	// Add the hosts if provided
// 	if service.Expose != nil && service.Expose.Hosts != nil {
// 		hosts = append(hosts, service.Expose.Hosts...)
// 	}
//
// 	// Add dummy host if none exist
// 	hosts = append(hosts, "dummy.com")
//
// 	return &v1alpha3.Gateway{
// 		ObjectMeta: metav1.ObjectMeta{ID: getGatewayName(service)},
// 		Spec: networkingv1alpha3.Gateway{
// 			Selector: map[string]string{"istio": "ingressgateway"},
// 			Servers: []*networkingv1alpha3.Server{
// 				{
// 					// TODO: make this https and load certificates dynamically
// 					Port: &networkingv1alpha3.Port{
// 						Number:   80,
// 						ID:     "http",
// 						Protocol: "HTTP",
// 					},
// 					Hosts: hosts,
// 				},
// 			},
// 		},
// 	}
// }

func generateResourceRequirements(c *model.Resources) *v1.ResourceRequirements {
	// Set default values if either value is absent
	if c.Memory == 0 || c.CPU == 0 {
		c.Memory = 512
		c.CPU = 250
	}

	resources := v1.ResourceRequirements{Limits: v1.ResourceList{}, Requests: v1.ResourceList{}}

	// Set the cpu limits
	// resources.Limits[v1.ResourceCPU] = *resource.NewMilliQuantity(c.CPU, resource.DecimalSI)
	resources.Requests[v1.ResourceCPU] = *resource.NewMilliQuantity(c.CPU, resource.DecimalSI)

	// Set the memory limits
	// resources.Limits[v1.ResourceMemory] = *resource.NewQuantity(c.Memory*1024*1024, resource.BinarySI)
	resources.Requests[v1.ResourceMemory] = *resource.NewQuantity(c.Memory*1024*1024, resource.BinarySI)

	// Set the GPU limits
	if c.GPU != nil && c.GPU.Value > 0 && c.GPU.Type != "" {
		resources.Limits[v1.ResourceName(fmt.Sprintf("%s.com/gpu", c.GPU.Type))] = *resource.NewQuantity(c.GPU.Value, resource.DecimalSI)
	}
	return &resources
}

func adjustMinScale(service *model.Service) {
	// Simply return if min replicas is greater than zero
	if service.Scale.MinReplicas > 0 {
		return
	}

	for _, task := range service.Tasks {
		for _, port := range task.Ports {
			if port.Protocol == model.TCP {
				service.Scale.MinReplicas = 1
				break
			}
		}
	}

}
