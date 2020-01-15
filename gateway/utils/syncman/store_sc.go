package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
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

type scTableScheam struct {
	clusterId string          `json:"clusterId"`
	ProjectId string          `json:"projectId"`
	Project   *config.Project `json:"project"`
}

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
	labels := fmt.Sprintf("app=%s,clusterId=%s", "project", s.clusterID) // todo verify this
	watcher, err := s.kube.CoreV1().ConfigMaps(fmt.Sprintf("sc/projects/%s", s.clusterID)).Watch(v12.ListOptions{Watch: true, LabelSelector: labels})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	projectMap := map[string]scTableScheam{}
	for ee := range watcher.ResultChan() {
		switch ee.Type {
		case watch.Added, watch.Modified:
			configMap := ee.Object.(*v1.ConfigMap)
			for key, value := range configMap.Data {
				v := scTableScheam{}
				if err := json.Unmarshal([]byte(value), &v); err != nil {
					logrus.Errorf("error while watching projects in kube store unable to unmarshal data - %v", err)
					continue
				}
				projectMap[key] = v
			}
			cb(s.getProjects(projectMap))

		case watch.Deleted:
			configMap := ee.Object.(*v1.ConfigMap)
			for key := range configMap.Data {
				delete(projectMap, key)
			}
			cb(s.getProjects(projectMap))
		}
	}
	return nil
}

func (s *KubeStore) WatchServices(cb func(scServices)) error {
	services := scServices{}
	labels := fmt.Sprintf("app=%s,clusterId=%s", "service", s.clusterID) // todo verify this
	watcher, err := s.kube.CoreV1().Pods(fmt.Sprintf("sc/instaces/%s", s.clusterID)).Watch(v12.ListOptions{Watch: true, LabelSelector: labels})
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for ee := range watcher.ResultChan() {
		switch ee.Type {
		case watch.Added:
			pod := ee.Object.(*v1.Pod)
			id, ok := pod.Annotations["id"] // "verify all off these
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find id in pod annotations while add event occurred")
				break
			}
			addr, ok := pod.Annotations["addr"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find addr in pod annotations while add event occurred")
				break
			}
			services = append(services, &service{id: id, addr: addr})
			sort.Stable(services)
			cb(services)

		case watch.Modified:
			pod := ee.Object.(*v1.Pod)
			id, ok := pod.Annotations["id"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find id in pod annotations while modified event occurred")
				break
			}
			addr, ok := pod.Annotations["addr"]
			if !ok {
				logrus.Errorf("error occurred watching services in kube store unable to find addr in pod annotations while modified event occurred")
				break
			}
			count := 0
			for _, service := range services {
				if service.id == id {
					count++
					service.addr = addr
					sort.Stable(services)
					cb(services)
					break
				}
			}
			// if doesn't exit add to services
			if count == 0 {
				cb(append(services, &service{id: id, addr: addr}))
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
					sort.Stable(services)
					cb(services)
					break
				}
			}
		}
	}
	return nil
}

func (s *KubeStore) SetProject(ctx context.Context, project *config.Project) error {
	// TODO WHO WILL CREATE THIS CONFIG MAP
	configMap, err := s.kube.CoreV1().ConfigMaps(fmt.Sprintf("sc/projects/%s", s.clusterID)).Get("", v12.GetOptions{})
	if err != nil {
		return fmt.Errorf("error while setting project in kube store unable to get config map - %v", err)
	}
	data, err := json.Marshal(&scTableScheam{clusterId: s.clusterID, ProjectId: project.ID, Project: project})
	if err != nil {
		return fmt.Errorf("error while setting project in kube store unable to marshal data - %v", err)
	}
	configMap.Data[project.ID] = string(data)
	_, err = s.kube.CoreV1().ConfigMaps("").Update(configMap)
	if err != nil {
		return fmt.Errorf("error while setting project in kube store unable to update config map - %v", err)
	}
	return nil
}

func (s *KubeStore) DeleteProject(ctx context.Context, projectId string) error {
	// TODO WHO WILL CREATE THIS CONFIG MAP
	configMap, err := s.kube.CoreV1().ConfigMaps(fmt.Sprintf("sc/projects/%s", s.clusterID)).Get("", v12.GetOptions{})
	if err != nil {
		return fmt.Errorf("error while deleting project in kube store unable to get config map - %v", err)
	}
	delete(configMap.Data, projectId)
	_, err = s.kube.CoreV1().ConfigMaps("").Update(configMap)
	if err != nil {
		return fmt.Errorf("error while deleting project of kube store unable to update config map - %v", err)
	}
	return nil
}

func (s *KubeStore) getProjects(v map[string]scTableScheam) []*config.Project {
	projects := []*config.Project{}
	for _, value := range v {
		projects = append(projects, value.Project)
	}
	return projects
}
