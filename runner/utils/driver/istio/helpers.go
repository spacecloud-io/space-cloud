package istio

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/kedacore/keda/api/v1alpha1"
	"github.com/spaceuptech/helpers"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	securityv1beta1 "istio.io/api/security/v1beta1"
	v1beta12 "istio.io/api/type/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

const defaultAPIGroup string = "rbac.authorization.k8s.io"

func (i *Istio) prepareContainers(service *model.Service, listOfSecrets map[string]*v1.Secret) ([]v1.Container, []v1.Volume, []v1.LocalObjectReference) {
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
		// if task.Runtime == model.Code {
		// 	artifactURL := v1.EnvVar{Name: model.ArtifactURL, Value: i.config.ArtifactAddr}
		// 	artifactToken := v1.EnvVar{Name: model.ArtifactToken, Value: token}
		// 	artifactProject := v1.EnvVar{Name: model.ArtifactProject, Value: service.ProjectID}
		// 	artifactService := v1.EnvVar{Name: model.ArtifactService, Value: service.ID}
		// 	artifactVersion := v1.EnvVar{Name: model.ArtifactVersion, Value: service.Version}
		// 	envVars = append(envVars, artifactURL, artifactToken, artifactProject, artifactService, artifactVersion)
		// }

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
		proto := v1.Protocol(p.Protocol)
		switch p.Protocol {
		case model.HTTP, model.TCP:
			proto = v1.ProtocolTCP
		}

		ports[i] = v1.ContainerPort{Name: p.Name, ContainerPort: p.Port, Protocol: proto}
	}

	return ports
}

func prepareServicePorts(tasks []model.Task) []v1.ServicePort {
	var ports []v1.ServicePort
	for _, task := range tasks {
		for _, p := range task.Ports {
			proto := v1.Protocol(p.Protocol)
			switch p.Protocol {
			case model.HTTP, model.TCP:
				proto = v1.ProtocolTCP
			}

			ports = append(ports, v1.ServicePort{Name: p.Name, Port: p.Port, Protocol: proto})
		}
	}

	return ports
}

