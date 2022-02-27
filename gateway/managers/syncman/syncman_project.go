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
func (s *Manager) ApplyProjectConfig(ctx context.Context, project *config.ProjectConfig, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

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
	p, ok := s.projectConfig.Projects[project.ID]
	if ok {
		p.ProjectConfig = project
	} else {
		s.projectConfig.Projects[project.ID] = config.GenerateEmptyProject(project)
		// Create a project in the runner as well
		if s.runnerAddr != "" {
			params := map[string]interface{}{"id": project.ID}
			if err := s.MakeHTTPRequest(ctx, "POST", fmt.Sprintf("http://%s/v1/runner/project/%s", s.runnerAddr, project.ID), token, "", params, &map[string]interface{}{}); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	if err := s.modules.SetProjectConfig(ctx, project); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetResource(ctx, config.GenerateResourceID(s.clusterID, project.ID, config.ResourceProject, project.ID), project); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteProjectConfig applies delete project config command to the raft log
func (s *Manager) DeleteProjectConfig(ctx context.Context, projectID string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Delete project in the runner as well
	if s.runnerAddr != "" {
		if err := s.MakeHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("http://%s/v1/runner/%s", s.runnerAddr, projectID), token, "", "", &map[string]interface{}{}); err != nil {
			return http.StatusInternalServerError, err
		}
	}
	// NOTE: we are not deleting project here as, the watcher of config maps will eventually delete the project
	s.modules.Delete(projectID)

	if err := s.store.DeleteProject(ctx, projectID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetProjectConfig returns the config of specified project
func (s *Manager) GetProjectConfig(ctx context.Context, projectID string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	v := []interface{}{}
	if projectID == "*" {
		for _, p := range s.projectConfig.Projects {
			if !p.ProjectConfig.IsIntegration {
				v = append(v, p.ProjectConfig)
			}
		}
		return http.StatusOK, v, nil
	}
	project, ok := s.projectConfig.Projects[projectID]
	if ok {
		return http.StatusOK, []interface{}{project.ProjectConfig}, nil
	}

	return http.StatusBadRequest, []interface{}{}, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Project (%s) not present in config", projectID), nil, nil)
}

// GetTokenForMissionControl returns the project token for internal use in mission control
func (s *Manager) GetTokenForMissionControl(ctx context.Context, projectID string, params model.RequestParams) (int, string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Get the auth module
	a, err := s.modules.GetAuthModuleForSyncMan(projectID)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	// Generate the token
	token, err := a.GetMissionControlToken(ctx, params.Claims)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}
