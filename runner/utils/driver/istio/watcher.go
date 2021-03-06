package istio

import (
	"context"
	"time"

	"github.com/spaceuptech/helpers"
	appsv1 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	// ResourceDeleteEvent syncman uses this const to represent a delete event
	ResourceDeleteEvent = "delete"
	// ResourceUpdateEvent syncman uses this const to represent a delete event
	ResourceUpdateEvent = "update"
	// ResourceAddEvent syncman uses this const to represent a delete event
	ResourceAddEvent = "add"
)

func onAddOrUpdateResource(eventType string, obj interface{}) (string, int32, int32, string, string) {
	deployment := obj.(*appsv1.Deployment)
	//helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Received watch event for service (%s:%s): available replicas - %d; ready replicas - %d", ns, service.ID, deployment.Status.AvailableReplicas, deployment.Status.ReadyReplicas), nil)
	return eventType, deployment.Status.AvailableReplicas, deployment.Status.ReadyReplicas, deployment.Namespace, deployment.Labels["app"]
}

// WatchDeployments maintains consistency over all Deployment
func WatchDeployments(cb func(eventType string, availableReplicas, readyReplicas int32, projectID, deploymentID string)) error {

	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = "app.kubernetes.io/managed-by=space-cloud"
		}
		kube := new(kubernetes.Clientset)
		informer := informers.NewSharedInformerFactoryWithOptions(kube, 15*time.Minute, informers.WithTweakListOptions(options)).Core().V1().Deployment().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash() // handles a crash & logs an error

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				evenType, availableReplicas, readyReplicas, projectID, deploymentID := onAddOrUpdateResource(ResourceAddEvent, obj)
				cb(evenType, availableReplicas, readyReplicas, projectID, deploymentID)
			},
			UpdateFunc: func(old, obj interface{}) {
				evenType, availableReplicas, readyReplicas, projectID, deploymentID := onAddOrUpdateResource(ResourceUpdateEvent, obj)
				cb(evenType, availableReplicas, readyReplicas, projectID, deploymentID)
			},
			DeleteFunc: func(obj interface{}) {
				evenType, availableReplicas, readyReplicas, projectID, deploymentID := onAddOrUpdateResource(ResourceDeleteEvent, obj)
				cb(evenType, availableReplicas, readyReplicas, projectID, deploymentID)
			},
		})

		go informer.Run(stopper)
		<-stopper
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Stopped watching over service", nil)
	}()
	return nil
}
