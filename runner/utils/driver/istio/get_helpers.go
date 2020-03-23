package istio

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func (i *Istio) getPreviousVirtualServiceIfExists(ns, service string) (*v1alpha3.VirtualService, error) {
	prevVirtualService, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(getVirtualServiceName(service), metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// We'll simple send `nil` if the virtual service did not actually exist. This is important since it indicates that
		// a virtual service needs to be created
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return prevVirtualService, nil
}

func (i *Istio) getServiceDeployments(ns, serviceID string) (*appsv1.DeploymentList, error) {
	return i.kube.AppsV1().Deployments(ns).List(metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", serviceID)})
}

func (i *Istio) getServiceDeploymentsCount(ns, serviceID string) (int, error) {
	deployments, err := i.getServiceDeployments(ns, serviceID)
	if err != nil {
		return 0, err
	}

	return len(deployments.Items), nil
}

func (i *Istio) getVirtualServices(ns string) (*v1alpha3.VirtualServiceList, error) {
	return i.istio.NetworkingV1alpha3().VirtualServices(ns).List(metav1.ListOptions{})
}

func (i *Istio) getAllVersionScaleConfig(ns, serviceID string) (map[string]model.ScaleConfig, error) {
	// Get all deployments of the provided service
	deployments, err := i.getServiceDeployments(ns, serviceID)
	if err != nil {
		return nil, err
	}

	// Throw error if the deployment contains no config at all
	if len(deployments.Items) == 0 {
		return nil, fmt.Errorf("no versions of service (%s) has been deployed", serviceID)
	}

	// Load the scale config of each version
	c := make(map[string]model.ScaleConfig, len(deployments.Items))
	for _, deployment := range deployments.Items {
		scale, err := getScaleConfigFromDeployment(deployment)
		if err != nil {
			return nil, err
		}
		c[deployment.Labels["version"]] = scale
	}

	return c, nil
}

func getScaleConfigFromDeployment(deployment appsv1.Deployment) (model.ScaleConfig, error) {
	concurrency, err := strconv.Atoi(deployment.Annotations["concurrency"])
	if err != nil {
		logrus.Errorf("Error getting service in istio - unable convert string to int annotation concurrency - %v", err)
		return model.ScaleConfig{}, err
	}
	minReplicas, err := strconv.Atoi(deployment.Annotations["minReplicas"])
	if err != nil {
		logrus.Errorf("Error getting service in istio - unable convert string to int annotation minReplicas - %v", err)
		return model.ScaleConfig{}, err
	}
	maxReplicas, err := strconv.Atoi(deployment.Annotations["maxReplicas"])
	if err != nil {
		logrus.Errorf("Error getting service in istio - unable convert string to int annotation maxReplicas - %v", err)
		return model.ScaleConfig{}, err
	}

	mode := deployment.Annotations["mode"]
	if mode == "" {
		mode = "per-second"
	}

	return model.ScaleConfig{
		Concurrency: int32(concurrency),
		MinReplicas: int32(minReplicas),
		MaxReplicas: int32(maxReplicas),
		Replicas:    *deployment.Spec.Replicas,
		Mode:        mode,
	}, nil
}
