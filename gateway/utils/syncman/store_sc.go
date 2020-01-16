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
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

type KubeStore struct {
	clusterID string
	kube      *kubernetes.Clientset
}

const spaceCloud string = spaceCloud

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

func (s *KubeStore) WatchProjects(cb func(projects []*config.Project)) error {
	labels := fmt.Sprintf("clusterId=%s", s.clusterID)
	watcher, err := s.kube.CoreV1().ConfigMaps(spaceCloud).Watch(v12.ListOptions{Watch: true, LabelSelector: labels})
	if err != nil {
		logrus.Errorf("error watching projects in kube store - %v", err)
		return err
	}
	defer watcher.Stop()

	projectMap := map[string]*config.Project{}
	for ee := range watcher.ResultChan() {
		switch ee.Type {
		case watch.Added, watch.Modified:
			configMap := ee.Object.(*v1.ConfigMap)
			projectJsonString, ok := configMap.Data["project"]
			if !ok {
				logrus.Errorf("error watching projects in kube store unable to find project in config map")
				continue
			}

			v := new(config.Project)
			if err := json.Unmarshal([]byte(projectJsonString), v); err != nil {
				logrus.Errorf("error while watching projects in kube store unable to unmarshal data - %v", err)
				continue
			}
			projectMap[v.ID] = v

		case watch.Deleted:
			configMap := ee.Object.(*v1.ConfigMap)
			projectId, ok := configMap.Data["id"]
			if !ok {
				logrus.Errorf("error watching project in kube store unable to find project id in config map")
				continue
			}
			delete(projectMap, projectId)
		}

		cb(s.getProjects(projectMap))
	}
	return nil
}

func (s *KubeStore) WatchServices(cb func(scServices)) error {
	services := scServices{}
	labels := fmt.Sprintf("app=%s,clusterId=%s", "gateway", s.clusterID)
	watcher, err := s.kube.CoreV1().Pods(spaceCloud).Watch(v12.ListOptions{Watch: true, LabelSelector: labels})
	if err != nil {
		logrus.Errorf("error watching services in kube store - %v", err)
		return err
	}
	defer watcher.Stop()
	for ee := range watcher.ResultChan() {
		switch ee.Type {
		case watch.Added, watch.Modified:
			pod := ee.Object.(*v1.Pod)
			id, ok := pod.Annotations["id"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find id in pod annotations while add event occurred")
				break
			}
			addr, ok := pod.Annotations["addr"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find addr in pod annotations while add event occurred")
				break
			}

			doesExist := false
			for _, service := range services {
				if service.id == id {
					doesExist = true
					service.addr = addr
					break
				}
			}

			// add service if it doesn't exist
			if !doesExist {
				services = append(services, &service{id: id, addr: addr})
			}

		case watch.Deleted:
			pod := ee.Object.(*v1.Pod)
			id, ok := pod.Annotations["id"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find id in pod annotations while delete event occurred")
				break
			}
			for index, service := range services {
				if service.id == id {
					services[index] = services[len(services)-1]
					services = services[:len(services)-1]
					break
				}
			}
		}
		sort.Stable(services)
		cb(services)
	}
	return nil
}

func (s *KubeStore) SetProject(ctx context.Context, project *config.Project) error {
	projectJsonString, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("error while setting project in kube store unable to marshal project config - %v", err)
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
			return fmt.Errorf("error while setting project in kube store unable to create config map - %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error while setting project in kube store unable to get config map - %v", err)
	}

	configMap.Data["id"] = project.ID
	configMap.Data["project"] = string(projectJsonString)

	_, err = s.kube.CoreV1().ConfigMaps(spaceCloud).Update(configMap)
	if err != nil {
		return fmt.Errorf("error while setting project in kube store unable to update config map - %v", err)
	}
	return nil
}

func (s *KubeStore) DeleteProject(ctx context.Context, projectId string) error {
	err := s.kube.CoreV1().ConfigMaps(spaceCloud).Delete(fmt.Sprintf("%s-%s", s.clusterID, projectId), &v12.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("error while deleting project in kube store unable to get config map - %v", err)
	}
	return nil
}

func (s *KubeStore) getProjects(v map[string]*config.Project) []*config.Project {
	projects := []*config.Project{}
	for _, value := range v {
		projects = append(projects, value)
	}
	return projects
}