func prepareVirtualServiceHTTPRoutes(ctx context.Context, projectID, serviceID string, services map[string]model.AutoScaleConfig, routes model.Routes, proxyPort uint32) ([]*networkingv1alpha3.HTTPRoute, error) {
	var httpRoutes []*networkingv1alpha3.HTTPRoute

	for _, route := range routes {
		// Before v0.21.0 space-cloud only supported HTTP routes, because of that we never specified protocol while creating service routes
		// From v0.21.0, we support both TCP & HTTP routes, the protocol to be used is specified in the protocol field.
		// If the protocol field is empty we assume it to be an HTTP route to be backward compatible.
		if route.Source.Protocol != "" && route.Source.Protocol != model.HTTP {
			continue
		}

		// Check if the port provided is correct
		if route.Source.Port == 0 {
			return nil, errors.New("port cannot be zero")
		}

		if route.RequestTimeout == 0 {
			route.RequestTimeout = model.DefaultRequestTimeout
		}
		if route.RequestRetries == 0 {
			route.RequestRetries = model.DefaultRequestRetries
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
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("version (%s) not found for service (%s)", target.Version, serviceID), nil, nil)
				}

				// Prepare variables
				destHost := getInternalServiceDomain(projectID, serviceID, target.Version)
				destPort := uint32(target.Port)

				// Redirect traffic to runner when no of replicas is equal to zero. The runner proxy will scale up the service to service incoming requests.
				if versionScaleConfig.MinReplicas == 0 {
					destHost = "runner-proxy.space-cloud.svc.cluster.local"
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
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid target type (%s) provided", target.Type), nil, nil)
			}
		}

		matchers := make([]*networkingv1alpha3.HTTPMatchRequest, 0)
		for _, matcher := range route.Matchers {
			tempMatcher := new(networkingv1alpha3.HTTPMatchRequest)

			// Add url matchers
			if matcher.URL != nil {
				tempMatcher.IgnoreUriCase = matcher.URL.IgnoreCase
				tempMatcher.Uri = new(networkingv1alpha3.StringMatch)
				switch matcher.URL.Type {
				case model.RouteHTTPMatchTypeExact:
					tempMatcher.Uri.MatchType = &networkingv1alpha3.StringMatch_Exact{Exact: matcher.URL.Value}
				case model.RouteHTTPMatchTypePrefix:
					tempMatcher.Uri.MatchType = &networkingv1alpha3.StringMatch_Prefix{Prefix: matcher.URL.Value}
				case model.RouteHTTPMatchTypeRegex:
					tempMatcher.Uri.MatchType = &networkingv1alpha3.StringMatch_Regex{Regex: matcher.URL.Value}
				}
			}

			// 	Add header matchers
			tempMatcher.Headers = map[string]*networkingv1alpha3.StringMatch{}
			for _, header := range matcher.Headers {
				switch header.Type {
				case model.RouteHTTPMatchTypeExact:
					tempMatcher.Headers[header.Key] = &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Exact{Exact: header.Value}}
				case model.RouteHTTPMatchTypePrefix:
					tempMatcher.Headers[header.Key] = &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Prefix{Prefix: header.Value}}
				case model.RouteHTTPMatchTypeRegex:
					tempMatcher.Headers[header.Key] = &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Regex{Regex: header.Value}}
				case model.RouteHTTPMatchTypeCheckPresence:
					tempMatcher.Headers[header.Key] = &networkingv1alpha3.StringMatch{}
				}
			}

			tempMatcher.Port = uint32(route.Source.Port)
			tempMatcher.Gateways = []string{"mesh"}
			matchers = append(matchers, tempMatcher)
		}
		if len(matchers) == 0 {
			matchers = append(matchers, &networkingv1alpha3.HTTPMatchRequest{Port: uint32(route.Source.Port), Gateways: []string{"mesh"}})
		}

		// Add the http route
		httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
			Name:    fmt.Sprintf("http-%d", route.Source.Port),
			Match:   matchers,
			Retries: &networkingv1alpha3.HTTPRetry{Attempts: route.RequestRetries, PerTryTimeout: &types.Duration{Seconds: route.RequestTimeout}},
			Route:   destinations,
		})
	}

	return httpRoutes, nil
}

func prepareVirtualServiceTCPRoutes(ctx context.Context, projectID, serviceID string, services map[string]model.AutoScaleConfig, routes model.Routes) ([]*networkingv1alpha3.TCPRoute, error) {
	var tcpRoutes []*networkingv1alpha3.TCPRoute

	for _, route := range routes {
		// Route protocol can be either TCP or HTTP
		// this function is intended to create TCP routes only, so we are skipping routes whose protocol is not TCP
		if route.Source.Protocol != model.TCP {
			continue
		}

		// Check if the port provided is correct
		if route.Source.Port == 0 {
			return nil, errors.New("port cannot be zero")
		}

		// Check if at least one target is provided
		if len(route.Targets) == 0 {
			return nil, errors.New("at least one target needs to be provided")
		}

		// Prepare an array of targets / destinations
		var destinations []*networkingv1alpha3.RouteDestination
		for _, target := range route.Targets {
			switch target.Type {
			case model.RouteTargetVersion:
				// Check if config for version exists
				_, p := services[target.Version]
				if !p {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("version (%s) not found for service (%s)", target.Version, serviceID), nil, nil)
				}

				// Prepare variables
				destHost := getInternalServiceDomain(projectID, serviceID, target.Version)
				destPort := uint32(target.Port)

				destinations = append(destinations, &networkingv1alpha3.RouteDestination{
					Destination: &networkingv1alpha3.Destination{
						Host: destHost,
						Port: &networkingv1alpha3.PortSelector{Number: destPort},
					},
					Weight: target.Weight,
				})

			case model.RouteTargetExternal:
				destinations = append(destinations, &networkingv1alpha3.RouteDestination{
					Destination: &networkingv1alpha3.Destination{
						Host: target.Host,
						Port: &networkingv1alpha3.PortSelector{Number: uint32(target.Port)},
					},
					Weight: target.Weight,
				})

			default:
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid target type (%s) provided", target.Type), nil, nil)
			}
		}

		// Add the http route
		tcpRoutes = append(tcpRoutes, &networkingv1alpha3.TCPRoute{
			Match: []*networkingv1alpha3.L4MatchAttributes{{Port: uint32(route.Source.Port)}},
			Route: destinations,
		})
	}

	return tcpRoutes, nil
}

