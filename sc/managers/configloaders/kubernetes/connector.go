package kube

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spaceuptech/helpers"
	"github.com/spf13/viper"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
)

// Connector is use as connector for KubeStore module
type Connector struct {
	Lock sync.RWMutex

	// Config related to cluster config
	clusterID      string
	ProjectsConfig *config.Config

	kube    *kubernetes.Clientset
	stopper chan struct{}
}

const spaceCloud string = "space-cloud"

// New create a new instance of the Module object
func New() (*Connector, error) {
	c := &Connector{}

	// Create the kubernetes client
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	// Set kube client and initialise stopper
	c.kube = kube
	c.stopper = make(chan struct{})

	// Set cluster-id
	c.clusterID = viper.GetString("cluster-id")

	// Set global config
	config, err := getGlobalConfig(c.kube, c.clusterID)
	if err != nil {
		return nil, err
	}
	c.ProjectsConfig = config

	if err := c.getSCConfig(); err != nil {
		return nil, err
	}

	return c, nil
}

// getSCConfig gets sc config from config map using shared informers
func (c *Connector) getSCConfig() error {
	// Start routine to observe space cloud project level resources
	if err := watchResources(c.kube, c.clusterID, c.stopper, func(eventType, resourceID string, resourceType config.Resource, resource interface{}) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_, projectID, _, err := splitResourceID(ctx, resourceID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to split resource id in watch resources", err, nil)
			return
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating resources", map[string]interface{}{"event": eventType, "resourceId": resourceID, "resource": resource, "projectId": projectID, "resourceType": resourceType})

		c.Lock.Lock()
		defer c.Lock.Unlock()
		if err := updateResource(ctx, eventType, c.ProjectsConfig, resourceID, resourceType, resource); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to update resources", err, nil)
			return
		}
	}); err != nil {
		return err
	}

	return nil
}

// Destruct destroys the kube store module
func (c *Connector) Destruct() error {
	// Acquire a lock
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.stopper <- struct{}{}

	return nil
}

// getGlobalConfig gets config of all resource required by a cluster
func getGlobalConfig(kube *kubernetes.Clientset, clusterID string) (*config.Config, error) {
	globalConfig := config.GenerateEmptyConfig()
	for _, resourceType := range config.ResourceFetchingOrder {
		configMaps, err := kube.CoreV1().ConfigMaps(spaceCloud).List(context.TODO(), v12.ListOptions{LabelSelector: fmt.Sprintf("clusterId=%s,kind=%s", clusterID, resourceType)})
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

	return globalConfig, nil
}

// watchResources get sc config from config map using shared informers
func watchResources(kube *kubernetes.Clientset, clusterID string, stopper chan struct{}, cb func(eventType, resourceID string, resourceType config.Resource, resource interface{})) error {
	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("clusterId=%s", clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(kube, 15*time.Minute, informers.WithTweakListOptions(options)).Core().V1().ConfigMaps().Informer()
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
