package syncman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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
)

// KubeStore is an object for storing kubestore information
type KubeStore struct {
	clusterID string
	kube      *kubernetes.Clientset
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

	return &KubeStore{clusterID: clusterID, kube: kube}, nil
}

// Register registers space cloud to the kube store
func (s *KubeStore) Register() {
	// kubernetes will handle this automatically
}

func onAddOrUpdateAdminConfig(obj interface{}, clusters []*config.Admin) {
	configMap := obj.(*v1.ConfigMap)
	clusterJSONString, ok := configMap.Data["cluster"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to watch kube store admin config as field (cluster) not found in config map", nil, nil)
		return
	}

	if err := json.Unmarshal([]byte(clusterJSONString), clusters[0]); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal config map data while watching kube store admin config", nil, nil)
		return
	}
	if clusters[0].ClusterConfig == nil {
		clusters[0].ClusterConfig = getDefaultAdminConfig().ClusterConfig
	}
}

func onAddOrUpdateProjects(obj interface{}, projectMap map[string]*config.Project) map[string]*config.Project {
	configMap := obj.(*v1.ConfigMap)
	projectJSONString, ok := configMap.Data["project"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to watch kube store project as field (cluster) not found in config map", nil, nil)
		return nil
	}

	v := new(config.Project)
	if err := json.Unmarshal([]byte(projectJSONString), v); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal config map data while watching kube store project", nil, nil)
		return nil
	}
	projectMap[v.ID] = v
	return projectMap
}

// WatchAdminConfig maintains consistency over all projects
func (s *KubeStore) WatchAdminConfig(cb func(clusters []*config.Admin)) error {
	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("adminConfig=adminConfig,clusterId=%s", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 0, informers.WithTweakListOptions(options)).Core().V1().ConfigMaps().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash() // handles a crash & logs an error

		clusters := []*config.Admin{getDefaultAdminConfig()}

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onAddOrUpdateAdminConfig(obj, clusters)
				cb(clusters)
			},
			UpdateFunc: func(old, obj interface{}) {
				onAddOrUpdateAdminConfig(obj, clusters)
				cb(clusters)
			},
			DeleteFunc: func(old interface{}) {
				clusters[0] = getDefaultAdminConfig()
				cb(clusters)
			},
		})

		go informer.Run(stopper)
		<-stopper
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "stopped watching over projects in kube store", nil)
	}()
	return nil
}
func getDefaultAdminConfig() *config.Admin {
	return &config.Admin{
		ClusterConfig: &config.ClusterConfig{
			EnableTelemetry: true,
		},
		LicenseKey:   "",
		LicenseValue: "",
		License:      "",
	}
}

// WatchProjects maintains consistency over all projects
func (s *KubeStore) WatchProjects(cb func(projects []*config.Project)) error {
	go func() {
		var options internalinterfaces.TweakListOptionsFunc = func(options *v12.ListOptions) {
			options.LabelSelector = fmt.Sprintf("kind=project,clusterId=%s", s.clusterID)
		}
		informer := informers.NewSharedInformerFactoryWithOptions(s.kube, 0, informers.WithTweakListOptions(options)).Core().V1().ConfigMaps().Informer()
		stopper := make(chan struct{})
		defer close(stopper)
		defer runtime.HandleCrash() // handles a crash & logs an error

		projectMap := map[string]*config.Project{}

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				cb(s.getProjects(onAddOrUpdateProjects(obj, projectMap)))
			},
			DeleteFunc: func(obj interface{}) {
				configMap := obj.(*v1.ConfigMap)
				projectID, ok := configMap.Data["id"]
				if !ok {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete project while watching kube store projects as field (id) not present in config map", nil, nil)
					return
				}
				delete(projectMap, projectID)
				cb(s.getProjects(projectMap))

			},
			UpdateFunc: func(old, obj interface{}) {
				cb(s.getProjects(onAddOrUpdateProjects(obj, projectMap)))
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

// SetProject sets the project of the kube store
func (s *KubeStore) SetProject(ctx context.Context, project *config.Project) error {
	projectJSONString, err := json.Marshal(project)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't unmarshal project config", err, nil)
		return err
	}

	name := fmt.Sprintf("%s-%s", s.clusterID, project.ID)
	configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(name, v12.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		configMap := &v1.ConfigMap{
			ObjectMeta: v12.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"kind":      "project",
					"projectId": project.ID,
					"clusterId": s.clusterID,
				},
			},
			Data: map[string]string{
				"id":      project.ID,
				"project": string(projectJSONString),
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

	configMap.Data["id"] = project.ID
	configMap.Data["project"] = string(projectJSONString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project in kube store couldn't update config map", err, nil)
	}
	return err
}

// GetAdminConfig returns the admin config present in the store
func (s *KubeStore) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	name := fmt.Sprintf("sc-admin-config-%s", s.clusterID)

	for i := 0; i < 3; i++ {
		configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(name, v12.GetOptions{})
		if kubeErrors.IsNotFound(err) {
			return getDefaultAdminConfig(), nil
		} else if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to fetch config map (%s) from kubernetes", name), err, map[string]interface{}{"namespace": spaceCloud})

			// Sleep for 5 seconds then try again
			time.Sleep(5 * time.Second)
			continue
		}

		clusterJSONString, ok := configMap.Data["cluster"]
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Admin config data is corrupted", errors.New("key (cluster) not found in config map data object"), map[string]interface{}{})
		}

		cluster := new(config.Admin)
		if err := json.Unmarshal([]byte(clusterJSONString), cluster); err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Admin config data is corrupted", err, map[string]interface{}{})
		}

		return cluster, nil
	}

	return nil, errors.New("admin config could not be fetched")
}

// SetAdminConfig sets the project of the kube store
func (s *KubeStore) SetAdminConfig(ctx context.Context, cluster *config.Admin) error {
	clusterJSONString, err := json.Marshal(cluster)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set admin config in kube store couldn't unmarshal admin config data", err, nil)
		return err
	}

	name := fmt.Sprintf("sc-admin-config-%s", s.clusterID)
	configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(name, v12.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		configMap := &v1.ConfigMap{
			ObjectMeta: v12.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"adminConfig": "adminConfig",
					"clusterId":   s.clusterID,
				},
			},
			Data: map[string]string{
				"cluster": string(clusterJSONString),
			},
		}
		_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Create(configMap)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set admin config in kube store couldn't create config map", err, nil)
		}
		return err
	} else if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set admin config in kube store couldn't get config map", err, nil)
		return err
	}

	configMap.Data["cluster"] = string(clusterJSONString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set admin config in kube store couldn't update config map", err, nil)
	}
	return err
}

// DeleteProject deletes the project from the kube store
func (s *KubeStore) DeleteProject(ctx context.Context, projectID string) error {
	err := s.kube.CoreV1().ConfigMaps(spaceCloud).Delete(fmt.Sprintf("%s-%s", s.clusterID, projectID), &v12.DeleteOptions{})
	if kubeErrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete project in kube store couldn't get config map", err, nil)
	}
	return err
}

func (s *KubeStore) getProjects(v map[string]*config.Project) []*config.Project {
	projects := []*config.Project{}
	for _, value := range v {
		projects = append(projects, value)
	}
	return projects
}