func updateOrCreateVirtualServiceRoutes(service *model.Service, proxyPort uint32, prevVirtualService *v1alpha3.VirtualService) ([]*networkingv1alpha3.HTTPRoute, []*networkingv1alpha3.TCPRoute) {
	// Update the existing destinations of this version if virtual service already exist. We only need to do this for http services.
	if prevVirtualService != nil {
		for _, httpRoute := range prevVirtualService.Spec.Http {
			for _, dest := range httpRoute.Route {

				// Check if the route was for a service with min scale 0. If the destination has the host of runner, it means it is communicating via the proxy.
				if dest.Destination.Host == "runner-proxy.space-cloud.svc.cluster.local" {
					// We are only interested in this case if the new min replica for this version is more than 0. If the min replica was zero there would be no change
					if service.AutoScale.MinReplicas == 0 {
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
				if service.AutoScale.MinReplicas > 0 {
					continue
				}

				// Update the destination to communicate via the proxy if its for our version
				if dest.Destination.Host == getInternalServiceDomain(service.ProjectID, service.ID, service.Version) {
					dest.Destination.Host = "runner-proxy.space-cloud.svc.cluster.local"
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
				if service.AutoScale.MinReplicas == 0 {
					destHost = "runner-proxy.space-cloud.svc.cluster.local"
					destPort = proxyPort
				}

				httpRoutes = append(httpRoutes, &networkingv1alpha3.HTTPRoute{
					Name:    fmt.Sprintf("http-%d%d-%d", j, i, port.Port),
					Match:   []*networkingv1alpha3.HTTPMatchRequest{{Port: uint32(port.Port), Gateways: []string{"mesh"}}},
					Retries: &networkingv1alpha3.HTTPRetry{Attempts: model.DefaultRequestRetries, PerTryTimeout: &types.Duration{Seconds: model.DefaultRequestTimeout}},
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
		Name: getServiceAccountName(service.ID),
		Labels: map[string]string{
			"app.kubernetes.io/name":       service.ID,
			"app.kubernetes.io/managed-by": "space-cloud",
			"space-cloud.io/version":       model.Version,
		},
	}}
}

func (i *Istio) generateKedaConfig(ctx context.Context, service *model.Service) (*v1alpha1.ScaledObject, []v1alpha1.TriggerAuthentication, error) {
	// Create an empty trigger authentication ref error
	triggerAuthRefs := make([]v1alpha1.TriggerAuthentication, 0)

	// Generate a default auto scale config if not provided
	if service.AutoScale == nil {
		service.AutoScale = getDefaultAutoScaleConfig()

		// Load value from the previous scale object
		if service.Scale != nil {
			mode := "requests-per-second"
			if service.Scale.Mode == "parallel" {
				mode = "active-requests"
			}
			service.AutoScale.MinReplicas = service.Scale.MinReplicas
			service.AutoScale.MaxReplicas = service.Scale.MaxReplicas
			service.AutoScale.Triggers = append(service.AutoScale.Triggers, model.AutoScaleTrigger{
				Type: mode,
				Name: mode,
				MetaData: map[string]string{
					"target": strconv.Itoa(int(service.Scale.Concurrency)),
				},
			})
		}
	}

	// Set default values for auto scale config
	if service.AutoScale.MaxReplicas == 0 {
		service.AutoScale.MaxReplicas = 100
	}
	if service.AutoScale.PollingInterval == 0 {
		service.AutoScale.PollingInterval = 15
	}
	if service.AutoScale.CoolDownInterval == 0 {
		service.AutoScale.CoolDownInterval = 120
	}

	// return nil value if no triggers are provided
	if len(service.AutoScale.Triggers) == 0 {
		return nil, triggerAuthRefs, nil
	}

	// A variable for the advanced config. We want the advanced config to be nil unless it is specifically needed.
	var advancedConfig *v1alpha1.AdvancedConfig

	// Prepare the triggers
	triggers := make([]v1alpha1.ScaleTriggers, 0)
	for _, trigger := range service.AutoScale.Triggers {
		switch trigger.Type {
		case "requests-per-second", "active-requests":
			// Check if target is provided
			target, p := trigger.MetaData["target"]
			if !p {
				return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Missing field (target) in scaling trigger (%s)", trigger.Type), nil, nil)
			}

			triggers = append(triggers, v1alpha1.ScaleTriggers{
				Type: "external-push",
				Name: trigger.Name,
				Metadata: map[string]string{
					"scalerAddress": "runner.space-cloud.svc.cluster.local:4060",
					"scaler":        "space-cloud.io/scaler",
					"type":          trigger.Type,
					"target":        target,
					"service":       service.ID,
					"version":       service.Version,
					"project":       service.ProjectID,
					"minReplicas":   strconv.Itoa(int(service.AutoScale.MinReplicas)),
				},
			})

		default:
			// Create a nil authRef object. We want it to be nil unless it is specifically needed.
			var authRef *v1alpha1.ScaledObjectAuthRef
			if trigger.AuthenticatedRef != nil {
				// Make the param mapping array
				secretTargetRefs := make([]v1alpha1.AuthSecretTargetRef, len(trigger.AuthenticatedRef.SecretMapping))
				for i, ref := range trigger.AuthenticatedRef.SecretMapping {
					key := strings.TrimPrefix(ref.Key, "secrets.")
					arr := strings.Split(key, ".")
					if len(arr) != 2 {
						return nil, nil, fmt.Errorf("invalid value (%s) provided for secret key", ref.Key)
					}

					secretTargetRefs[i] = v1alpha1.AuthSecretTargetRef{
						Name:      arr[0],
						Key:       arr[1],
						Parameter: ref.Parameter,
					}
				}

				// Generate a unique name for the trigger auth
				name := getKedaTriggerAuthName(service.ID, service.Version, trigger.Name)

				// Add the trigger authentication object
				triggerAuthRefs = append(triggerAuthRefs, v1alpha1.TriggerAuthentication{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
						Labels: map[string]string{
							"app":                          service.ID,
							"version":                      service.Version,
							"app.kubernetes.io/name":       service.ID,
							"app.kubernetes.io/version":    service.Version,
							"app.kubernetes.io/managed-by": "space-cloud",
							"space-cloud.io/version":       model.Version,
						},
					},
					Spec: v1alpha1.TriggerAuthenticationSpec{
						SecretTargetRef: secretTargetRefs,
					},
				})

				// Don't forget to populate the auth ref object
				authRef = &v1alpha1.ScaledObjectAuthRef{
					Name: name,
				}
			}

			// Add the trigger to the list of triggers
			triggers = append(triggers, v1alpha1.ScaleTriggers{
				Type:              trigger.Type,
				Name:              trigger.Name,
				Metadata:          trigger.MetaData,
				AuthenticationRef: authRef,
			})
		}
	}

	// Prepare the keda config
	kedaConfig := &v1alpha1.ScaledObject{
		ObjectMeta: metav1.ObjectMeta{
			Name: getKedaScaledObjectName(service.ID, service.Version),
			Labels: map[string]string{
				"app":                          service.ID,
				"version":                      service.Version,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/version":    service.Version,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: v1alpha1.ScaledObjectSpec{
			Triggers:        triggers,
			PollingInterval: &service.AutoScale.PollingInterval,
			CooldownPeriod:  &service.AutoScale.CoolDownInterval,
			MinReplicaCount: &service.AutoScale.MinReplicas,
			MaxReplicaCount: &service.AutoScale.MaxReplicas,
			Advanced:        advancedConfig,
			ScaleTargetRef: &v1alpha1.ScaleTarget{
				Name: getDeploymentName(service.ID, service.Version),
				Kind: "Deployment", // Change this to stateful set when necessary
			},
		},
	}

	return kedaConfig, triggerAuthRefs, nil
}

func (i *Istio) generateDeployment(service *model.Service, listOfSecrets map[string]*v1.Secret) *appsv1.Deployment {
	preparedContainer, volumes, imagePull := i.prepareContainers(service, listOfSecrets)

	// Set the default stats inclusion prefix
	if service.StatsInclusionPrefixes == "" {
		service.StatsInclusionPrefixes = "http.inbound,cluster_manager,listener_manager"
	}
	if !strings.Contains(service.StatsInclusionPrefixes, "http.inbound") {
		service.StatsInclusionPrefixes += ",http.inbound"
	}

	var nodeAffinity *v1.NodeAffinity
	var podAffinity *v1.PodAffinity
	var podAntiAffinity *v1.PodAntiAffinity
	for _, affinity := range service.Affinity {
		switch affinity.Type {
		case model.AffinityTypeService:
			// affinity
			if affinity.Weight > 0 {
				if podAffinity == nil {
					podAffinity = &v1.PodAffinity{}
				}
				required, preferred := getServiceAffinityObject(affinity, 1)
				if preferred != nil {
					podAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(podAffinity.PreferredDuringSchedulingIgnoredDuringExecution, *preferred)
				}
				if required != nil {
					podAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(podAffinity.RequiredDuringSchedulingIgnoredDuringExecution, *required)
				}
			} else {
				if podAntiAffinity == nil {
					podAntiAffinity = &v1.PodAntiAffinity{}
				}
				required, preferred := getServiceAffinityObject(affinity, -1)
				if preferred != nil {
					podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution, *preferred)
				}
				if required != nil {
					podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, *required)
				}
			}
		case model.AffinityTypeNode:
			// affinity
			if nodeAffinity == nil {
				nodeAffinity = &v1.NodeAffinity{}
			}
			required, preferred := getNodeAffinityObject(affinity)
			if preferred != nil {
				nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution, *preferred)
			}
			if required != nil {
				if nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
					nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &v1.NodeSelector{
						NodeSelectorTerms: []v1.NodeSelectorTerm{},
					}
				}
				nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms, *required)
			}
		}
	}

	// Set default labels if not present already
	if service.Labels == nil {
		service.Labels = map[string]string{}
	}

	labels := service.Labels
	labels["app"] = service.ID
	labels["version"] = service.Version
	labels["app.kubernetes.io/name"] = service.ID
	labels["app.kubernetes.io/version"] = service.Version
	labels["app.kubernetes.io/managed-by"] = "space-cloud"
	labels["space-cloud.io/version"] = model.Version

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   getDeploymentName(service.ID, service.Version),
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &service.AutoScale.MinReplicas,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app":     service.ID,
				"version": service.Version,
			}},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"sidecar.istio.io/statsInclusionPrefixes": service.StatsInclusionPrefixes},
					Labels:      labels,
				},
				Spec: v1.PodSpec{
					ServiceAccountName: getServiceAccountName(service.ID),
					Containers:         preparedContainer,
					Volumes:            volumes,
					ImagePullSecrets:   imagePull,
					Affinity: &v1.Affinity{
						NodeAffinity:    nodeAffinity,
						PodAffinity:     podAffinity,
						PodAntiAffinity: podAntiAffinity,
					},
				},
			},
		},
	}
}

