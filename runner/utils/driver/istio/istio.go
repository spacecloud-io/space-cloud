package istio

import (
	"fmt"
	"math"
	"strconv"
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

	return &Istio{auth: auth, config: c, kube: kube, istio: istio}, nil
}

// ApplyService deploys the service on istio
func (i *Istio) ApplyService(service *model.Service) error {
	// TODO: do we need to rollback on failure? rollback to previous version if it existed else remove
	// TODO: Add support for custom runtime
	// TODO: Add support for running multiple versions
	// We are hard coding the version right now. But we need to create rules for running multiple versions of
	// the same service. Also the traffic splitting between the versions needs to be configurable
	service.Version = "v1"

	ns := getNamespaceName(service.ProjectID, service.Environment)

	// Set the default concurrency value to 50
	if service.Scale.Concurrency == 0 {
		service.Scale.Concurrency = 50
	}

	// Create necessary resources
	// We do not need to setup destination rules since we will be enabling global mtls as described by this guide:
	// https://istio.io/docs/tasks/security/authentication/authn-policy/#globally-enabling-istio-mutual-tls
	// However we will need destination rules when routing between various versions
	kubeServiceAccount := generateServiceAccount(service)
	kubeDeployment := i.generateDeployment(service)
	kubeService := generateService(service)
	istioVirtualService := i.generateVirtualService(service)
	istioDestRule := generateDestinationRule(service)
	istioGateway := generateGateways(service)
	istioAuthPolicy := generateAuthPolicy(service)
	istioSidecar := generateSidecarConfig(service)

	// Create a service account if it doesn't already exist. This is used as the identity of the service.
	_, err := i.kube.CoreV1().ServiceAccounts(ns).Get(getServiceAccountName(service), metav1.GetOptions{})
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

		logrus.Debugf("Creating gateway for %s in %s", service.ID, ns)
		if _, err := i.istio.NetworkingV1alpha3().Gateways(ns).Create(istioGateway); err != nil {
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

		logrus.Debugf("Updating gateway for %s in %s", service.ID, ns)
		prevGateway, err := i.istio.NetworkingV1alpha3().Gateways(ns).Get(istioGateway.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		prevGateway.Spec = istioGateway.Spec
		prevGateway.Labels = istioGateway.Labels
		if _, err := i.istio.NetworkingV1alpha3().Gateways(ns).Update(prevGateway); err != nil {
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

// AdjustScale adjusts the number of instances based on the number of active requests. It tries to make sure that
// no instance has more than the desired concurrency level. We simply change the number of replicas in the deployment
func (i *Istio) AdjustScale(service *model.Service, activeReqs int32) error {
	// We will process a single adjust scale request for a given service at any given time. We might miss out on some updates,
	// but the adjust scale routine will eventually make sure we reach the desired scale
	ns := getNamespaceName(service.ProjectID, service.Environment)
	uniqueName := getServiceUniqueName(service.ProjectID, service.ID, service.Environment, service.Version)
	if _, loaded := i.adjustScaleLock.LoadOrStore(uniqueName, struct{}{}); loaded {
		logrus.Infof("Ignoring adjust scale request for service (%s:%s) since another request is already in progress", ns, service.ID)
		return nil
	}
	// Remove the lock once processing is done
	defer i.adjustScaleLock.Delete(uniqueName)

	logrus.Debugf("Adjusting scale of service (%s:%s): Active reqs - %d", ns, service.ID, activeReqs)
	deployment, err := i.kube.AppsV1().Deployments(ns).Get(getDeploymentName(service), metav1.GetOptions{})
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
	// the galaxy runner proxy. The proxy is responsible to scale the service back up from zero.
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

	logrus.Infof("Scale of of service (%s:%s) adjusted to %d successfully", ns, service.ID, replicaCount)
	return nil
}

// WaitForService adjusts scales, up the service to scale up the number of nodes from zero to one
// TODO: Do one watch per service. Right now its possible to have multiple watches for the same service
func (i *Istio) WaitForService(service *model.Service) error {
	ns := getNamespaceName(service.ProjectID, service.Environment)
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
func (i *Istio) CreateProject(project *model.Environment) error {
	ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: project.ID, Labels: map[string]string{"istio-injection": "enabled"}}}
	_, err := i.kube.CoreV1().Namespaces().Create(ns)
	return err
}

// Type returns the type of the driver
func (i *Istio) Type() model.DriverType {
	return model.TypeIstio
}
