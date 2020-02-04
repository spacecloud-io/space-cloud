package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

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
)

type KubeStore struct {
	clusterID string
	kube      *kubernetes.Clientset
}

const spaceCloud string = "space-cloud"

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

func (s *KubeStore) Register() {
	// kubernetes will handle this automatically
}

func onAddOrUpdateProjects(obj interface{}, projectMap map[string]*config.Project) map[string]*config.Project {
	configMap := obj.(*v1.ConfigMap)
	projectJsonString, ok := configMap.Data["project"]
	if !ok {
		logrus.Errorf("error watching projects in kube store unable to find project in config map")
		return nil
	}

	v := new(config.Project)
	if err := json.Unmarshal([]byte(projectJsonString), v); err != nil {
		logrus.Errorf("error while watching projects in kube store unable to unmarshal data - %v", err)
		return nil
	}
	projectMap[v.ID] = v
	return projectMap
}

func (s *KubeStore) WatchProjects(cb func(projects []*config.Project)) error {
	var options internalinterfaces.TweakListOptionsFunc
	// labels := fmt.Sprintf("clusterId=%s", s.clusterID)
	factory := informers.NewFilteredSharedInformerFactory(s.kube, 0, spaceCloud, options)
	informer := factory.Core().V1().ConfigMaps().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	projectMap := map[string]*config.Project{}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			projMap := onAddOrUpdateProjects(obj, projectMap)
			cb(s.getProjects(projMap))
		},
		DeleteFunc: func(obj interface{}) {
			configMap := obj.(*v1.ConfigMap)
			projectId, ok := configMap.Data["id"]
			if !ok {
				logrus.Errorf("error watching project in kube store unable to find project id in config map")
				return
			}
			delete(projectMap, projectId)
			cb(s.getProjects(projectMap))
		},
		UpdateFunc: func(old, obj interface{}) {
			projMap := onAddOrUpdateProjects(obj, projectMap)
			cb(s.getProjects(projMap))
		},
	})

	go informer.Run(stopper)

	return nil
}

func onAddOrUpdateServices(obj interface{}, services []*service) []*service {
	pod := obj.(*v1.Pod)
	id, ok := pod.Annotations["id"]
	if !ok {
		logrus.Errorf("error occurred watching services in kube store unable to find id in pod annotations while add event occurred")
		return nil
	}
	addr, ok := pod.Annotations["addr"]
	if !ok {
		logrus.Errorf("error occurred watching services in kube store unable to find addr in pod annotations while add event occurred")
		return nil
	}

	doesExist := false
	for _, service := range services {
		if service.id == id {
			doesExist = true
			service.addr = addr
			return nil
		}
	}

	// add service if it doesn't exist
	if !doesExist {
		services = append(services, &service{id: id, addr: addr})
	}
	return services
}

func (s *KubeStore) WatchServices(cb func(scServices)) error {
	services := scServices{}
	// labels := fmt.Sprintf("app=%s,clusterId=%s", "gateway", s.clusterID)
	var options internalinterfaces.TweakListOptionsFunc
	factory := informers.NewFilteredSharedInformerFactory(s.kube, 0, spaceCloud, options)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onAddOrUpdateServices(obj, services)
			sort.Stable(services)
			cb(services)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			id, ok := pod.Annotations["id"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find id in pod annotations while delete event occurred")
				return
			}
			for index, service := range services {
				if service.id == id {
					services[index] = services[len(services)-1]
					services = services[:len(services)-1]
					break
				}
			}
			sort.Stable(services)
			cb(services)
		},
		UpdateFunc: func(old, obj interface{}) {
			onAddOrUpdateServices(obj, services)
			sort.Stable(services)
			cb(services)
		},
	})

	go informer.Run(stopper)

	return nil
}

func (s *KubeStore) SetProject(ctx context.Context, project *config.Project) error {
	projectJsonString, err := json.Marshal(project)
	if err != nil {
		logrus.Errorf("error while setting project in kube store unable to marshal project config - %v", err)
		return err
	}

	configMap, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Get(fmt.Sprintf("%s-%s", s.clusterID, project.ID), v12.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		configMap := &v1.ConfigMap{
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"projectId": project.ID,
					"clusterId": s.clusterID,
				},
			},
			Data: map[string]string{
				"id":      project.ID,
				"project": string(projectJsonString),
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
	configMap.Data["project"] = string(projectJsonString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		logrus.Errorf("error while setting project in kube store unable to update config map - %v", err)
	}
	return err
}

func (s *KubeStore) DeleteProject(ctx context.Context, projectId string) error {
	err := s.kube.CoreV1().ConfigMaps(spaceCloud).Delete(fmt.Sprintf("%s-%s", s.clusterID, projectId), &v12.DeleteOptions{})
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