func getServiceAffinityObject(affinity model.Affinity, multiplier int32) (*v1.PodAffinityTerm, *v1.WeightedPodAffinityTerm) {
	expressions := []metav1.LabelSelectorRequirement{}
	for _, expression := range affinity.MatchExpressions {
		expressions = append(expressions, metav1.LabelSelectorRequirement{
			Key:      expression.Key,
			Operator: metav1.LabelSelectorOperator(expression.Operator),
			Values:   expression.Values,
		})
	}
	switch affinity.Operator {
	case model.AffinityOperatorRequired:
		return &v1.PodAffinityTerm{
			LabelSelector: &metav1.LabelSelector{
				MatchLabels:      nil,
				MatchExpressions: expressions,
			},
			Namespaces:  affinity.Projects,
			TopologyKey: affinity.TopologyKey,
		}, nil
	case model.AffinityOperatorPreferred:
		return nil, &v1.WeightedPodAffinityTerm{
			Weight: affinity.Weight * multiplier,
			PodAffinityTerm: v1.PodAffinityTerm{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels:      nil,
					MatchExpressions: expressions,
				},
				Namespaces:  affinity.Projects,
				TopologyKey: affinity.TopologyKey,
			},
		}
	}
	return nil, nil
}

func getNodeAffinityObject(affinity model.Affinity) (*v1.NodeSelectorTerm, *v1.PreferredSchedulingTerm) {
	expressions := []v1.NodeSelectorRequirement{}
	for _, expression := range affinity.MatchExpressions {
		expressions = append(expressions, v1.NodeSelectorRequirement{
			Key:      expression.Key,
			Operator: v1.NodeSelectorOperator(expression.Operator),
			Values:   expression.Values,
		})
	}
	switch affinity.Operator {
	case model.AffinityOperatorRequired:
		return &v1.NodeSelectorTerm{MatchExpressions: expressions}, nil
	case model.AffinityOperatorPreferred:
		return nil, &v1.PreferredSchedulingTerm{
			Weight: affinity.Weight,
			Preference: v1.NodeSelectorTerm{
				MatchExpressions: expressions,
				MatchFields:      nil,
			},
		}
	}
	return nil, nil
}

