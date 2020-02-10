package istio

import (
	"fmt"
	"strconv"
	"strings"

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
	// TODO: Add support for private repos
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
		if task.Docker.Cmd != nil {
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
			ImagePullPolicy: v1.PullAlways, // TODO: make this configurable
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
			Env:  []v1.EnvVar{{Name: "TOKEN", Value: token}},

			// Resource Related
			Resources: *generateResourceRequirements(&model.Resources{CPU: 20, Memory: 50}),

			// Docker related
			Image:           "spaceuptech/metric-proxy:latest", // TODO: Lets use the version tag here to make sure we pull the latest image
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

func makeOriginalVirtualService(service *model.Service, virtualService *v1alpha3.VirtualService) {
	ogHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID)

	// Redo the http routes. The tcp routes are lost anyways so we don't really care about them.
	for _, httpRoute := range virtualService.Spec.Http {
		for _, route := range httpRoute.Route {
			// Revert the destination to original
			port, _ := strconv.Atoi(strings.Split(httpRoute.Name, "-")[2])
			route.Destination.Host = ogHost
			route.Destination.Port.Number = uint32(port)

			// Reset the headers
			route.Headers = nil
		}
	}
}

func makeScaleZeroVirtualService(service *model.Service, virtualService *v1alpha3.VirtualService, proxyPort uint32) {
	ogHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID)

	// Redirect traffic to runner when no of replicas is equal to zero. The runner proxy will scale up the service
	// to service incoming requests.
	for _, httpRoute := range virtualService.Spec.Http {
		for _, route := range httpRoute.Route {
			// Set the destination to runner proxy
			route.Destination.Host = "runner.space-cloud.svc.cluster.local"
			route.Destination.Port.Number = proxyPort

			// Set the headers
			route.Headers = &networkingv1alpha3.Headers{
				Request: &networkingv1alpha3.Headers_HeaderOperations{
					Add: map[string]string{
						"x-og-project": service.ProjectID,
						"x-og-service": service.ID,
						"x-og-host":    ogHost,
						"x-og-port":    strings.Split(httpRoute.Name, "-")[2],
						"x-og-version": service.Version,
					},
				},
			}
		}
	}
}

