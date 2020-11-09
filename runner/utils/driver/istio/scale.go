package istio

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// WaitForService adjusts scales, up the service to scale up the number of nodes from zero to one
// TODO: Do one watch per service. Right now its possible to have multiple watches for the same service
func (i *Istio) WaitForService(ctx context.Context, service *model.Service) error {
	ns := service.ProjectID
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Scaling up service (%s:%s:%s) from zero", ns, service.ID, service.Version), nil)

	timeout := int64(5 * 60)
	labels := fmt.Sprintf("app=%s,version=%s", service.ID, service.Version)
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Watching for service (%s:%s:%s) to scale up and enter ready state", ns, service.ID, service.Version), nil)
	watcher, err := i.kube.AppsV1().Deployments(ns).Watch(ctx, metav1.ListOptions{Watch: true, LabelSelector: labels, TimeoutSeconds: &timeout})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for ev := range watcher.ResultChan() {
		if ev.Type == watch.Error {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("service (%s:%s:%s) could not be scaled up", ns, service.ID, service.Version), nil, nil)
		}
		deployment := ev.Object.(*appsv1.Deployment)
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Received watch event for service (%s:%s): available replicas - %d; ready replicas - %d", ns, service.ID, deployment.Status.AvailableReplicas, deployment.Status.ReadyReplicas), nil)
		if deployment.Status.AvailableReplicas >= 1 && deployment.Status.ReadyReplicas >= 1 {
			return nil
		}
	}

	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("service (%s:%s) could not be started", ns, service.ID), nil, nil)
}

// ScaleUp is notifies keda to scale up a service
func (i *Istio) ScaleUp(ctx context.Context, projectID, serviceID, version string) error {
	if err := i.kedaScaler.ScaleUp(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to scale up service", err, map[string]interface{}{"project": projectID, "service": serviceID, "version": version})
	}

	return nil
}