func generateGeneralService(service *model.Service) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getServiceName(service.ID),
			Labels: map[string]string{
				"app":                          service.ID,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: prepareServicePorts(service.Tasks),
			Selector: map[string]string{
				"app": service.ID,
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}
}

func generateInternalService(service *model.Service) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getInternalServiceName(service.ID, service.Version),
			Labels: map[string]string{
				"app":                          service.ID,
				"version":                      service.Version,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/version":    service.Version,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: prepareServicePorts(service.Tasks),
			Selector: map[string]string{
				"app":     service.ID,
				"version": service.Version,
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}
}

func (i *Istio) updateVirtualService(service *model.Service, prevVirtualService *v1alpha3.VirtualService) *v1alpha3.VirtualService {
	httpRoutes, tcpRoutes := updateOrCreateVirtualServiceRoutes(service, i.config.ProxyPort, prevVirtualService)
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getVirtualServiceName(service.ID),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"app":                          service.ID,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: prepareVirtualServiceHosts(service),
			// Gateways: prepareVirtualServiceGateways(service),
			Http: httpRoutes,
			Tcp:  tcpRoutes,
		},
	}
}
func (i *Istio) generateVirtualServiceBasedOnRoutes(ctx context.Context, projectID, serviceID string, scaleConfig map[string]model.AutoScaleConfig, routes model.Routes) (*v1alpha3.VirtualService, error) {
	// Generate the httpRoutes based on the routes provided
	httpRoutes, err := prepareVirtualServiceHTTPRoutes(ctx, projectID, serviceID, scaleConfig, routes, i.config.ProxyPort)
	if err != nil {
		return nil, err
	}

	tcpRoutes, err := prepareVirtualServiceTCPRoutes(ctx, projectID, serviceID, scaleConfig, routes)
	if err != nil {
		return nil, err
	}

	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name: getVirtualServiceName(serviceID),
			Labels: map[string]string{
				"app":                          serviceID,
				"app.kubernetes.io/name":       serviceID,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: []string{getServiceDomainName(projectID, serviceID)},
			Http:  httpRoutes,
			Tcp:   tcpRoutes,
		},
	}, nil
}

