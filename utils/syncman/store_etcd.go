package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/hashicorp/consul/api"

	"github.com/spaceuptech/space-cloud/config"
)

type ETCDStore struct {
	etcdClient                       *clientv3.Client
	kv                               clientv3.KV
	nodeID, clusterID, advertiseAddr string
}

func NewETCDStore(nodeID, clusterID, advertiseAddr string) (*ETCDStore, error) {
	client, err := clientv3.New(clientv3.Config{})
	if err != nil {
		return nil, err
	}

	return &ETCDStore{etcdClient: client, nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, kv: clientv3.NewKV(client)}, nil
}

func (s *ETCDStore) Register() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	lease, err := s.etcdClient.Grant(ctx, 10)
	if err != nil {
		log.Fatal("Could not create a new session with etcd:", err)
	}

	if _, err := s.kv.Put(ctx, fmt.Sprintf("sc/instances/%s/%s", s.clusterID, s.nodeID), s.advertiseAddr, clientv3.WithLease(lease.ID)); err != nil {
		log.Fatal("Could not register space cloud with etcd:", err)
	}

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if _, err := s.etcdClient.KeepAlive(context.Background(), lease.ID); err != nil {
					log.Println("Could not renew consul session:", err)
				}
			}
		}
	}()
}

func (s *ETCDStore) WatchProjects(cb func(projects []*config.Project)) error {

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/projects/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {
			var projects []*config.Project
			for _, event := range watchResponse.Events {
				kv := event.Kv
				a := strings.Split(string(kv.Key), "/")
				if a[2] != s.clusterID {
					continue
				}

				project := new(config.Project)
				if err := json.Unmarshal(kv.Value, project); err != nil {
					log.Println("Sync manager: Could not parse project received -", err)
					continue
				}

				projects = append(projects, project)
			}
			cb(projects)
		}
	}()

	return nil
}

func (s *ETCDStore) WatchServices(cb func(scServices)) error {

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/instances/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {
			var services scServices
			for _, event := range watchResponse.Events {
				kv := event.Kv
				a := strings.Split(string(kv.Key), "/")
				if a[2] != s.clusterID {
					continue
				}

				service := new(service)
				service.id = string(kv.Key)
				service.addr = string(kv.Value)

				services = append(services, service)
			}

			// Sort and store
			sort.Stable(services)
			cb(services)
		}
	}()

	return nil
}

func (s *ETCDStore) SetProject(project *config.Project) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	_, err := s.kv.Put(ctx, fmt.Sprintf("sc/projects/%s/%s", s.clusterID, project.ID), project.ID)

	return err
}

func (s *ETCDStore) DeleteProject(projectID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	_,err := s.kv.Delete(ctx, fmt.Sprintf("sc/projects/%s/%s", s.clusterID, projectID))
	return err
}
