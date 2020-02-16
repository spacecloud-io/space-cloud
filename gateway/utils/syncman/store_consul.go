package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

type ConsulStore struct {
	consulClient                     *api.Client
	nodeID, clusterID, advertiseAddr string
}

func NewConsulStore(nodeID, clusterID, advertiseAddr string) (*ConsulStore, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &ConsulStore{consulClient: client, nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr}, nil
}

func (s *ConsulStore) Register() {
	opts := &api.WriteOptions{}
	opts = opts.WithContext(context.Background())

	session := s.consulClient.Session()
	id, _, err := session.Create(&api.SessionEntry{
		Name:     s.nodeID,
		Behavior: "delete",
		TTL:      "10s",
	}, opts)
	if err != nil {
		log.Fatal("Could not create a new session with consul:", err)
	}

	data := []byte(s.advertiseAddr)
	if _, _, err := s.consulClient.KV().Acquire(&api.KVPair{Session: id, Key: fmt.Sprintf("sc/instances/%s/%s", s.clusterID, s.nodeID), Value: data}, opts); err != nil {
		log.Fatal("Could not register space cloud with consul:", err)
	}

	ticker := time.NewTicker(4 * time.Second)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if _, _, err := session.Renew(id, opts); err != nil {
				log.Println("Could not renew consul session:", err)
				// register again
				s.Register()
				return
			}
		}
	}()
}

func (s *ConsulStore) WatchProjects(cb func(projects []*config.Project)) error {
	watchParams := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "sc/projects/" + s.clusterID,
	}
	p, err := watch.Parse(watchParams)
	if err != nil {
		return err
	}

	p.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		kvPairs := data.(api.KVPairs)
		var projects []*config.Project

		for _, kv := range kvPairs {
			a := strings.Split(kv.Key, "/")
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

	go func() {
		if err := p.Run(""); err != nil {
			log.Println("Sync Manager: could not start watcher -", err)
			os.Exit(-1)
		}
	}()
	return nil
}

func (s *ConsulStore) WatchServices(cb func(scServices)) error {
	watchParams := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "sc/instances/" + s.clusterID,
	}
	p, err := watch.Parse(watchParams)
	if err != nil {
		return err
	}

	p.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		kvPairs := data.(api.KVPairs)

		var services scServices

		// Filter out failing nodes
		for _, kv := range kvPairs {
			a := strings.Split(kv.Key, "/")
			if a[2] != s.clusterID {
				continue
			}

			service := new(service)
			service.id = a[3]
			service.addr = string(kv.Value)
			services = append(services, service)
		}

		// Sort and store
		sort.Stable(services)
		cb(services)
	}

	go func() {
		if err := p.Run(""); err != nil {
			log.Println("Sync Manager: could not start watch -", err)
			os.Exit(-1)
		}
	}()

	return nil
}

func (s *ConsulStore) SetProject(ctx context.Context, project *config.Project) error {
	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	data, _ := json.Marshal(project)

	_, err := s.consulClient.KV().Put(&api.KVPair{
		Key:   fmt.Sprintf("sc/projects/%s/%s", s.clusterID, project.ID),
		Value: data,
	}, opts)

	return err
}

func (s *ConsulStore) DeleteProject(ctx context.Context, projectID string) error {
	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	_, err := s.consulClient.KV().Delete(fmt.Sprintf("sc/projects/%s/%s", s.clusterID, projectID), opts)
	return err
}
