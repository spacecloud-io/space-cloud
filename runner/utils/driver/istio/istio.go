package istio

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spaceuptech/space-cloud/runner/model"

	"github.com/spaceuptech/space-cloud/runner/utils/auth"
)

// Istio manages the istio on kubernetes deployment target
type Istio struct {
	// For internal use
	auth   *auth.Module
	config *Config

	// For tacking invocations to adjust scale
	adjustScaleLock sync.Map

	// Drivers to talk to k8s and istio
	kube  *kubernetes.Clientset
	istio *versionedclient.Clientset

	// For caching deployments
	cache *cache
}

// NewIstioDriver creates a new instance of the istio driver
func NewIstioDriver(auth *auth.Module, c *Config) (*Istio, error) {
	var restConfig *rest.Config
	var err error

	if c.IsInsideCluster {
		restConfig, err = rest.InClusterConfig()
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", c.KubeConfigPath)
	}
	if err != nil {
		return nil, err
	}

	// Create the kubernetes client
	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	// Create the istio client
	istio, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	// Create a cache
	cache, err := newCache(kube)
	if err != nil {
		return nil, err
	}

	return &Istio{auth: auth, config: c, kube: kube, istio: istio, cache: cache}, nil
}

func (i *Istio) getSecrets(service *model.Service) (map[string]*v1.Secret, error) {
	listOfSecrets := map[string]*v1.Secret{}
	tasks := service.Tasks
	for _, task := range tasks {
		for _, secretName := range task.Secrets {
			if _, p := listOfSecrets[secretName]; p {
				continue
			}
			secrets, err := i.kube.CoreV1().Secrets(service.ProjectID).Get(secretName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			listOfSecrets[secretName] = secrets
		}
	}
	return listOfSecrets, nil
}

// ApplyService deploys the service on istio
func (i *Istio) ApplyService(ctx context.Context, service *model.Service) error {
	// TODO: do we need to rollback on failure? rollback to previous version if it existed else remove. We also need to rollback the cache in this case
	// TODO: Add support for custom runtime
	// TODO: Add support for running multiple versions
	// We are hard coding the version right now. But we need to create rules for running multiple versions of
	// the same service. Also the traffic splitting between the versions needs to be configurable
	service.Version = "v1"

	ns := service.ProjectID

	// Set the default concurrency value to 50
	if service.Scale.Concurrency == 0 {
		service.Scale.Concurrency = 50
	}

	token, err := i.auth.GenerateTokenForArtifactStore(service.ID, service.ProjectID, service.Version)
	if err != nil {
		return err
	}

	// Create necessary resources
	// We do not need to setup destination rules since we will be enabling global mtls as described by this guide:
	// https://istio.io/docs/tasks/security/authentication/authn-policy/#globally-enabling-istio-mutual-tls
	// However we will need destination rules when routing between various versions

	// Get the list of secrets required for this service
	listOfSecrets, err := i.getSecrets(service)
	if err != nil {
		return err
	}

	// Create the appropriate kubernetes and istio objects
	kubeServiceAccount := generateServiceAccount(service)
	kubeDeployment := i.generateDeployment(service, token, listOfSecrets)
	kubeService := generateService(service)
	istioVirtualService := i.generateVirtualService(service)
	istioDestRule := generateDestinationRule(service)
	istioAuthPolicy := generateAuthPolicy(service)
	istioSidecar := generateSidecarConfig(service)

	// Create a service account if it doesn't already exist. This is used as the identity of the service.
	_, err = i.kube.CoreV1().ServiceAccounts(ns).Get(getServiceAccountName(service), metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create the resources since they dont exist
		logrus.Debugf("Creating service account for %s in %s", service.ID, ns)
		if _, err := i.kube.CoreV1().ServiceAccounts(ns).Create(kubeServiceAccount); err != nil {
			return err
		}

		logrus.Debugf("Creating deployment for %s in %s", service.ID, ns)
		if _, err := i.kube.AppsV1().Deployments(ns).Create(kubeDeployment); err != nil {
			return err
		}
		_ = i.cache.setDeployment(ns, kubeDeployment.Name, kubeDeployment)

		logrus.Debugf("Creating service for %s in %s", service.ID, ns)
		if _, err := i.kube.CoreV1().Services(ns).Create(kubeService); err != nil {
			return err
		}

		logrus.Debugf("Creating virtual service for %s in %s", service.ID, ns)
		if _, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Create(istioVirtualService); err != nil {
			return err
		}

		logrus.Debugf("Creating destination rule for %s in %s", service.ID, ns)
		if _, err := i.istio.NetworkingV1alpha3().DestinationRules(ns).Create(istioDestRule); err != nil {
			return err
		}

		logrus.Debugf("Creating auth policy for %s in %s", service.ID, ns)
		if _, err := i.istio.SecurityV1beta1().AuthorizationPolicies(ns).Create(istioAuthPolicy); err != nil {
			return err
		}

		logrus.Debugf("Creating sidecar config for %s in %s", service.ID, ns)
		if _, err := i.istio.NetworkingV1alpha3().Sidecars(ns).Create(istioSidecar); err != nil {
			return err
		}
	} else if err == nil {
		// Update the resources
		logrus.Debugf("Updating service for %s in %s", service.ID, ns)
		if _, err := i.kube.AppsV1().Deployments(ns).Update(kubeDeployment); err != nil {
			return err
		}
		if err := i.cache.setDeployment(ns, kubeDeployment.Name, kubeDeployment); err != nil {
			return err
		}

		logrus.Debugf("Updating service for %s in %s", service.ID, ns)
		prevService, err := i.kube.CoreV1().Services(ns).Get(kubeService.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		prevService.Spec.Ports = kubeService.Spec.Ports
		prevService.Labels = kubeService.Labels
		if _, err := i.kube.CoreV1().Services(ns).Update(prevService); err != nil {
			return err
		}

		logrus.Debugf("Updating virtual service for %s in %s", service.ID, ns)
		prevVirtualService, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(istioVirtualService.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		prevVirtualService.Spec = istioVirtualService.Spec
		prevVirtualService.Labels = istioVirtualService.Labels
		if _, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Update(prevVirtualService); err != nil {
			return err
		}

		logrus.Debugf("Updating auth policy for %s in %s", service.ID, ns)
		prevAuthPolicy, err := i.istio.SecurityV1beta1().AuthorizationPolicies(ns).Get(istioAuthPolicy.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		prevAuthPolicy.Spec = istioAuthPolicy.Spec
		prevAuthPolicy.Labels = istioAuthPolicy.Labels
		if _, err := i.istio.SecurityV1beta1().AuthorizationPolicies(ns).Update(prevAuthPolicy); err != nil {
			return err
		}

		logrus.Debugf("Updating sidecar config for %s in %s", service.ID, ns)
		prevSidecar, err := i.istio.NetworkingV1alpha3().Sidecars(ns).Get(istioSidecar.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		prevSidecar.Spec = istioSidecar.Spec
		prevSidecar.Labels = istioSidecar.Labels
		if _, err := i.istio.NetworkingV1alpha3().Sidecars(ns).Update(prevSidecar); err != nil {
			return err
		}
	} else {
		// Return error for unknown error
		return err
	}

	logrus.Infof("Service %s in %s applied successfully", service.ID, ns)
	return nil
}

// DeleteService deletes a service version
func (i *Istio) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	service := &model.Service{ID: serviceID, ProjectID: projectID, Version: version}

	// This will delete the service level service account. This will work fine as long as a service has only one version.
	// TODO: Add support for multiple versions
	service.Version = "v1"
	err := i.kube.CoreV1().ServiceAccounts(projectID).Delete(getServiceAccountName(service), &metav1.DeleteOptions{})
	if kubeErrors.IsNotFound(err) {
		// service account not found meaning no service present
		logrus.Errorf("Service does not exist - %v", err)
		return nil
	} else if err == nil {
		if err = i.kube.AppsV1().Deployments(projectID).Delete(getDeploymentName(service), &metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("error deleting service in istio unable to find deployment got error message - %v", err)
			return err
		}
		if err = i.kube.CoreV1().Services(projectID).Delete(getServiceName(serviceID), &metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("error deleting service in istio unable to find services got error message - %v", err)
			return err
		}
		// when we add versioning support, the destination rules and virtual services can only be removed when all service versions are removed
		if err = i.istio.NetworkingV1alpha3().VirtualServices(projectID).Delete(getVirtualServiceName(serviceID), &metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("error deleting service in istio unable to find virtual services got error message - %v", err)
			return err
		}
		if err = i.istio.NetworkingV1alpha3().DestinationRules(projectID).Delete(getDestRuleName(serviceID), &metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("error deleting service in istio unable to find destination rule got error message - %v", err)
			return err
		}
		if err = i.istio.SecurityV1beta1().AuthorizationPolicies(projectID).Delete(getAuthorizationPolicyName(service), &metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("error deleting service in istio unable to find authorization policies got error message - %v", err)
			return err
		}
		if err = i.istio.NetworkingV1alpha3().Sidecars(projectID).Delete(getSidecarName(serviceID), &metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("error deleting service in istio unable to find sidecars got error message - %v", err)
			return err
		}
	} else {
		logrus.Errorf("error deleting service in istio unknown error got error message - %v", err)
		return err
	}

	return nil
}

func (i *Istio) GetServices(ctx context.Context, projectId string) ([]*model.Service, error) {
	deploymentList, err := i.kube.AppsV1().Deployments(projectId).List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("error getting service in istio unable to find deployment got error message - %v", err)
		return nil, err
	}
	services := []*model.Service{}
	for _, deployment := range deploymentList.Items {
		service := new(model.Service)
		service.ProjectID = projectId
		service.ID = deployment.Labels["app"]
		service.Version = deployment.Labels["version"]
		s1, err := strconv.Atoi(deployment.Annotations["concurrency"])
		if err != nil {
			logrus.Errorf("error getting service in istio unable convert string to int annotation concurrency got error message - %v", err)
			return nil, err
		}
		s2, err := strconv.Atoi(deployment.Annotations["minReplicas"])
		if err != nil {
			logrus.Errorf("error getting service in istio unable convert string to int annotation minReplicas got error message - %v", err)
			return nil, err
		}
		s3, err := strconv.Atoi(deployment.Annotations["maxReplicas"])
		if err != nil {
			logrus.Errorf("error getting service in istio unable convert string to int annotation maxReplicas got error message - %v", err)
			return nil, err
		}
		service.Scale.Concurrency = int32(s1)
		service.Scale.MinReplicas = int32(s2)
		service.Scale.MaxReplicas = int32(s3)
		service.Scale.Replicas = *deployment.Spec.Replicas

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
			var secrets []string

			// get environment variables
			envs := map[string]string{}
			for _, env := range containerInfo.Env {
				if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
					secrets = append(secrets, env.ValueFrom.SecretKeyRef.LocalObjectReference.Name)
					continue
				}
				envs[env.Name] = env.Value
			}

			// Range over the file mounts for secrets
			for _, volume := range containerInfo.VolumeMounts {
				if checkIfVolumeIsSecret(volume.Name, deployment.Spec.Template.Spec.Volumes) {
					secrets = append(secrets, volume.Name)
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
					Image:  containerInfo.Image,
					Cmd:    containerInfo.Command,
					Secret: dockerSecret,
				},
				Env:     envs,
				Runtime: runtime,
				Secrets: secrets,
			})
		}

		// set whitelist
		authPolicy, _ := i.istio.SecurityV1beta1().AuthorizationPolicies(projectId).Get(getAuthorizationPolicyName(service), metav1.GetOptions{})
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
		sideCar, _ := i.istio.NetworkingV1alpha3().Sidecars(projectId).Get(service.ID, metav1.GetOptions{})
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

func checkIfVolumeIsSecret(name string, volumes []v1.Volume) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// AdjustScale adjusts the number of instances based on the number of active requests. It tries to make sure that
// no instance has more than the desired concurrency level. We simply change the number of replicas in the deployment
func (i *Istio) AdjustScale(service *model.Service, activeReqs int32) error {
	// We will process a single adjust scale request for a given service at any given time. We might miss out on some updates,
	// but the adjust scale routine will eventually make sure we reach the desired scale
	ns := service.ProjectID
	uniqueName := getServiceUniqueName(service.ProjectID, service.ID, service.Version)
	if _, loaded := i.adjustScaleLock.LoadOrStore(uniqueName, struct{}{}); loaded {
		logrus.Infof("Ignoring adjust scale request for service (%s:%s) since another request is already in progress", ns, service.ID)
		return nil
	}
	// Remove the lock once processing is done
	defer i.adjustScaleLock.Delete(uniqueName)

	logrus.Debugf("Adjusting scale of service (%s:%s): Active reqs - %d", ns, service.ID, activeReqs)
	deployment, err := i.cache.getDeployment(ns, getDeploymentName(service))
	if err != nil {
		return err
	}

	// Get the min and max replica numbers
	minReplicasString := deployment.Annotations["minReplicas"]
	maxReplicasString := deployment.Annotations["maxReplicas"]
	minReplicas, _ := strconv.Atoi(minReplicasString)
	maxReplicas, _ := strconv.Atoi(maxReplicasString)

	// Calculate the desired replica count
	concurrencyString := deployment.Annotations["concurrency"]
	concurrency, _ := strconv.Atoi(concurrencyString)
	replicaCount := int32(math.Ceil(float64(activeReqs) / float64(concurrency)))

	// Make sure the desired replica count doesn't cross the min and max range
	if replicaCount < int32(minReplicas) {
		replicaCount = int32(minReplicas)
	}
	if replicaCount > int32(maxReplicas) {
		replicaCount = int32(maxReplicas)
	}

	// Return if the existing replica count is the same
	if *deployment.Spec.Replicas == replicaCount {
		logrus.Debugf("Desired scale of service (%s:%s) is same as current scale (%d). Making no changes", ns, service.ID, replicaCount)
		return nil
	}

	// Update the virtual service if the new replica count is zero. This is required to redirect incoming http requests to
	// the runner proxy. The proxy is responsible to scale the service back up from zero.
	if replicaCount == 0 {
		virtualService, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(service.ID, metav1.GetOptions{})
		if err != nil {
			logrus.Errorf("Could not fetch virtual service (%s:%s) to adjust scale: %s", ns, service.ID, err.Error())
			return err
		}

		// Apply scale zero config to virtual service
		makeScaleZeroVirtualService(service, virtualService, i.config.ProxyPort)
		if _, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Update(virtualService); err != nil {
			logrus.Errorf("Could not revert virtual service (%s:%s) back to original to adjust scale: %s", ns, service.ID, err.Error())
			return err
		}
	}

	// Update the replica count
	deployment.Spec.Replicas = &replicaCount
	if _, err := i.kube.AppsV1().Deployments(ns).Update(deployment); err != nil {
		logrus.Errorf("Could not adjust scale: %s", err.Error())
		return err
	}
	if err := i.cache.setDeployment(ns, deployment.Name, deployment); err != nil {
		logrus.Errorf("Could not update cache in adjust scale: %s", err.Error())
		return err
	}

	logrus.Infof("Scale of service (%s:%s) adjusted to %d successfully", ns, service.ID, replicaCount)
	return nil
}

// WaitForService adjusts scales, up the service to scale up the number of nodes from zero to one
// TODO: Do one watch per service. Right now its possible to have multiple watches for the same service
func (i *Istio) WaitForService(service *model.Service) error {
	ns := service.ProjectID
	logrus.Debugf("Scaling up service (%s:%s) from zero", ns, service.ID)

	// Scale up the service
	if err := i.AdjustScale(service, 1); err != nil {
		return err
	}

	timeout := int64(3 * 60)
	labels := fmt.Sprintf("app=%s,version=%s", service.ID, service.Version)
	logrus.Debugf("Watching for service (%s:%s) to scale up and enter ready state", ns, service.ID)
	watcher, err := i.kube.AppsV1().Deployments(ns).Watch(metav1.ListOptions{Watch: true, LabelSelector: labels, TimeoutSeconds: &timeout})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for ev := range watcher.ResultChan() {
		deployment := ev.Object.(*appsv1.Deployment)
		logrus.Debugf("Received watch event for service (%s:%s): available replicas - %d; ready replicas - %d", ns, service.ID, deployment.Status.AvailableReplicas, deployment.Status.ReadyReplicas)
		if deployment.Status.AvailableReplicas >= 1 && deployment.Status.ReadyReplicas >= 1 {
			go func() {
				// Update the `virtual service` config of this service back to the original
				virtualService, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(service.ID, metav1.GetOptions{})
				if err != nil {
					logrus.Errorf("Could not fetch virtual service (%s:%s): %s", ns, service.ID, err.Error())
					return
				}

				// Revert back to the original configuration and apply that
				logrus.Debugf("Reverting routing rules back to original for service (%s:%s)", ns, service.ID)
				makeOriginalVirtualService(service, virtualService)
				if _, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Update(virtualService); err != nil {
					logrus.Errorf("Could not revert virtual service (%s:%s) back to original: %s", ns, service.ID, err.Error())
				}
				logrus.Infof("Routing rules reverted back to original for service (%s:%s) successfully", ns, service.ID)
			}()
			return nil
		}
	}

	return fmt.Errorf("service (%s:%s) could not be started", ns, service.ID)
}

// CreateProject creates a new namespace for the client
func (i *Istio) CreateProject(ctx context.Context, project *model.Project) error {
	// Project ID provided here is already in the form `project-env`
	namespace := project.ID
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: map[string]string{"istio-injection": "enabled"},
		},
	}
	_, err := i.kube.CoreV1().Namespaces().Create(ns)
	return err
}

// DeleteProject deletes a namespace for the client
func (i *Istio) DeleteProject(ctx context.Context, projectID string) error {
	return i.kube.CoreV1().Namespaces().Delete(projectID, &metav1.DeleteOptions{})
}

// Type returns the type of the driver
func (i *Istio) Type() model.DriverType {
	return model.TypeIstio
}
