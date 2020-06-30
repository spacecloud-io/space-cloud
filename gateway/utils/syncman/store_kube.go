package syncman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
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
	"github.com/spaceuptech/space-cloud/gateway/utils"
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
		logrus.Errorf("error watching projects in kube store unable to find field project in config map")
		return
	}

	if err := json.Unmarshal([]byte(clusterJSONString), clusters[0]); err != nil {
		logrus.Errorf("error while watching projects in kube store unable to unmarshal data - %v", err)
		return
	}
}

func onAddOrUpdateProjects(obj interface{}, projectMap map[string]*config.Project) map[string]*config.Project {
	configMap := obj.(*v1.ConfigMap)
	projectJSONString, ok := configMap.Data["project"]
	if !ok {
		logrus.Errorf("error watching projects in kube store unable to find field project in config map")
		return nil
	}

	v := new(config.Project)
	if err := json.Unmarshal([]byte(projectJSONString), v); err != nil {
		logrus.Errorf("error while watching projects in kube store unable to unmarshal data - %v", err)
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

		clusters := []*config.Admin{
			{
				ClusterID:  "",
				ClusterKey: "",
				License:    "",
			},
		}

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onAddOrUpdateAdminConfig(obj, clusters)
				cb(clusters)
			},
			UpdateFunc: func(old, obj interface{}) {
				onAddOrUpdateAdminConfig(obj, clusters)
				cb(clusters)
			},
		})

		go informer.Run(stopper)
		<-stopper
		logrus.Debug("stopped watching over projects in kube store")
	}()
	return nil
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
					logrus.Errorf("error watching project in kube store unable to find project id in config map")
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
		logrus.Debug("stopped watching over projects in kube store")
	}()
	return nil
}

func onAddOrUpdateServices(obj interface{}, services scServices) scServices {
	pod := obj.(*v1.Pod)
	id := pod.Name

	// Ignore if pod isn't running
	if pod.Status.Phase != v1.PodRunning || pod.Status.PodIP == "" {
		logrus.Debugf("Pod (%s) isn't running yet. Current status - %s", id, pod.Status.Phase)
		for index, service := range services {
			if service.id == id {
				utils.LogDebug("Removing service", "syncman", "onAddOrUpdateServices", map[string]interface{}{"id": id})
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
			utils.LogDebug("Updating service", "syncman", "onAddOrUpdateServices", map[string]interface{}{"id": id, "addr": addr})
			doesExist = true
			service.addr = addr
			break
		}
	}

	// add service if it doesn't exist
	if !doesExist {
		utils.LogDebug("Adding service", "syncman", "onAddOrUpdateServices", map[string]interface{}{"id": id, "addr": addr})
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
						utils.LogDebug("Removing service", "syncman", "WatchServices", map[string]interface{}{"id": id})
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
		logrus.Debug("stopped watching over services in kube store channel closed")
	}()

	return nil
}

// SetProject sets the project of the kube store
func (s *KubeStore) SetProject(ctx context.Context, project *config.Project) error {
	projectJSONString, err := json.Marshal(project)
	if err != nil {
		logrus.Errorf("error while setting project in kube store unable to marshal project config - %v", err)
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
			logrus.Errorf("error while setting project in kube store unable to create config map - %v", err)
		}
		return err
	} else if err != nil {
		logrus.Errorf("error while setting project in kube store unable to get config map - %v", err)
		return err
	}

	configMap.Data["id"] = project.ID
	configMap.Data["project"] = string(projectJSONString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		logrus.Errorf("error while setting project in kube store unable to update config map - %v", err)
	}
	return err
}

// GetAdminConfig returns the admin config present in the store
func (s *KubeStore) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	name := fmt.Sprintf("sc-admin-config-%s", s.clusterID)

	for i := 0; i < 3; i++ {
		configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(name, v12.GetOptions{})
		if kubeErrors.IsNotFound(err) {
			return new(config.Admin), nil
		} else if err != nil {
			_ = utils.LogError("Unable to fetch admin config", "syncman", "GetAdminConfig", err)

			// Sleep for 5 seconds then try again
			time.Sleep(5 * time.Second)
			continue
		}

		clusterJSONString, ok := configMap.Data["cluster"]
		if !ok {
			return nil, utils.LogError("Admin config data is corrupted", "syncman", "GetAdminConfig", errors.New("key not found in config map"))
		}

		cluster := new(config.Admin)
		if err := json.Unmarshal([]byte(clusterJSONString), cluster); err != nil {
			return nil, utils.LogError("Admin config data is corrupted", "syncman", "GetAdminConfig", err)
		}

		return cluster, nil
	}

	return nil, errors.New("admin config could not be fetched")
}

// SetAdminConfig sets the project of the kube store
func (s *KubeStore) SetAdminConfig(ctx context.Context, cluster *config.Admin) error {
	clusterJSONString, err := json.Marshal(cluster)
	if err != nil {
		logrus.Errorf("error while setting project in kube store unable to marshal project config - %v", err)
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
			logrus.Errorf("error while setting project in kube store unable to create config map - %v", err)
		}
		return err
	} else if err != nil {
		logrus.Errorf("error while setting project in kube store unable to get config map - %v", err)
		return err
	}

	configMap.Data["cluster"] = string(clusterJSONString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		logrus.Errorf("error while setting project in kube store unable to update config map - %v", err)
	}
	return err
}

// DeleteProject deletes the project from the kube store
func (s *KubeStore) DeleteProject(ctx context.Context, projectID string) error {
	err := s.kube.CoreV1().ConfigMaps(spaceCloud).Delete(fmt.Sprintf("%s-%s", s.clusterID, projectID), &v12.DeleteOptions{})
	if kubeErrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		logrus.Errorf("error while deleting project in kube store unable to get config map - %v", err)
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
