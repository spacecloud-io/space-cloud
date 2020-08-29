package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ApplyProjectConfig creates the config for the project
func (s *Manager) ApplyProjectConfig(ctx context.Context, project *config.Project, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.adminMan.ValidateProjectSyncOperation(s.projectConfig, project) {
		return http.StatusUpgradeRequired, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Upgrade your plan to create more projects", nil, nil)
	}

	// set default context time
	if project.ContextTimeGraphQL == 0 {
		project.ContextTimeGraphQL = 10
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
			p.AESKey = project.AESKey
			p.Secrets = project.Secrets
			p.SecretSource = project.SecretSource
			p.DockerRegistry = project.DockerRegistry
			p.ContextTimeGraphQL = project.ContextTimeGraphQL
			p.IsIntegration = project.IsIntegration
			// Mark project as existing
			doesProjectExists = true
			project = p
		}
	}

	if !doesProjectExists {
		// Append project with default modules to projects array
		project.Modules = &config.Modules{
			FileStore:    &config.FileStore{},
			Services:     &config.ServicesModule{},
			Auth:         map[string]*config.AuthStub{},
			Crud:         map[string]*config.CrudStub{},
			Routes:       []*config.Route{},
			GlobalRoutes: &config.GlobalRoutesConfig{},
			LetsEncrypt:  config.LetsEncrypt{WhitelistedDomains: []string{}},
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
	_ = s.modules.SetProjectConfig(project)

	if err := s.store.SetProject(ctx, project); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteProjectConfig applies delete project config command to the raft log
func (s *Manager) DeleteProjectConfig(ctx context.Context, projectID string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Delete project in the runner as well
	if s.runnerAddr != "" {
		if err := s.MakeHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("http://%s/v1/runner/%s", s.runnerAddr, projectID), token, "", "", &map[string]interface{}{}); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	s.delete(projectID)
	s.modules.Delete(projectID)

	if err := s.store.DeleteProject(ctx, projectID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetProjectConfig returns the config of specified project
func (s *Manager) GetProjectConfig(ctx context.Context, projectID string, params model.RequestParams) (int, []interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	v := []interface{}{}
	for _, p := range s.projectConfig.Projects {
		if projectID == "*" {
			// get all projects
			v = append(v, config.Project{DockerRegistry: p.DockerRegistry, AESKey: p.AESKey, ContextTimeGraphQL: p.ContextTimeGraphQL, Secrets: p.Secrets, SecretSource: p.SecretSource, IsIntegration: p.IsIntegration, Name: p.Name, ID: p.ID})
			continue
		}

		if projectID == p.ID {
			return http.StatusOK, []interface{}{config.Project{DockerRegistry: p.DockerRegistry, AESKey: p.AESKey, ContextTimeGraphQL: p.ContextTimeGraphQL, Secrets: p.Secrets, SecretSource: p.SecretSource, IsIntegration: p.IsIntegration, Name: p.Name, ID: p.ID}}, nil
		}
	}
	if len(v) > 0 || projectID == "*" {
		return http.StatusOK, v, nil
	}
	return http.StatusBadRequest, []interface{}{}, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Project (%s) not present in config", projectID), nil, nil)
}

// GetTokenForMissionControl returns the project token for internal use in mission control
func (s *Manager) GetTokenForMissionControl(ctx context.Context, projectID string, params model.RequestParams) (int, string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Get the auth module
	a := s.modules.GetAuthModuleForSyncMan()

	// Generate the token
	token, err := a.GetMissionControlToken(params.Claims)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}
