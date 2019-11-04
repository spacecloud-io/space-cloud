package syncman

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// GetEventSource returns the source id for the space cloud instance
func (s *Manager) GetEventSource() string {
	return fmt.Sprintf("sc-%s", s.nodeID)
}

// GetAssignedSpaceCloudURL returns the space cloud url assigned for the provided token
func (s *Manager) GetAssignedSpaceCloudURL(ctx context.Context, project string, token int) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if !s.isConsulEnabled {
		return fmt.Sprintf("http://localhost:4122/v1/api/%s/eventing/process", project), nil
	}

	opts := &api.QueryOptions{AllowStale: true}
	opts = opts.WithContext(ctx)

	index := calcIndex(token, utils.MaxEventTokens, len(s.services))

	return fmt.Sprintf("http://%s:%d/v1/api/%s/eventing/process", s.services[index].Node.Address, s.services[index].Service.Port, project), nil
}

// GetSpaceCloudNodeURLs returns the array of space cloud urls
func (s *Manager) GetSpaceCloudNodeURLs(project string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if !s.isConsulEnabled {
		return []string{fmt.Sprintf("http://localhost:4122/v1/api/%s/realtime/process", project)}
	}

	urls := make([]string, len(s.services))

	for i, addr := range s.services {
		urls[i] = fmt.Sprintf("http://%s:%d/v1/api/%s/realtime/process", addr.Node.Address, addr.Service.Port, project)
	}

	return urls
}

// GetAssignedTokens returns the array or tokens assigned to this node
func (s *Manager) GetAssignedTokens() (start, end int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Always return true if running in single mode
	if !s.isConsulEnabled {
		return calcTokens(1, utils.MaxEventTokens, 0)
	}

	index := 0

	for i, v := range s.services {
		if v.Service.ID == s.nodeID || (s.nodeID == v.Node.ID && s.port == v.Service.Port) {
			index = i
			break
		}
	}

	totalMembers := len(s.services)
	return calcTokens(totalMembers, utils.MaxEventTokens, index)
}

// GetClusterSize returns the size of the cluster
func (s *Manager) GetClusterSize(ctxParent context.Context) (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Return 1 if not running with consul
	if !s.isConsulEnabled {
		return 1, nil
	}

	return len(s.services), nil
}

func (s *Manager) CreateProjectConfig(project *config.Project) (error, int) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, p := range s.projectConfig.Projects {
		if p.ID == project.ID {
			return errors.New("project already exists in config"), http.StatusConflict
		}
	}

	s.projectConfig.Projects = append(s.projectConfig.Projects, project)

	s.cb(s.projectConfig)

	if !s.isConsulEnabled {
		return config.StoreConfigToFile(s.projectConfig, s.configFile), http.StatusInternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	data, _ := json.Marshal(project)

	_, err := s.consulClient.KV().Put(&api.KVPair{
		Key:   fmt.Sprintf("sc/projects/%s/%s", s.clusterID, project.ID),
		Value: data,
	}, opts)
	return err, http.StatusInternalServerError
}

// SetProjectGlobalConfig applies the set project config command to the raft log
func (s *Manager) SetProjectGlobalConfig(project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project.ID)
	if err != nil {
		return err
	}

	projectConfig.Secret = project.Secret
	projectConfig.Name = project.Name

	return s.setProject(projectConfig)
}

// SetProjectConfig applies the set project config command to the raft log
func (s *Manager) SetProjectConfig(project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.setProject(project)
}

func (s *Manager) setProject(project *config.Project) error {
	if err := s.cb(&config.Config{Projects: []*config.Project{project}}); err != nil {
		return err
	}

	s.setProjectConfig(project)

	if !s.isConsulEnabled {
		return config.StoreConfigToFile(s.projectConfig, s.configFile)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	data, _ := json.Marshal(project)

	_, err := s.consulClient.KV().Put(&api.KVPair{
		Key:   fmt.Sprintf("sc/projects/%s/%s", s.clusterID, project.ID),
		Value: data,
	}, opts)
	return err
}

// DeleteProjectConfig applies delete project config command to the raft log
func (s *Manager) DeleteProjectConfig(projectID string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.isConsulEnabled {
		s.delete(projectID)
		if err := s.cb(s.projectConfig); err != nil {
			return err
		}

		return config.StoreConfigToFile(s.projectConfig, s.configFile)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)

	_, err := s.consulClient.KV().Delete(fmt.Sprintf("sc/projects/%s/%s", s.clusterID, projectID), opts)
	return err
}

// GetConfig returns the config present in the state
func (s *Manager) GetConfig(projectID string) (*config.Project, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	for _, p := range s.projectConfig.Projects {
		if projectID == p.ID {
			return p, nil
		}
	}

	return nil, errors.New("given project is not present in state")
}