func prepareVirtualServiceRoutes(service *model.Service, proxyPort uint32) ([]*networkingv1alpha3.HTTPRoute, []*networkingv1alpha3.TCPRoute) {
	var httpRoutes []*networkingv1alpha3.HTTPRoute
	var tcpRoutes []*networkingv1alpha3.TCPRoute

	for i, task := range service.Tasks {
		for j, port := range task.Ports {
			switch port.Protocol {
			case model.HTTP:
				// Prepare variables
				var headers *networkingv1alpha3.Headers
				retries := &networkingv1alpha3.HTTPRetry{Attempts: 3, PerTryTimeout: &types.Duration{Seconds: 90}}
				destHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID)
				destPort := uint32(port.Port)

				// Redirect traffic to runner when no of replicas is equal to zero. The runner proxy will scale up the service
				// to service incoming requests.
				if service.Scale.Replicas == 0 {
					headers = &networkingv1alpha3.Headers{
						Request: &networkingv1alpha3.Headers_HeaderOperations{
							Add: map[string]string{
								"x-og-project": service.ProjectID,
								"x-og-service": service.ID,
								"x-og-host":    destHost,
								"x-og-port":    strconv.Itoa(int(destPort)),
								"x-og-version": service.Version,
							},
						},
					}
					retries = &networkingv1alpha3.HTTPRetry{Attempts: 1, PerTryTimeout: &types.Duration{Seconds: 180}}
					destHost = "runner.space-cloud.svc.cluster.local"
					destPort = proxyPort
				}

				httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
					Name:    fmt.Sprintf("http-%d%d-%d", j, i, port.Port),
					Match:   []*networkingv1alpha3.HTTPMatchRequest{{Port: uint32(port.Port), Gateways: []string{"mesh"}}},
					Retries: retries,
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						{
							Headers: headers,
							Destination: &networkingv1alpha3.Destination{
								Host: destHost,
								Port: &networkingv1alpha3.PortSelector{Number: destPort},
							},
						},
					},
				})

			case model.TCP:
				// Ignore tcp routes if scale is zero
				if service.Scale.Replicas == 0 {
					continue
				}

				tcpRoutes = append(tcpRoutes, &networkingv1alpha3.TCPRoute{
					Match: []*networkingv1alpha3.L4MatchAttributes{{Port: uint32(port.Port)}},
					Route: []*networkingv1alpha3.RouteDestination{
						{
							Destination: &networkingv1alpha3.Destination{
								Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(port.Port)},
							},
						},
					},
				})
			}
		}
	}

	// We dont need to expose services since space cloud will take care of it
	// // Add http routes for the exposed http routes. Exposing a service is only supported for http services
	// if service.Expose != nil && service.Expose.Rules != nil && len(service.Expose.Rules) > 0 {
	// 	for i, rule := range service.Expose.Rules {
	// 		// Prepare variables
	// 		var headers *networkingv1alpha3.Headers
	// 		retries := &networkingv1alpha3.HTTPRetry{Attempts: 3, PerTryTimeout: &types.Duration{Seconds: 90}}
	// 		destHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID)
	// 		destPort := uint32(rule.Port)
	//
	// 		// Redirect traffic to runner when no of replicas is equal to zero. The runner proxy will scale up the service
	// 		// to service incoming requests.
	// 		if service.Scale.Replicas == 0 {
	// 			headers = &networkingv1alpha3.Headers{
	// 				Request: &networkingv1alpha3.Headers_HeaderOperations{
	// 					Set: map[string]string{
	// 						"x-og-project": service.ProjectID,
	// 						"x-og-service": service.ID,
	// 						"x-og-host":    destHost,
	// 						"x-og-port":    strconv.Itoa(int(destPort)),
	// 						"x-og-env":     service.Environment,
	// 						"x-og-version": service.Version,
	// 					},
	// 				},
	// 			}
	// 			retries = &networkingv1alpha3.HTTPRetry{Attempts: 1, PerTryTimeout: &types.Duration{Seconds: 180}}
	// 			destHost = "runner.space-cloud.svc.cluster.local"
	// 			destPort = proxyPort
	// 		}
	//
	// 		match := prepareHTTPMatch(&rule)
	// 		match[0].Gateways = []string{getGatewayName(service)}
	// 		httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
	// 			Match:   match,
	// 			Rewrite: prepareHTTPMatchRewrite(&rule),
	// 			Name:    fmt.Sprintf("expose-%d-%d", i, rule.Port),
	// 			Retries: retries,
	// 			Route: []*networkingv1alpha3.HTTPRouteDestination{
	// 				{
	// 					Headers: headers,
	// 					Destination: &networkingv1alpha3.Destination{
	// 						Host: destHost,
	// 						Port: &networkingv1alpha3.PortSelector{Number: destPort},
	// 					},
	// 				},
	// 			},
	// 		})
	// 	}
	// }

	return httpRoutes, tcpRoutes
}

// func prepareHTTPMatch(rule *model.ExposeRule) []*networkingv1alpha3.HTTPMatchRequest {
// 	// TODO: Add project level host
// 	if rule.URI.Exact != nil {
// 		return []*networkingv1alpha3.HTTPMatchRequest{
// 			{Uri: &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Exact{Exact: *rule.URI.Exact}}},
// 		}
//
// 	}
// 	if rule.URI.Prefix != nil {
// 		return []*networkingv1alpha3.HTTPMatchRequest{
// 			{Uri: &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Prefix{Prefix: *rule.URI.Prefix}}},
// 		}
// 	}
//
// 	return nil
// }

// func prepareHTTPMatchRewrite(rule *model.ExposeRule) *networkingv1alpha3.HTTPRewrite {
// 	if rule.URI.Rewrite != nil {
// 		return &networkingv1alpha3.HTTPRewrite{Uri: *rule.URI.Rewrite}
// 	}
// 	return nil
// }

func prepareVirtualServiceHosts(service *model.Service) []string {
	hosts := []string{fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID)}

	// if service.Expose != nil && service.Expose.Hosts != nil {
	// 	hosts = append(hosts, service.Expose.Hosts...)
	// }

	return hosts
}

