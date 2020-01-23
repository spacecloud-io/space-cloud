package syncman

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetEventSource returns the source id for the space cloud instance
func (s *Manager) GetEventSource() string {
	return fmt.Sprintf("sc-%s", s.nodeID)
}

// GetAssignedSpaceCloudURL returns the space cloud url assigned for the provided token
func (s *Manager) GetAssignedSpaceCloudURL(ctx context.Context, project string, token int) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.storeType == "none" {
		return fmt.Sprintf("http://localhost:%d/v1/api/%s/eventing/process", s.port, project), nil
	}

	index := calcIndex(token, utils.MaxEventTokens, len(s.services))

	return fmt.Sprintf("http://%s/v1/api/%s/eventing/process", s.services[index].addr, project), nil
}

// GetSpaceCloudNodeURLs returns the array of space cloud urls
func (s *Manager) GetSpaceCloudNodeURLs(project string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.storeType == "none" {
		return []string{fmt.Sprintf("http://localhost:%d/v1/api/%s/realtime/process", s.port, project)}
	}

	urls := make([]string, len(s.services))

	for i, svc := range s.services {
		urls[i] = fmt.Sprintf("http://%s/v1/api/%s/realtime/process", svc.addr, project)
	}

	return urls
}

func (s *Manager) GetRealtimeUrl(project string) string {
	return string(fmt.Sprintf("http://localhost:%d/v1/api/%s/realtime/handle", s.port, project))
}

// GetAssignedTokens returns the array or tokens assigned to this node
func (s *Manager) GetAssignedTokens() (start, end int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Always return true if running in single mode
	if s.storeType == "none" {
		return calcTokens(1, utils.MaxEventTokens, 0)
	}

	index := 0

	for i, v := range s.services {
		if v.id == s.nodeID {
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
	if s.storeType == "none" {
		return 1, nil
	}

	return len(s.services), nil
}

func (s *Manager) CreateProjectConfig(ctx context.Context, project *config.Project) (error, int) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	for _, p := range s.projectConfig.Projects {
		if p.ID == project.ID {
			return errors.New("project already exists in config"), http.StatusConflict
		}
	}

	s.projectConfig.Projects = append(s.projectConfig.Projects, project)

	// Create a project in the runner as well
	if s.runnerAddr != "" {
		params := map[string]interface{}{"id": project.ID}
		if err := s.MakeHTTPRequest(ctx, "POST", fmt.Sprintf("http://%s/v1/runner/project", s.runnerAddr), token, "", params, &map[string]interface{}{}); err != nil {
			return err, http.StatusInternalServerError
		}
	}

	// We will ignore the error for the create project request
	_ = s.cb(s.projectConfig)

	if s.storeType == "none" {
		return config.StoreConfigToFile(s.projectConfig, s.configFile), http.StatusInternalServerError
	}

	return s.store.SetProject(ctx, project), http.StatusInternalServerError
}

// SetProjectGlobalConfig applies the set project config command to the raft log
func (s *Manager) SetProjectGlobalConfig(ctx context.Context, project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project.ID)
	if err != nil {
		return err
	}

	projectConfig.Secret = project.Secret
	projectConfig.Name = project.Name

	return s.setProject(ctx, projectConfig)
}

// SetProjectConfig applies the set project config command to the raft log
func (s *Manager) SetProjectConfig(ctx context.Context, project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cb(s.projectConfig)
	return s.setProject(ctx, project)
}

func (s *Manager) setProject(ctx context.Context, project *config.Project) error {
	if err := s.cb(&config.Config{Projects: []*config.Project{project}}); err != nil {
		return err
	}

	s.setProjectConfig(project)

	if s.storeType == "none" {
		return config.StoreConfigToFile(s.projectConfig, s.configFile)
	}

	return s.store.SetProject(ctx, project)
}

// DeleteProjectConfig applies delete project config command to the raft log
func (s *Manager) DeleteProjectConfig(ctx context.Context, projectID string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	s.delete(projectID)
	if err := s.cb(s.projectConfig); err != nil {
		return err
	}

	if s.storeType == "none" {
		return config.StoreConfigToFile(s.projectConfig, s.configFile)
	}

	return s.store.DeleteProject(ctx, projectID)
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
