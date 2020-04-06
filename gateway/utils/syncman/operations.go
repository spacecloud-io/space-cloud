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

// GetRealtimeURL get the url of realtime
func (s *Manager) GetRealtimeURL(project string) string {
	return fmt.Sprintf("http://localhost:%d/v1/api/%s/realtime/handle", s.port, project)
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

// ApplyProjectConfig creates the config for the project
func (s *Manager) ApplyProjectConfig(ctx context.Context, project *config.Project) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// set default context time
	if project.ContextTime == 0 {
		project.ContextTime = 10
	}

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var doesProjectExists bool
	for _, p := range s.projectConfig.Projects {
		if p.ID == project.ID {
			// override the existing config
			p.Name = project.Name
			p.AESkey = project.AESkey
			p.Secrets = project.Secrets
			p.DockerRegistry = project.DockerRegistry
			p.ContextTime = project.ContextTime

			// Mark project as existing
			doesProjectExists = true
		}
	}

	if !doesProjectExists {
		// Append project with default modules to projects array
		project.Modules = &config.Modules{
			FileStore:   &config.FileStore{},
			Services:    &config.ServicesModule{},
			Auth:        map[string]*config.AuthStub{},
			Crud:        map[string]*config.CrudStub{},
			Routes:      []*config.Route{},
			LetsEncrypt: config.LetsEncrypt{WhitelistedDomains: []string{}},
		}
		s.projectConfig.Projects = append(s.projectConfig.Projects, project)

		// Create a project in the runner as well
		if s.runnerAddr != "" {
			params := map[string]interface{}{"id": project.ID}
			if err := s.MakeHTTPRequest(ctx, "POST", fmt.Sprintf("http://%s/v1/runner/project/%s", s.runnerAddr, project.ID), token, "", params, &map[string]interface{}{}); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// We will ignore the error for the create project request
	s.modules.SetProjectConfig(s.projectConfig, s.letsencrypt, s.routing)

	if s.storeType == "none" {
		return http.StatusInternalServerError, config.StoreConfigToFile(s.projectConfig, s.configFile)
	}

	return http.StatusInternalServerError, s.store.SetProject(ctx, project)
}

// SetProjectGlobalConfig applies the set project config command to the raft log
func (s *Manager) SetProjectGlobalConfig(ctx context.Context, project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.modules.SetGlobalConfig(project.Name, project.Secrets, project.AESkey); err != nil {
		return err
	}

	projectConfig, err := s.getConfigWithoutLock(project.ID)
	if err != nil {
		return err
	}

	projectConfig.Secrets = project.Secrets
	projectConfig.AESkey = project.AESkey
	projectConfig.Name = project.Name
	projectConfig.ContextTime = project.ContextTime

	return s.setProject(ctx, projectConfig)
}

// SetProjectConfig applies the set project config command to the raft log
func (s *Manager) SetProjectConfig(ctx context.Context, project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	s.modules.SetProjectConfig(s.projectConfig, s.letsencrypt, s.routing)

	return s.setProject(ctx, project)
}

func (s *Manager) setProject(ctx context.Context, project *config.Project) error {
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
	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return err
	}
	s.delete(projectID)

	// Create a project in the runner as well
	if s.runnerAddr != "" {
		if err := s.MakeHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("http://%s/v1/runner/%s", s.runnerAddr, projectID), token, "", "", &map[string]interface{}{}); err != nil {
			return err
		}
	}
	s.modules.SetProjectConfig(s.projectConfig, s.letsencrypt, s.routing)

	if s.storeType == "none" {
		return config.StoreConfigToFile(s.projectConfig, s.configFile)
	}

	return s.store.DeleteProject(ctx, projectID)
}

// GetProjectConfig returns the config of specified project
func (s *Manager) GetProjectConfig(projectID string) ([]interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	for _, p := range s.projectConfig.Projects {
		if projectID == p.ID {
			return []interface{}{config.Project{AESkey: p.AESkey, ContextTime: p.ContextTime, Secrets: p.Secrets, Name: p.Name, ID: p.ID}}, nil
		}
	}

	return []interface{}{}, errors.New("given project is not present in state")
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