// func prepareVirtualServiceGateways(service *model.Service) []string {
// 	gateways := []string{"mesh"}
//
// 	// Add gateway if the service is exposed
// 	if service.Expose != nil {
// 		gateways = append(gateways, getGatewayName(service))
// 	}
//
// 	return gateways
// }

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

	froms = append(froms, &securityv1beta1.Rule_From{Source: &securityv1beta1.Source{}})

	return []*securityv1beta1.Rule{{From: froms}}
}

func prepareUpstreamHosts(service *model.Service) []string {
	hosts := make([]string, len(service.Upstreams)+1)

	// First entry will always be space-cloud
	hosts[0] = "space-cloud/*"

	for i, upstream := range service.Upstreams {
		hosts[i+1] = upstream.ProjectID + "/" + upstream.Service
	}

	return hosts
}

func generateServiceAccount(service *model.Service) *v1.ServiceAccount {
	saName := getServiceAccountName(service)
	return &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: saName, Labels: map[string]string{"account": service.ID}}}
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
			Name: getDeploymentName(service),
			Labels: map[string]string{
				"app":     service.ID,
				"version": service.Version,
			},
			Annotations: map[string]string{
				"concurrency": strconv.Itoa(int(service.Scale.Concurrency)),
				"minReplicas": strconv.Itoa(int(service.Scale.MinReplicas)),
				"maxReplicas": strconv.Itoa(int(service.Scale.MaxReplicas)),
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
					ServiceAccountName: getServiceAccountName(service),
					Containers:         preparedContainer,
					Volumes:            volumes,
					ImagePullSecrets:   imagePull,
					// TODO: Add config for affinity
				},
			},
		},
	}
}

func generateService(service *model.Service) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   getServiceName(service.ID),
			Labels: map[string]string{"app": service.ID, "service": service.ID},
		},
		Spec: v1.ServiceSpec{
			Ports:    prepareServicePorts(service.Tasks),
			Selector: map[string]string{"app": service.ID},
			Type:     v1.ServiceTypeClusterIP,
		},
	}
}

func (i *Istio) generateVirtualService(service *model.Service) *v1alpha3.VirtualService {
	httpRoutes, tcpRoutes := prepareVirtualServiceRoutes(service, i.config.ProxyPort)
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{Name: getVirtualServiceName(service.ID)},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: prepareVirtualServiceHosts(service),
			// Gateways: prepareVirtualServiceGateways(service),
			Http: httpRoutes,
			Tcp:  tcpRoutes,
		},
	}
}

func generateDestinationRule(service *model.Service) *v1alpha3.DestinationRule {
	return &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{Name: getDestRuleName(service.ID)},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, service.ProjectID),
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				Tls: &networkingv1alpha3.TLSSettings{Mode: networkingv1alpha3.TLSSettings_ISTIO_MUTUAL},
			},
		},
	}
}

func generateAuthPolicy(service *model.Service) *v1beta1.AuthorizationPolicy {
	return &v1beta1.AuthorizationPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: getAuthorizationPolicyName(service)},
		Spec: securityv1beta1.AuthorizationPolicy{
			Selector: &v1beta12.WorkloadSelector{MatchLabels: map[string]string{"app": service.ID}},
			Rules:    prepareAuthPolicyRules(service),
		},
	}
}

func generateSidecarConfig(service *model.Service) *v1alpha3.Sidecar {
	return &v1alpha3.Sidecar{
		ObjectMeta: metav1.ObjectMeta{Name: getSidecarName(service.ID)},
		Spec: networkingv1alpha3.Sidecar{
			WorkloadSelector:      &networkingv1alpha3.WorkloadSelector{Labels: map[string]string{"app": service.ID}},
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
// 		ObjectMeta: metav1.ObjectMeta{Name: getGatewayName(service)},
// 		Spec: networkingv1alpha3.Gateway{
// 			Selector: map[string]string{"istio": "ingressgateway"},
// 			Servers: []*networkingv1alpha3.Server{
// 				{
// 					// TODO: make this https and load certificates dynamically
// 					Port: &networkingv1alpha3.Port{
// 						Number:   80,
// 						Name:     "http",
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

	return &resources
}
