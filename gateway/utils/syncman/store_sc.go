package syncman

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/sirupsen/logrus"
	sc "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"
	"github.com/spaceuptech/space-api-go/types"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

type SCStore struct {
	clusterID  string
	collection string
	scDatabase *db.DB
	kube       *kubernetes.Clientset
}

type scTableScheam struct {
	clusterId string `json:"cluster_id"`
	ProjectId string `json:"project_id"`
	Project   string `json:"project"`
}

func NewSCStore(clusterID, scStoreProjectName, scStoreAddr, scStoreDatabaseName string) (*SCStore, error) {
	// Create the kubernetes client
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &SCStore{
		clusterID:  clusterID,
		kube:       kube,
		scDatabase: sc.New(scStoreProjectName, scStoreAddr, false).DB(scStoreDatabaseName)}, nil
}

func (s *SCStore) Register() {
	// kubernetes will handle this automatically
}

func (s *SCStore) WatchProjects(cb func(projects []*config.Project)) error {
	subObj := s.scDatabase.LiveQuery(s.collection).Where(types.M{"cluster_id": s.clusterID}).Subscribe()
	for value := range subObj.C() {
		if err := value.Err(); err != nil {
			logrus.Errorf("error watching projects in sc store - %v", err)
			continue
		}
		projects, err := s.getProjects(subObj.GetSnapshot())
		if err != nil {
			logrus.Errorf("error watching projects in sc store unable to get projects - %v", err)
			continue
		}
		cb(projects)
	}
	return nil
}

func (s *SCStore) WatchServices(cb func(scServices)) error {
	// services := scServices{}
	// services = append(services, &service{id: "localhost", addr: "4122"})
	// sort.Stable(services)
	// cb(services)
	services := scServices{}
	labels := fmt.Sprintf("app=%s,clusterId=%s", "app", s.clusterID) // todo verify this
	watcher, err := s.kube.CoreV1().Pods("").Watch(v12.ListOptions{Watch: true, LabelSelector: labels})
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
				logrus.Errorf("error watching services id not found in pod annotations")
			}
			addr, ok := pod.Annotations["addr"]
			if !ok {
				logrus.Errorf("error watching services addr not found in pod annotations")
			}
			services = append(services, &service{id: id, addr: addr})
			sort.Stable(services)
			cb(services)

		case watch.Modified:
			pod := ee.Object.(*v1.Pod)
			id, ok := pod.Annotations["id"]
			if !ok {
				logrus.Errorf("error watching services id not found in pod annotations")
			}
			addr, ok := pod.Annotations["addr"]
			if !ok {
				logrus.Errorf("error watching services addr not found in pod annotations")
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
			if count == 0 {
				cb(append(services, &service{id: id, addr: addr}))
			}

		case watch.Deleted:
			pod := ee.Object.(*v1.Pod)
			id, ok := pod.Annotations["id"]
			if !ok {
				logrus.Errorf("error watching services id not found in pod annotations")
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

func (s *SCStore) SetProject(ctx context.Context, project *config.Project) error {
	res, err := s.scDatabase.Upsert(s.collection).Set(types.M{"project": project, "cluster_id": s.clusterID, "project_id": project.ID}).Apply(ctx)
	if res.Status != http.StatusOK || err != nil {
		return fmt.Errorf("error unable to set project in sc store -%v", err)
	}
	return nil
}

func (s *SCStore) DeleteProject(ctx context.Context, projectId string) error {
	res, err := s.scDatabase.Delete(s.collection).Where(types.M{"cluster_id": s.clusterID, "project_id": projectId}).Apply(ctx)
	if res.Status != http.StatusOK || err != nil {
		return fmt.Errorf("error unable to get project from sc store -%v", err)
	}
	return nil
}

func (s *SCStore) getProjects(v []types.DocumentSnapshot) ([]*config.Project, error) {
	projects := []*config.Project{}
	for _, value := range v {
		scScheam := new(scTableScheam)
		if err := value.Unmarshal(scScheam); err != nil {
			return nil, fmt.Errorf("error unable to marshal in sc store - %v", err)
		}
	}
	return projects, nil
}
