package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

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
	"github.com/spaceuptech/space-cloud/gateway/model"
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

func onAddOrUpdateResource(eventType string, obj interface{}) (string, string, config.Resource, interface{}) {
	configMap := obj.(*v1.ConfigMap)
	resourceID, ok := configMap.Data["id"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (id) field was not found in config map data", eventType), nil, nil)
		return "", "", "", nil
	}

	resourceType, ok := configMap.Labels["kind"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (kind) label was not found in config map", eventType), nil, nil)
		return "", "", "", nil
	}

	dataJSONString, ok := configMap.Data["data"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (resource) field was not found in config map data", eventType), nil, nil)
		return "", "", "", nil
	}

	v := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dataJSONString), &v); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal resource config map data while watching kube store project", nil, map[string]interface{}{"resourceId": resourceID, "eventType": eventType})
		return "", "", "", nil
	}
	return eventType, resourceID, config.Resource(resourceType), v
}

// WatchResources maintains consistency over all projects
func (s *KubeStore) WatchResources(cb func(eventType, resourceID string, resourceType config.Resource, resource interface{})) error {
	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("clusterId=%s", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 15*time.Minute, informers.WithTweakListOptions(options)).Core().V1().ConfigMaps().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash() // handles a crash & logs an error

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
				if resource == nil || resourceID == "" {
					return
				}
				cb(evenType, resourceID, resourceType, resource)
			},
			UpdateFunc: func(old, obj interface{}) {
				evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
				if resource == nil || resourceID == "" {
					return
				}
				cb(evenType, resourceID, resourceType, resource)
			},
			DeleteFunc: func(obj interface{}) {
				evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
				if resource == nil || resourceID == "" {
					return
				}
				cb(evenType, resourceID, resourceType, resource)
			},
		})

		go informer.Run(stopper)
		<-stopper
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Stopped watching over projects in kube store", nil)
	}()
	return nil
}

func onAddOrUpdateServices(obj interface{}, services model.ScServices) (string, model.ScServices) {
	pod := obj.(*v1.Pod)
	id := string(pod.UID)

	// Ignore if pod isn't running
	if pod.Status.Phase != v1.PodRunning || pod.Status.PodIP == "" {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Pod (%s) isn't running yet. Current status - %s", id, pod.Status.Phase), nil)
		for index, service := range services {
			if service.ID == id {
				helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Removing space cloud service from kubernetes", map[string]interface{}{"id": id})
				services[index] = services[len(services)-1]
				services = services[:len(services)-1]
				break
			}
		}
		return id, services
	}

	doesExist := false
	for _, service := range services {
		if service.ID == id {
			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating space cloud service in kubernetes", map[string]interface{}{"id": id})
			doesExist = true
			break
		}
	}

	// add service if it doesn't exist
	if !doesExist {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Adding a space cloud service in kubernetes", map[string]interface{}{"id": id})
		services = append(services, &model.Service{ID: id})
	}
	return id, services
}

// WatchServices maintains consistency over all services
func (s *KubeStore) WatchServices(cb func(string, string, model.ScServices)) error {
	go func() {
		services := model.ScServices{}
		var serviceID string
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("app=%s,clusterId=%s", "gateway", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 15*time.Minute, informers.WithTweakListOptions(options)).Core().V1().Pods().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash()

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				serviceID, services = onAddOrUpdateServices(obj, services)
				sort.Stable(services)
				cb(config.ResourceAddEvent, serviceID, services)
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				id := string(pod.UID)
				for index, service := range services {
					if service.ID == id {
						// remove service
						helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Removing space cloud service from kubernetes", map[string]interface{}{"id": id})
						services[index] = services[len(services)-1]
						services = services[:len(services)-1]
						break
					}
				}
				sort.Stable(services)
				cb(config.ResourceDeleteEvent, id, services)
			},
			UpdateFunc: func(old, obj interface{}) {
				serviceID, services = onAddOrUpdateServices(obj, services)
				sort.Stable(services)
				cb(config.ResourceUpdateEvent, serviceID, services)
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
	if err := updateResource(ctx, config.ResourceAddEvent, s.projectsConfig, resourceID, resourceType, resource); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to validate resource", err, map[string]interface{}{"project": projectID})
	}

	resourceJSONString, err := json.Marshal(resource)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't unmarshal project config", err, nil)
		return err
	}

	name := makeIDConfigMapCompatible(resourceID)
	configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(ctx, name, v12.GetOptions{})
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
		_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Create(ctx, configMap, v12.CreateOptions{})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't create config map", err, nil)
		}
		return err
	} else if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't set config map", err, nil)
		return err
	}

	configMap.Data["data"] = string(resourceJSONString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(ctx, configMap, v12.UpdateOptions{})
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't update config map", err, nil)
	}
	return err
}

// DeleteResource deletes a resource from cluster
func (s *KubeStore) DeleteResource(ctx context.Context, resourceID string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Deleting resource", map[string]interface{}{"resourceId": resourceID})
	err := s.kube.CoreV1().ConfigMaps(spaceCloud).Delete(ctx, makeIDConfigMapCompatible(resourceID), v12.DeleteOptions{})
	if kubeErrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete project in kube store couldn't get config map", err, nil)
	}
	return err
}

// DeleteProject deletes all the config resources which matches label projectId
func (s *KubeStore) DeleteProject(ctx context.Context, projectID string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Deleting entire project", map[string]interface{}{"projectId": projectID})
	// iterate over the map in reverse order
	for i := len(config.ResourceFetchingOrder) - 1; i >= 0; i-- {
		list, err := s.kube.CoreV1().ConfigMaps(spaceCloud).List(ctx, v12.ListOptions{LabelSelector: fmt.Sprintf("projectId=%s,kind=%s", projectID, config.ResourceFetchingOrder[i])})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list config maps which has label projectId = (%s)", projectID), err, nil)
		}
		for _, item := range list.Items {
			if err := s.DeleteResource(ctx, item.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetGlobalConfig gets config of all resource required by a cluster
func (s *KubeStore) GetGlobalConfig() (*config.Config, error) {
	globalConfig := config.GenerateEmptyConfig()
	for _, resourceType := range config.ResourceFetchingOrder {
		configMaps, err := s.kube.CoreV1().ConfigMaps(spaceCloud).List(context.TODO(), v12.ListOptions{LabelSelector: fmt.Sprintf("clusterId=%s,kind=%s", s.clusterID, resourceType)})
		if err != nil {
			return nil, err
		}
		for _, configMap := range configMaps.Items {
			eventType, resourceID, _, resource := onAddOrUpdateResource(config.ResourceAddEvent, &configMap)
			if err := updateResource(context.TODO(), eventType, globalConfig, resourceID, resourceType, resource); err != nil {
				return nil, err
			}
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
