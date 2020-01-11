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

func (i *Istio) prepareContainers(service *model.Service) []v1.Container {
	// There will be n + 1 containers in the pod. Each task will have it's own container. Along with that,
	// there will be a metric collection container as well which pushes metric data to the autoscaler.
	// TODO: Add support for private repos
	tasks := service.Tasks
	containers := make([]v1.Container, len(tasks))
	for i, task := range tasks {
		// Prepare env variables
		var envVars []v1.EnvVar
		for k, v := range task.Env {
			envVars = append(envVars, v1.EnvVar{Name: k, Value: v})
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

		containers[i] = v1.Container{
			Name: task.ID,
			Env:  envVars,

			// Resource Related
			Ports:     ports,
			Resources: *generateResourceRequirements(&task.Resources),

			// Docker related
			Image:           task.Docker.Image,
			Command:         cmd,
			Args:            args,
			ImagePullPolicy: v1.PullIfNotPresent,
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
		token, _ := i.auth.SignProxyToken(ksuid.New().String(), service.ProjectID, service.ID, service.Environment, service.Version)
		containers = append(containers, v1.Container{
			Name: "galaxy-metrics",
			Env:  []v1.EnvVar{{Name: "TOKEN", Value: token}},

			// Resource Related
			Resources: *generateResourceRequirements(&model.Resources{CPU: 20, Memory: 50}),

			// Docker related
			Image:           "spaceuptech/space-cloud/runner:latest",
			Command:         []string{"./galaxy"},
			Args:            []string{"proxy"},
			ImagePullPolicy: v1.PullIfNotPresent,
		})
	}

	return containers
}

func prepareContainerPorts(taskPorts []model.Port) []v1.ContainerPort {
	ports := make([]v1.ContainerPort, len(taskPorts))
	for i, p := range taskPorts {
		ports[i] = v1.ContainerPort{Name: p.Name, ContainerPort: p.Port}
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
	ogHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment))

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
	ogHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment))

	// Redirect traffic to galaxy runner when no of replicas is equal to zero. The galaxy proxy will scale up the service
	// to service incoming requests.
	for _, httpRoute := range virtualService.Spec.Http {
		for _, route := range httpRoute.Route {
			// Set the destination to galaxy runner proxy
			route.Destination.Host = "runner.galaxy.svc.cluster.local"
			route.Destination.Port.Number = proxyPort

			// Set the headers
			route.Headers = &networkingv1alpha3.Headers{
				Request: &networkingv1alpha3.Headers_HeaderOperations{
					Add: map[string]string{
						"x-og-project": service.ProjectID,
						"x-og-service": service.ID,
						"x-og-host":    ogHost,
						"x-og-port":    strings.Split(httpRoute.Name, "-")[2],
						"x-og-env":     service.Environment,
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
				destHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment))
				destPort := uint32(port.Port)

				// Redirect traffic to galaxy runner when no of replicas is equal to zero. The galaxy proxy will scale up the service
				// to service incoming requests.
				if service.Scale.Replicas == 0 {
					headers = &networkingv1alpha3.Headers{
						Request: &networkingv1alpha3.Headers_HeaderOperations{
							Add: map[string]string{
								"x-og-project": service.ProjectID,
								"x-og-service": service.ID,
								"x-og-host":    destHost,
								"x-og-port":    strconv.Itoa(int(destPort)),
								"x-og-env":     service.Environment,
								"x-og-version": service.Version,
							},
						},
					}
					retries = &networkingv1alpha3.HTTPRetry{Attempts: 1, PerTryTimeout: &types.Duration{Seconds: 180}}
					destHost = "runner.galaxy.svc.cluster.local"
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
								Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment)),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(port.Port)},
							},
						},
					},
				})
			}
		}
	}

	// Add http routes for the exposed http routes. Exposing a service is only supported for http services
	if service.Expose != nil && service.Expose.Rules != nil && len(service.Expose.Rules) > 0 {
		for i, rule := range service.Expose.Rules {
			// Prepare variables
			var headers *networkingv1alpha3.Headers
			retries := &networkingv1alpha3.HTTPRetry{Attempts: 3, PerTryTimeout: &types.Duration{Seconds: 90}}
			destHost := fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment))
			destPort := uint32(rule.Port)

			// Redirect traffic to galaxy runner when no of replicas is equal to zero. The galaxy proxy will scale up the service
			// to service incoming requests.
			if service.Scale.Replicas == 0 {
				headers = &networkingv1alpha3.Headers{
					Request: &networkingv1alpha3.Headers_HeaderOperations{
						Set: map[string]string{
							"x-og-project": service.ProjectID,
							"x-og-service": service.ID,
							"x-og-host":    destHost,
							"x-og-port":    strconv.Itoa(int(destPort)),
							"x-og-env":     service.Environment,
							"x-og-version": service.Version,
						},
					},
				}
				retries = &networkingv1alpha3.HTTPRetry{Attempts: 1, PerTryTimeout: &types.Duration{Seconds: 180}}
				destHost = "runner.galaxy.svc.cluster.local"
				destPort = proxyPort
			}

			match := prepareHTTPMatch(&rule)
			match[0].Gateways = []string{getGatewayName(service)}
			httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
				Match:   match,
				Rewrite: prepareHTTPMatchRewrite(&rule),
				Name:    fmt.Sprintf("expose-%d-%d", i, rule.Port),
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
		}
	}

	return httpRoutes, tcpRoutes
}

