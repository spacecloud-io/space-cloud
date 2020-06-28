package istio

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// AdjustScale adjusts the number of instances based on the number of active requests. It tries to make sure that
// no instance has more than the desired concurrency level. We simply change the number of replicas in the deployment
func (i *Istio) AdjustScale(ctx context.Context, service *model.Service, activeReqs int32) error {
	// We will process a single adjust scale request for a given service at any given time. We might miss out on some updates,
	// but the adjust scale routine will eventually make sure we reach the desired scale
	ns := service.ProjectID
	uniqueName := getServiceUniqueName(service.ProjectID, service.ID, service.Version)
	if _, loaded := i.adjustScaleLock.LoadOrStore(uniqueName, struct{}{}); loaded {
		logrus.Infof("Ignoring adjust scale request for service (%s:%s:%s) since another request is already in progress", ns, service.ID, service.Version)
		return nil
	}
	// Remove the lock once processing is done
	defer i.adjustScaleLock.Delete(uniqueName)

	logrus.Debugf("Adjusting scale of service (%s:%s:%s): Active reqs - %d", ns, service.ID, service.Version, activeReqs)
	deployment, err := i.cache.getDeployment(ctx, ns, getDeploymentName(service.ID, service.Version))
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

	// Update the replica count
	deployment.Spec.Replicas = &replicaCount
	if err := i.applyDeployment(ctx, ns, deployment); err != nil {
		logrus.Errorf("Could not adjust scale: %s", err.Error())
		return err
	}

	logrus.Infof("Scale of service (%s:%s) adjusted to %d successfully", ns, service.ID, replicaCount)
	return nil
}

// WaitForService adjusts scales, up the service to scale up the number of nodes from zero to one
// TODO: Do one watch per service. Right now its possible to have multiple watches for the same service
func (i *Istio) WaitForService(ctx context.Context, service *model.Service) error {
	ns := service.ProjectID
	logrus.Debugf("Scaling up service (%s:%s:%s) from zero", ns, service.ID, service.Version)

	// Scale up the service
	if err := i.AdjustScale(ctx, service, 1); err != nil {
		return err
	}

	timeout := int64(5 * 60)
	labels := fmt.Sprintf("app=%s,version=%s", service.ID, service.Version)
	logrus.Debugf("Watching for service (%s:%s:%s) to scale up and enter ready state", ns, service.ID, service.Version)
	watcher, err := i.kube.AppsV1().Deployments(ns).Watch(ctx, metav1.ListOptions{Watch: true, LabelSelector: labels, TimeoutSeconds: &timeout})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for ev := range watcher.ResultChan() {
		if ev.Type == watch.Error {
			return fmt.Errorf("service (%s:%s:%s) could not be scaled up", ns, service.ID, service.Version)
		}
		deployment := ev.Object.(*appsv1.Deployment)
		logrus.Debugf("Received watch event for service (%s:%s): available replicas - %d; ready replicas - %d", ns, service.ID, deployment.Status.AvailableReplicas, deployment.Status.ReadyReplicas)
		if deployment.Status.AvailableReplicas >= 1 && deployment.Status.ReadyReplicas >= 1 {
			return nil
		}
	}

	return fmt.Errorf("service (%s:%s) could not be started", ns, service.ID)
}
