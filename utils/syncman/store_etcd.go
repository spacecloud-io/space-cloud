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
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/hashicorp/consul/api"

	"github.com/spaceuptech/space-cloud/config"
)

type ETCDStore struct {
	etcdClient                       *clientv3.Client
	kv                               clientv3.KV
	nodeID, clusterID, advertiseAddr string
}

type trackedItemMeta struct {
	createRevision int64
	modRevision    int64
	service        *service
	project        *config.Project
}

func NewETCDStore(nodeID, clusterID, advertiseAddr, storeAddr string) (*ETCDStore, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{storeAddr},
	})
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

	go func() {
		for {
			select {
			case <-ticker.C:
				if _, err := s.etcdClient.KeepAlive(context.Background(), lease.ID); err != nil {
					log.Println("Could not renew consul session:", err)
					// register again
					s.Register()
				}
			}
		}
	}()
}

func (s *ETCDStore) WatchProjects(cb func(projects []*config.Project)) error {
	idxID := 3
	itemsMeta := map[string]*trackedItemMeta{}

	// Query all KVs with prefix
	res, err := s.etcdClient.Get(context.Background(), fmt.Sprintf("sc/projects/%s", s.clusterID), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range res.Kvs {
		// Get the id of the item
		id := strings.Split(string(kv.Key), "/")[idxID]
		project := new(config.Project)
		if err := json.Unmarshal(kv.Value, project); err != nil {
			log.Println("Sync manager: Could not parse project received -", err)
			continue
		}
		// Store the item
		itemsMeta[id] = &trackedItemMeta{createRevision: kv.CreateRevision, modRevision: kv.ModRevision, project: project}
	}

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/projects/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {

			for _, event := range watchResponse.Events {
				// Return if watch response contains an error
				if watchResponse.Err() != nil {
					log.Fatal(watchResponse.Err())
				}
				kv := event.Kv
				a := strings.Split(string(kv.Key), "/")
				id := a[idxID]
				if a[2] != s.clusterID {
					continue
				}

				switch event.Type {
				case mvccpb.PUT:
					project := new(config.Project)
					if err := json.Unmarshal(kv.Value, project); err != nil {
						log.Println("Sync manager: Could not parse project received -", err)
						continue
					}
					meta, p := itemsMeta[id]
					if !p {
						// AddStateless node if doesn't already exists
						itemsMeta[id] = &trackedItemMeta{createRevision: event.Kv.CreateRevision, modRevision: event.Kv.ModRevision, project: project}
					}

					// Ignore if incoming create revision is smaller
					if event.Kv.CreateRevision < meta.createRevision {
						break
					}

					// Update if incoming create revision or mod revision is greater
					if event.Kv.CreateRevision > meta.createRevision || event.Kv.ModRevision > meta.modRevision {
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						meta.project = project
						itemsMeta[id] = meta
						break
					}

				case mvccpb.DELETE:
					meta, p := itemsMeta[id]
					if !p {
						// Ignore if node does not exist
						break
					}

					// Remove if incoming mod revision is greater
					if event.Kv.ModRevision > meta.modRevision {
						// AddStateless node if doesn't already exists
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						delete(itemsMeta, id)
					}
				}
			}
			var arrProjects []*config.Project
			for _, item := range itemsMeta {
				arrProjects = append(arrProjects, item.project)
			}
			cb(arrProjects)
		}
	}()
	return nil
}

func (s *ETCDStore) WatchServices(cb func(scServices)) error {
	idxID := 3
	itemsMeta := map[string]*trackedItemMeta{}

	// Query all KVs with prefix
	res, err := s.etcdClient.Get(context.Background(), fmt.Sprintf("sc/instances/%s", s.clusterID), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range res.Kvs {
		// Get the id of the item
		id := strings.Split(string(kv.Key), "/")[idxID]
		service := &service{id: id, addr: string(kv.Value)}
		// Store the item
		itemsMeta[id] = &trackedItemMeta{createRevision: kv.CreateRevision, modRevision: kv.ModRevision, service: service}
	}

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/instances/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {
			for _, event := range watchResponse.Events {
				// Return if watch response contains an error
				if watchResponse.Err() != nil {
					log.Fatal(watchResponse.Err())
				}
				kv := event.Kv
				a := strings.Split(string(kv.Key), "/")
				id := a[idxID]
				if a[2] != s.clusterID {
					continue
				}

				switch event.Type {
				case mvccpb.PUT:
					meta, p := itemsMeta[id]
					if !p {
						// AddStateless node if doesn't already exists
						itemsMeta[id] = &trackedItemMeta{createRevision: event.Kv.CreateRevision, modRevision: event.Kv.ModRevision, service: &service{id: id, addr: string(kv.Value)}}
					}

					// Ignore if incoming create revision is smaller
					if event.Kv.CreateRevision < meta.createRevision {
						break
					}

					// Update if incoming create revision or mod revision is greater
					if event.Kv.CreateRevision > meta.createRevision || event.Kv.ModRevision > meta.modRevision {
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						meta.service = &service{id: id, addr: string(kv.Value)}
						itemsMeta[id] = meta
						break
					}

				case mvccpb.DELETE:
					meta, p := itemsMeta[id]
					if !p {
						// Ignore if node does not exist
						break
					}

					// Remove if incoming mod revision is greater
					if event.Kv.ModRevision > meta.modRevision {
						// AddStateless node if doesn't already exists
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						delete(itemsMeta, id)
					}
				}
			}

			// Sort and store
			var services scServices
			for _, item := range itemsMeta {
				services = append(services, item.service)
			}
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

	_, err := s.kv.Delete(ctx, fmt.Sprintf("sc/projects/%s/%s", s.clusterID, projectID))
	return err
}