func prepareHTTPMatch(rule *model.ExposeRule) []*networkingv1alpha3.HTTPMatchRequest {
	// TODO: Add project level host
	if rule.URI.Exact != nil {
		return []*networkingv1alpha3.HTTPMatchRequest{
			{Uri: &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Exact{Exact: *rule.URI.Exact}}},
		}

	}
	if rule.URI.Prefix != nil {
		return []*networkingv1alpha3.HTTPMatchRequest{
			{Uri: &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Prefix{Prefix: *rule.URI.Prefix}}},
		}
	}

	return nil
}

func prepareHTTPMatchRewrite(rule *model.ExposeRule) *networkingv1alpha3.HTTPRewrite {
	if rule.URI.Rewrite != nil {
		return &networkingv1alpha3.HTTPRewrite{Uri: *rule.URI.Rewrite}
	}
	return nil
}

func prepareVirtualServiceHosts(service *model.Service) []string {
	hosts := []string{fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment))}

	if service.Expose != nil && service.Expose.Hosts != nil {
		hosts = append(hosts, service.Expose.Hosts...)
	}

	return hosts
}

func prepareVirtualServiceGateways(service *model.Service) []string {
	gateways := []string{"mesh"}

	// Add gateway if the service is exposed
	if service.Expose != nil {
		gateways = append(gateways, getGatewayName(service))
	}

	return gateways
}

func prepareAuthPolicyRules(service *model.Service) []*securityv1beta1.Rule {
	var froms []*securityv1beta1.Rule_From
	var namespaces []string
	var principals []string

	for _, whitelist := range service.Whitelist {
		array := strings.Split(whitelist, ":")
		projectID, service := array[0], array[1]

		if projectID == "*" {
			// This means this is an open service
			return []*securityv1beta1.Rule{{}}
		}

		if service == "*" {
			// This means that the service can be accessed from everyone in the project
			namespaces = append(namespaces, projectID)
		} else {
			// This means that the service can be accessed only from that service in that project
			principals = append(principals, fmt.Sprintf("cluster.local/ns/%s/sa/%s", projectID, service))
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

	// First entry will always be galaxy
	hosts[0] = "galaxy/*"

	for i, upstream := range service.Upstreams {
		hosts[i+1] = upstream.ProjectID + "/" + upstream.Service
	}

	return hosts
}

func generateServiceAccount(service *model.Service) *v1.ServiceAccount {
	saName := getServiceAccountName(service)
	return &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: saName, Labels: map[string]string{"account": service.ID}}}
}

func (i *Istio) generateDeployment(service *model.Service) *appsv1.Deployment {
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
					Annotations: map[string]string{"sidecar.istio.io/statsInclusionPrefixes": "cluster.outbound,listener,http,cluster_manager,listener_manager,http_mixer_filter,tcp_mixer_filter,server,cluster.xds-grpc"},
					Labels:      map[string]string{"app": service.ID, "version": service.Version},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: getServiceAccountName(service),
					Containers:         i.prepareContainers(service),
					// TODO: Add config for affinity
				},
			},
		},
	}
}

func generateService(service *model.Service) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   service.ID,
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
		ObjectMeta: metav1.ObjectMeta{Name: service.ID},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    prepareVirtualServiceHosts(service),
			Gateways: prepareVirtualServiceGateways(service),
			Http:     httpRoutes,
			Tcp:      tcpRoutes,
		},
	}
}

func generateDestinationRule(service *model.Service) *v1alpha3.DestinationRule {
	return &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{Name: service.ID},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.ID, getNamespaceName(service.ProjectID, service.Environment)),
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
		ObjectMeta: metav1.ObjectMeta{Name: service.ID},
		Spec: networkingv1alpha3.Sidecar{
			WorkloadSelector:      &networkingv1alpha3.WorkloadSelector{Labels: map[string]string{"app": service.ID}},
			Egress:                []*networkingv1alpha3.IstioEgressListener{{Hosts: prepareUpstreamHosts(service)}},
			OutboundTrafficPolicy: &networkingv1alpha3.OutboundTrafficPolicy{Mode: networkingv1alpha3.OutboundTrafficPolicy_ALLOW_ANY},
		},
	}
}

func generateGateways(service *model.Service) *v1alpha3.Gateway {
	hosts := make([]string, 0)

	// Add the hosts if provided
	if service.Expose != nil && service.Expose.Hosts != nil {
		hosts = append(hosts, service.Expose.Hosts...)
	}

	// Add dummy host if none exist
	hosts = append(hosts, "dummy.com")

	return &v1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{Name: getGatewayName(service)},
		Spec: networkingv1alpha3.Gateway{
			Selector: map[string]string{"istio": "ingressgateway"},
			Servers: []*networkingv1alpha3.Server{
				{
					// TODO: make this https and load certificates dynamically
					Port: &networkingv1alpha3.Port{
						Number:   80,
						Name:     "http",
						Protocol: "HTTP",
					},
					Hosts: hosts,
				},
			},
		},
	}
}

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