func generateGeneralDestinationRule(service *model.Service) *v1alpha3.DestinationRule {
	return &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: getGeneralDestRuleName(service.ID),
			Labels: map[string]string{
				"app":                          service.ID,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
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
			Name: getInternalDestRuleName(service.ID, service.Version),
			Labels: map[string]string{
				"app":                          service.ID,
				"version":                      service.Version,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/version":    service.Version,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
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
			Name: getAuthorizationPolicyName(service.ProjectID, service.ID, service.Version),
			Labels: map[string]string{
				"app":                          service.ID,
				"version":                      service.Version,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/version":    service.Version,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: securityv1beta1.AuthorizationPolicy{
			Selector: &v1beta12.WorkloadSelector{MatchLabels: map[string]string{
				"app":     service.ID,
				"version": service.Version,
			}},
			Rules: prepareAuthPolicyRules(service),
		},
	}
	return authPolicy
}

func generateSidecarConfig(service *model.Service) *v1alpha3.Sidecar {
	return &v1alpha3.Sidecar{
		ObjectMeta: metav1.ObjectMeta{
			Name: getSidecarName(service.ID, service.Version),
			Labels: map[string]string{
				"app":                          service.ID,
				"version":                      service.Version,
				"app.kubernetes.io/name":       service.ID,
				"app.kubernetes.io/version":    service.Version,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/version":       model.Version,
			},
		},
		Spec: networkingv1alpha3.Sidecar{
			WorkloadSelector: &networkingv1alpha3.WorkloadSelector{Labels: map[string]string{
				"app":     service.ID,
				"version": service.Version,
			}},
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

func getDefaultAutoScaleConfig() *model.AutoScaleConfig {
	return &model.AutoScaleConfig{
		PollingInterval:  15,
		CoolDownInterval: 120,
		MinReplicas:      1,
		MaxReplicas:      100,
		Triggers:         []model.AutoScaleTrigger{},
	}
}

func (i *Istio) generateServiceRole(ctx context.Context, role *model.Role) (*v12.Role, *v12.RoleBinding) {
	rules := make([]v12.PolicyRule, 0)
	for _, rule := range role.Rules {
		rules = append(rules, v12.PolicyRule{APIGroups: rule.APIGroups, Verbs: rule.Verbs, Resources: rule.Resources})

	}

	return &v12.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      role.ID,
				Namespace: role.Project,
				Labels:    i.generateServiceRoleLabels(role),
			},
			Rules: rules,
		}, &v12.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      role.ID,
				Namespace: role.Project,
				Labels:    i.generateServiceRoleLabels(role),
			},
			Subjects: []v12.Subject{
				{
					Kind:     "ServiceAccount",
					Name:     getServiceAccountName(role.Service),
					APIGroup: "",
				},
			},
			RoleRef: v12.RoleRef{
				Kind:     "Role",
				Name:     role.ID,
				APIGroup: defaultAPIGroup,
			},
		}

}

