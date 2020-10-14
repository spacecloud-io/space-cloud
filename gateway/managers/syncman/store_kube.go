package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// KubeStore is an object for storing kubestore information
type KubeStore struct {
	clusterID      string
	projectsConfig *config.Config
	kube           *kubernetes.Clientset
}

const spaceCloud string = "space-cloud"

// NewKubeStore creates a new Kube store
func NewKubeStore(clusterID string) (*KubeStore, error) {
	// Create the kubernetes client
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &KubeStore{clusterID: clusterID, kube: kube, projectsConfig: new(config.Config)}, nil
}

// Register registers space cloud to the kube store
func (s *KubeStore) Register() {
	// kubernetes will handle this automatically
}

func onAddOrUpdateResource(eventType string, obj interface{}) (string, string, interface{}) {
	configMap := obj.(*v1.ConfigMap)
	resourceID, ok := configMap.Data["id"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (id) field was not found in config map data", eventType), nil, nil)
		return "", "", nil
	}
	dataJSONString, ok := configMap.Data["data"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (resource) field was not found in config map data", eventType), nil, nil)
		return "", "", nil
	}

	var v map[string]interface{}
	if err := json.Unmarshal([]byte(dataJSONString), &v); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal resource config map data while watching kube store project", nil, map[string]interface{}{"resourceId": resourceID, "eventType": eventType})
		return "", "", nil
	}
	return eventType, resourceID, v
}

// WatchResources maintains consistency over all projects
func (s *KubeStore) WatchResources(cb func(eventType, resourceID string, resource interface{})) error {
	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("clusterId=%s", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 0, informers.WithTweakListOptions(options)).Core().V1().ConfigMaps().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash() // handles a crash & logs an error

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				cb(onAddOrUpdateResource(config.ResourceAddEvent, obj))
			},
			UpdateFunc: func(old, obj interface{}) {
				cb(onAddOrUpdateResource(config.ResourceUpdateEvent, obj))
			},
			DeleteFunc: func(obj interface{}) {
				cb(onAddOrUpdateResource(config.ResourceDeleteEvent, obj))
			},
		})

		go informer.Run(stopper)
		<-stopper
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Stopped watching over projects in kube store", nil)
	}()
	return nil
}

func onAddOrUpdateServices(obj interface{}, services scServices) scServices {
	pod := obj.(*v1.Pod)
	id := pod.Name

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

	addr := fmt.Sprintf("%s.gateway.space-cloud.svc.cluster.local:4122", pod.Name)

	doesExist := false
	for _, service := range services {
		if service.id == id {
			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating space cloud service in kubernetes", map[string]interface{}{"id": id, "addr": addr})
			doesExist = true
			service.addr = addr
			break
		}
	}

	// add service if it doesn't exist
	if !doesExist {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Adding a space cloud service in kubernetes", map[string]interface{}{"id": id, "addr": addr})
		services = append(services, &service{id: id, addr: addr})
	}
	return services
}

// WatchServices maintains consistency over all services
func (s *KubeStore) WatchServices(cb func(scServices)) error {
	go func() {
		services := scServices{}
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("app=%s,clusterId=%s", "gateway", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 0, informers.WithTweakListOptions(options)).Core().V1().Pods().Informer()
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

// SetResource sets the project of the kube store
func (s *KubeStore) SetResource(ctx context.Context, resourceID string, resource interface{}) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting resource", map[string]interface{}{"resourceId": resourceID})
	clusterID, projectID, resourceType, err := splitResourceID(ctx, resourceID)
	if err != nil {
		return err
	}

	// validate if the resource value is according to the resource type
	if err := validateResource(ctx, config.ResourceAddEvent, s.projectsConfig, resourceID, resource); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to validate resource", err, map[string]interface{}{"project": projectID})
	}

	resourceJSONString, err := json.Marshal(resource)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't unmarshal project config", err, nil)
		return err
	}

	name := makeIDConfigMapCompatible(resourceID)
	configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(name, v12.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		configMap := &v1.ConfigMap{
			ObjectMeta: v12.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"kind":      string(resourceType),
					"projectId": projectID,
					"clusterId": clusterID,
				},
			},
			Data: map[string]string{
				"id":   resourceID,
				"data": string(resourceJSONString),
			},
		}
		_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Create(configMap)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't create config map", err, nil)
		}
		return err
	} else if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't set config map", err, nil)
		return err
	}

	configMap.Data["data"] = string(resourceJSONString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't update config map", err, nil)
	}
	return err
}

// DeleteResource deletes a resource from cluster
func (s *KubeStore) DeleteResource(ctx context.Context, resourceID string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Deleting resource", map[string]interface{}{"resourceId": resourceID})
	err := s.kube.CoreV1().ConfigMaps(spaceCloud).Delete(makeIDConfigMapCompatible(resourceID), &v12.DeleteOptions{})
	if kubeErrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete project in kube store couldn't get config map", err, nil)
	}
	return err
}

// GetGlobalConfig gets config of all resource required by a cluster
func (s *KubeStore) GetGlobalConfig() (*config.Config, error) {
	configMaps, err := s.kube.CoreV1().ConfigMaps(spaceCloud).List(v12.ListOptions{LabelSelector: fmt.Sprintf("clusterId=%s", s.clusterID)})
	if err != nil {
		return nil, err
	}
	globalConfig := config.GenerateEmptyConfig()
	for _, configMap := range configMaps.Items {
		eventType, resourceID, resource := onAddOrUpdateResource(config.ResourceAddEvent, &configMap)
		if err := validateResource(context.TODO(), eventType, globalConfig, resourceID, resource); err != nil {
			return nil, err
		}
	}
	s.projectsConfig = globalConfig
	return globalConfig, nil
}

// name cannot contain have underscore (_) but some resource have underscore in their name e.g --> event_logs, invocation_logs
// NOTE: use labels for getting the correct resource id
func makeIDConfigMapCompatible(resourceID string) string {
	return strings.ToLower(strings.Replace(resourceID, "_", "-", -1))
}
