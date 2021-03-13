package syncman

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// KubeService is an object for storing kubeservice information
type KubeService struct {
	clusterID      string
	projectsConfig *config.Config
	kube           *kubernetes.Clientset
}

// NewKubeService creates a new Kube service
func NewKubeService(clusterID string) (*KubeService, error) {
	// Create the kubernetes client
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &KubeService{clusterID: clusterID, kube: kube, projectsConfig: new(config.Config)}, nil
}

func onAddOrUpdateServices(obj interface{}, services scServices) scServices {
	pod := obj.(*v1.Pod)
	id := string(pod.UID)

	// Ignore if pod isn't running
	if pod.Status.Phase != v1.PodRunning || pod.Status.PodIP == "" {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Pod (%s) isn't running yet. Current status - %s", id, pod.Status.Phase), nil)
		for index, service := range services {
			if service.id == id {
				helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Removing space cloud service from kubernetes", map[string]interface{}{"id": id})
				services[index] = services[len(services)-1]
				services = services[:len(services)-1]
				break
			}
		}
		return services
	}

	doesExist := false
	for _, service := range services {
		if service.id == id {
			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating space cloud service in kubernetes", map[string]interface{}{"id": id})
			doesExist = true
			break
		}
	}

	// add service if it doesn't exist
	if !doesExist {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Adding a space cloud service in kubernetes", map[string]interface{}{"id": id})
		services = append(services, &service{id: id})
	}
	return services
}

// WatchServices maintains consistency over all services
func (s *KubeService) WatchServices(cb func(scServices)) error {
	go func() {
		services := scServices{}
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("app=%s,clusterId=%s", "gateway", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 15*time.Minute, informers.WithTweakListOptions(options)).Core().V1().Pods().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash()

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				services = onAddOrUpdateServices(obj, services)
				sort.Stable(services)
				cb(services)
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				id := pod.Name
				for index, service := range services {
					if service.id == id {
						// remove service
						helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Removing space cloud service from kubernetes", map[string]interface{}{"id": id})
						services[index] = services[len(services)-1]
						services = services[:len(services)-1]
						break
					}
				}
				sort.Stable(services)
				cb(services)
			},
			UpdateFunc: func(old, obj interface{}) {
				services = onAddOrUpdateServices(obj, services)
				sort.Stable(services)
				cb(services)
			},
		})

		go informer.Run(stopper)
		<-stopper
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Stopped watching over services in kube store channel closed", nil)
	}()

	return nil
}