func (i *Istio) generateServiceClusterRole(ctx context.Context, role *model.Role) (*v12.ClusterRole, *v12.ClusterRoleBinding) {
	rules := make([]v12.PolicyRule, 0)
	for _, rule := range role.Rules {
		rules = append(rules, v12.PolicyRule{APIGroups: rule.APIGroups, Verbs: rule.Verbs, Resources: rule.Resources})

	}

	return &v12.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:   role.ID,
				Labels: i.generateServiceRoleLabels(role),
			},
			Rules: rules,
		}, &v12.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:   role.ID,
				Labels: i.generateServiceRoleLabels(role),
			},
			Subjects: []v12.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      getServiceAccountName(role.Service),
					APIGroup:  "",
					Namespace: role.Project,
				},
			},
			RoleRef: v12.RoleRef{
				Kind:     "ClusterRole",
				Name:     role.ID,
				APIGroup: defaultAPIGroup,
			},
		}
}

func (i *Istio) generateServiceRoleLabels(role *model.Role) map[string]string {
	labels := make(map[string]string)
	labels["app"] = role.Service
	labels["app.kubernetes.io/name"] = role.Service
	labels["app.kubernetes.io/managed-by"] = "space-cloud"
	labels["space-cloud.io/version"] = model.Version
	return labels
}
