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
	"k8s.io/client-go/tools/cache"
)

const (
	// resourceDeleteEvent istio uses this const to represent a delete event
	resourceDeleteEvent = "delete"
	// resourceUpdateEvent istio uses this const to represent a delete event
	resourceUpdateEvent = "update"
	// resourceAddEvent istio uses this const to represent a delete event
	resourceAddEvent = "add"
)

func onAddOrUpdateResource(eventType string, obj interface{}) (string, int32, int32, string, string) {
	deployment := obj.(*appsv1.Deployment)
	return eventType, deployment.Status.AvailableReplicas, deployment.Status.ReadyReplicas, deployment.Namespace, deployment.Labels["app"]
}

// WatchDeployments maintains consistency over all Deployment
func (i *Istio) WatchDeployments(cb func(eventType string, availableReplicas, readyReplicas int32, projectID, deploymentID string)) error {

	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = "app.kubernetes.io/managed-by=space-cloud"
		}

		informer := informers.NewSharedInformerFactoryWithOptions(i.kube, 15*time.Minute, informers.WithTweakListOptions(options)).Apps().V1().Deployments().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash() // handles a crash & logs an error

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				cb(onAddOrUpdateResource(resourceAddEvent, obj))
			},
			UpdateFunc: func(old, obj interface{}) {
				cb(onAddOrUpdateResource(resourceAddEvent, obj))
			},
			DeleteFunc: func(obj interface{}) {
				cb(onAddOrUpdateResource(resourceAddEvent, obj))
			},
		})

		go informer.Run(stopper)
		<-stopper
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Stopped watching over service", nil)
	}()
	return nil
}
