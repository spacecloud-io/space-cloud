package syncman

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// ConsulStore is an object for storing consul information
type ConsulStore struct {
	consulClient                     *api.Client
	nodeID, clusterID, advertiseAddr string
}

// GetAdminConfig returns the admin config
func (s *ConsulStore) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	kvPair, _, err := s.consulClient.KV().Get(fmt.Sprintf("sc/admin-config/%s", s.clusterID), &api.QueryOptions{})
	if err != nil {
		return nil, utils.LogError("couldn't fetch key in consul", "syncman", "getAdminConfig", err)
	}

	// kvPair will be nil if key doesn't exists
	if kvPair == nil {
		return nil, utils.LogError("specified key doesn't exists in consul", "syncman", "getAdminConfig", err)
	}

	var cluster *config.Admin
	if err := json.Unmarshal(kvPair.Value, &cluster); err != nil {
		return nil, utils.LogError("couldn't unmarshal cluster config to json", "syncman", "getAdminConfig", err)
	}

	return cluster, nil
}

// NewConsulStore creates new consul store
func NewConsulStore(nodeID, clusterID, advertiseAddr string) (*ConsulStore, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &ConsulStore{consulClient: client, nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr}, nil
}

// Register registers space cloud to the consul store
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

// WatchAdminConfig maintains consistency between all instances of sc
func (s *ConsulStore) WatchAdminConfig(cb func(cluster []*config.Admin)) error {
	watchParams := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "sc/admin-config/" + s.clusterID,
	}
	p, err := watch.Parse(watchParams)
	if err != nil {
		return err
	}

	p.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		kvPairs := data.(api.KVPairs)
		clusters := []*config.Admin{
			{
				LicenseKey:   "",
				LicenseValue: "",
				License:      "",
			},
		}

		for _, kv := range kvPairs {
			if err := json.Unmarshal(kv.Value, clusters[0]); err != nil {
				log.Println("Sync manager: Could not parse project received -", err)
				continue
			}
		}
		cb(clusters)
	}

	go func() {
		if err := p.Run(""); err != nil {
			log.Println("Sync Manager: could not start watcher -", err)
			os.Exit(-1)
		}
	}()
	return nil
}

// WatchProjects maintains consistency between all instances of sc
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

// WatchServices maintains consistency between all instances of sc
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

// SetProject sets the project of the consul store
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

// SetAdminConfig sets the admin config
func (s *ConsulStore) SetAdminConfig(ctx context.Context, cluster *config.Admin) error {
	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	data, _ := json.Marshal(cluster)

	_, err := s.consulClient.KV().Put(&api.KVPair{
		Key:   fmt.Sprintf("sc/admin-config/%s", s.clusterID),
		Value: data,
	}, opts)

	return err
}

// DeleteProject deletes the project from the consul store
func (s *ConsulStore) DeleteProject(ctx context.Context, projectID string) error {
	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	_, err := s.consulClient.KV().Delete(fmt.Sprintf("sc/projects/%s/%s", s.clusterID, projectID), opts)
	return err
}
