package syncman

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/integration"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (s *Manager) EnableIntegration(ctx context.Context, integrationConfig *config.IntegrationConfig, params model.RequestParams) (int, error) {
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

	s.lock.Lock()
	defer s.lock.Unlock()
	resourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegration, integrationConfig.ID)
	if err := s.adminMan.ValidateIntegrationSyncOperation(config.Integrations{resourceID: integrationConfig}); err != nil {
		return http.StatusUpgradeRequired, err
	}

	// Create a project if it doesn't already exist
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Creating a new project for integration (%s)", integrationConfig.ID), nil)

	// Seed the random byte generator
	rand.Seed(time.Now().UnixNano())

	// We need runner enabled for this one. We need to create a new project, secret and deployment in runner
	if s.runnerAddr == "" {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Runner must be enabled for integrations to work", nil, nil)
	}

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Instruct runner to create project
	reqPayload := map[string]interface{}{"id": integrationConfig.ID, "kind": "integration"}
	if err := s.MakeHTTPRequest(ctx, "POST", fmt.Sprintf("http://%s/v1/runner/project/%s", s.runnerAddr, integrationConfig.ID), token, "", reqPayload, &map[string]interface{}{}); err != nil {
		return http.StatusInternalServerError, err
	}

	// Create login credentials for the token
	integrationKeyTemp := make([]byte, 32)
	_, _ = rand.Read(integrationKeyTemp)
	integrationKey := base64.StdEncoding.EncodeToString(integrationKeyTemp)

	// Instruct runner to create secret with appropriate login credentials
	reqPayload = map[string]interface{}{"id": integrationConfig.ID, "type": "env", "data": map[string]string{"PROJECT": integrationConfig.ID, "LOGIN_ID": integrationConfig.ID, "LOGIN_KEY": integrationKey}}
	if err := s.MakeHTTPRequest(ctx, "POST", fmt.Sprintf("http://%s/v1/runner/%s/secrets/%s", s.runnerAddr, integrationConfig.ID, integrationConfig.ID), token, "", reqPayload, &map[string]interface{}{}); err != nil {
		return http.StatusInternalServerError, err
	}

	// Attach the integrationKey to the integration config
	integrationConfig.Key = integrationKey

	// Instruct runner to create deployment
	license, err := s.adminMan.ParseLicense(integrationConfig.License)
	if err != nil {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to parse integration license", err, nil)
	}
	integrationConfig.Deployments = license["deployments"].([]interface{})
	for _, service := range integrationConfig.Deployments {
		obj := service.(map[string]interface{})
		if err := s.MakeHTTPRequest(ctx, "POST", fmt.Sprintf("http://%s/v1/runner/%s/services/%s/%s", s.runnerAddr, integrationConfig.ID, obj["id"], obj["version"]), token, "", obj, &map[string]interface{}{}); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Generate AES key for integration project
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to create new aes key for integration", nil, nil)
	}

	// Check if integration already exists
	integrations := make(config.Integrations)
	var proj *config.ProjectConfig
	if _, p := s.projectConfig.Integrations.Get(integrationConfig.ID); !p {
		// Prepare a unique
		proj = &config.ProjectConfig{
			ID:                 integrationConfig.ID,
			Name:               integrationConfig.Name,
			SecretSource:       integrationConfig.SecretSource,
			ContextTimeGraphQL: 20,
			AESKey:             base64.StdEncoding.EncodeToString(key),
			IsIntegration:      true,
		}

		integrations[resourceID] = integrationConfig
	} else {
		integrations = s.projectConfig.Integrations
		// update existing integration
		integrations[resourceID] = integrationConfig
	}

	if err := s.integrationMan.SetIntegrations(integrations); err != nil {
		return http.StatusUpgradeRequired, err
	}
	s.adminMan.SetIntegrationConfig(integrations)
	s.projectConfig.Integrations = integrations

	if err := s.store.SetResource(ctx, resourceID, integrationConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	if proj != nil {
		s.projectConfig.Projects[proj.ID] = &config.Project{ProjectConfig: proj}
		if err := s.modules.SetProjectConfig(ctx, proj); err != nil {
			return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to create new project for integration", nil, nil)
		}

		// Update the store
		rID := config.GenerateResourceID(s.clusterID, proj.ID, config.ResourceProject, proj.ID)
		if err := s.store.SetResource(ctx, rID, proj); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

// RemoveIntegration removes an integration from space cloud
func (s *Manager) RemoveIntegration(ctx context.Context, id string, params model.RequestParams) (int, error) {
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

	s.lock.Lock()
	defer s.lock.Unlock()

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Delete project in runner
	if err := s.MakeHTTPRequest(ctx, http.MethodDelete, fmt.Sprintf("http://%s/v1/runner/%s", s.runnerAddr, id), token, "", &map[string]interface{}{}, &map[string]interface{}{}); err != nil {
		return http.StatusInternalServerError, err
	}

	// Remove project from projects array
	projectResourceID := config.GenerateResourceID(s.clusterID, id, config.ResourceProject, id)
	delete(s.projectConfig.Projects, projectResourceID)

	// Remove integration from integrations array
	integrationResourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegration, id)
	delete(s.projectConfig.Integrations, integrationResourceID)

	// Remove integration hooks
	for _, hook := range s.projectConfig.IntegrationHooks {
		// remove integration hook belonging to particular integration
		if hook.IntegrationID == id {
			resourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegrationHook, hook.ID)
			delete(s.projectConfig.IntegrationHooks, resourceID)
			if err := s.store.DeleteResource(ctx, resourceID); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// Update the modules and integration manager
	s.modules.Delete(id)
	_ = s.integrationMan.SetConfig(s.projectConfig.Integrations, s.projectConfig.IntegrationHooks)

	// Update the stores
	if err := s.store.DeleteProject(ctx, id); err != nil {
		return http.StatusInternalServerError, err
	}
	if err := s.store.DeleteResource(ctx, integrationResourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetIntegrations returns the list of integrations from gateway
func (s *Manager) GetIntegrations(ctx context.Context, id string, params model.RequestParams) (int, []interface{}, error) {
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

	result := make([]interface{}, 0)
	for _, i := range s.projectConfig.Integrations {
		if id == "*" {
			result = append(result, i)
			continue
		}

		if id == i.ID {
			result = append(result, i)
		}
	}

	if len(result) > 0 || id == "*" {
		return http.StatusOK, result, nil
	}

	return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) not found", id), nil, nil)
}

// AddIntegrationHook adds an integration hook
func (s *Manager) AddIntegrationHook(ctx context.Context, integrationID string, hookConfig *config.IntegrationHook, params model.RequestParams) (int, error) {
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

	s.lock.Lock()
	defer s.lock.Unlock()

	// Find the right integration
	resourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegration, integrationID)
	integrationConfig, ok := s.projectConfig.Integrations[resourceID]
	if !ok || integrationConfig == nil {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) does not exist", integrationID), nil, nil)
	}

	// Check if the integration has the permissions to create this hook
	if !integration.HasPermissionForHook(integrationConfig, hookConfig) {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) does not have necessary permissions to create the hook", integrationID), nil, nil)
	}

	// Create an empty map if nil
	if s.projectConfig.IntegrationHooks == nil {
		s.projectConfig.IntegrationHooks = make(config.IntegrationHooks)
	}

	// Add the hook and store the config
	resourceID = config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegrationHook, hookConfig.ID)
	s.projectConfig.IntegrationHooks[resourceID] = hookConfig
	s.integrationMan.SetIntegrationHooks(s.projectConfig.IntegrationHooks)

	// Store the config in the store
	if err := s.store.SetResource(ctx, resourceID, hookConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemoveIntegrationHook removes an integration hook
func (s *Manager) RemoveIntegrationHook(ctx context.Context, integrationID, hookID string, params model.RequestParams) (int, error) {
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

	s.lock.Lock()
	defer s.lock.Unlock()

	// Find the right integration
	resourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegration, integrationID)
	integrationConfig, ok := s.projectConfig.Integrations[resourceID]
	if !ok || integrationConfig == nil {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) does not exist", integrationID), nil, nil)
	}

	// Delete the hook
	resourceID = config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegrationHook, hookID)
	delete(s.projectConfig.IntegrationHooks, resourceID)
	s.integrationMan.SetIntegrationHooks(s.projectConfig.IntegrationHooks)

	// Store the config in the store
	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetIntegrationHook removes an integration hook
func (s *Manager) GetIntegrationHooks(ctx context.Context, integrationID, hookID string, params model.RequestParams) (int, []interface{}, error) {
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

	// Find the right integration
	resourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegration, integrationID)
	integrationConfig, ok := s.projectConfig.Integrations[resourceID]
	if !ok || integrationConfig == nil {
		return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) does not exist", integrationID), nil, nil)
	}

	// Return the provided hook if id is present
	if hookID != "*" {
		resourceID = config.GenerateResourceID(s.clusterID, "noProject", config.ResourceIntegrationHook, hookID)
		hook, ok := s.projectConfig.IntegrationHooks[resourceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration hook (%s) does not exist", hookID), nil, nil)
		}

		hook.ID = hookID
		hook.IntegrationID = integrationID
		return http.StatusOK, []interface{}{hook}, nil
	}

	// Create an array of hooks
	hooks := make([]interface{}, 0)
	for k, v := range s.projectConfig.IntegrationHooks {
		if v.IntegrationID != integrationID {
			continue
		}
		v.ID = k
		v.IntegrationID = integrationID
		hooks = append(hooks, v)
	}

	// Return the hooks for that integration
	return http.StatusOK, hooks, nil
}

// GetIntegrationTokens returns the tokens required for an integration
func (s *Manager) GetIntegrationTokens(ctx context.Context, id, key string) (int, interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Check if an integration by that id exists
	if s.projectConfig.Integrations == nil {
		return http.StatusNotFound, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) not found", id), nil, nil)
	}

	i, p := s.projectConfig.Integrations.Get(id)
	if !p {
		return http.StatusNotFound, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Integration (%s) not found", id), nil, nil)
	}

	// Check if the credentials are right
	if i.Key != key {
		return http.StatusNotFound, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid credentials provided for integration (%s)", id), nil, nil)
	}

	// Create the admin token
	adminToken, err := s.adminMan.GetIntegrationToken(id)
	if err != nil {
		return http.StatusInternalServerError, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create admin token for integration (%s)", id), err, nil)
	}

	// Create a token for each project
	projects := map[string]string{}
	for _, p := range s.projectConfig.Projects {
		// Skip if the project is an integration
		if p.ProjectConfig.IsIntegration && p.ProjectConfig.ID != id {
			continue
		}

		// Get auth module of that project
		a, err := s.modules.GetAuthModuleForSyncMan(p.ProjectConfig.ID)
		if err != nil {
			return http.StatusInternalServerError, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get auth module of project (%s) for integration (%s)", p.ProjectConfig.ID, id), err, nil)
		}

		// Create the token for the project
		projects[p.ProjectConfig.ID], err = a.GetIntegrationToken(ctx, id)
		if err != nil {
			return http.StatusInternalServerError, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create token for project (%s) for integration (%s)", p.ProjectConfig.ID, id), err, nil)
		}
	}

	// Return the tokens we have created for the integration
	return http.StatusOK, map[string]interface{}{"admin": adminToken, "api": projects}, nil
}
