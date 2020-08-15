package syncman

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/integration"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
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
	if err := s.adminMan.ValidateIntegrationSyncOperation(config.Integrations{integrationConfig}); err != nil {
		return http.StatusUpgradeRequired, err
	}

	// Create a project if it doesn't already exist
	utils.LogDebug(fmt.Sprintf("Creating a new project for integration (%s)", integrationConfig.ID), "syncman", "enable-integration", nil)

	// Seed the random byte generator
	rand.Seed(time.Now().UnixNano())

	// We need runner enabled for this one. We need to create a new project, secret and deployment in runner
	if s.runnerAddr == "" {
		return http.StatusBadRequest, utils.LogError("Runner must be enabled for integrations to work", "syncman", "enable-integration", nil)
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
		return http.StatusBadRequest, utils.LogError("Unable to parse integration license", "syncman", "enable-integration", err)
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
		return http.StatusInternalServerError, utils.LogError("Unable to create new aes key for integration", "syncman", "enable-integration", nil)
	}

	// Check if integration already exists
	var integrations config.Integrations
	var proj *config.Project
	if _, p := s.projectConfig.Admin.Integrations.Get(integrationConfig.ID); !p {
		// Prepare a unique
		proj = &config.Project{
			ID:                 integrationConfig.ID,
			Name:               integrationConfig.Name,
			SecretSource:       integrationConfig.SecretSource,
			ContextTimeGraphQL: 20,
			AESKey:             base64.StdEncoding.EncodeToString(key),
			IsIntegration:      true,
		}

		integrationConfig.Hooks = map[string]*config.IntegrationHook{}

		integrations = append(s.projectConfig.Admin.Integrations, integrationConfig)
	} else {
		for _, i := range s.projectConfig.Admin.Integrations {
			if i.ID != integrationConfig.ID {
				integrations = append(integrations, i)
			}
		}
		integrations = append(integrations, integrationConfig)
	}

	if err := s.integrationMan.SetConfig(integrations); err != nil {
		return http.StatusUpgradeRequired, err
	}

	s.projectConfig.Admin.Integrations = integrations

	if err := s.store.SetAdminConfig(ctx, s.projectConfig.Admin); err != nil {
		return http.StatusInternalServerError, err
	}

	if proj != nil {
		s.projectConfig.Projects = append(s.projectConfig.Projects, proj)
		if err := s.modules.SetProjectConfig(proj); err != nil {
			return http.StatusInternalServerError, utils.LogError("Unable to create new project for integration", "syncman", "enable-integration", nil)
		}

		// Update the store
		if err := s.store.SetProject(ctx, proj); err != nil {
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
	for i, p := range s.projectConfig.Projects {
		if p.ID == id && p.IsIntegration {
			length := len(s.projectConfig.Projects)
			s.projectConfig.Projects[i] = s.projectConfig.Projects[length-1]
			s.projectConfig.Projects = s.projectConfig.Projects[:length-1]
			break
		}
	}

	// Remove integration from integrations array
	for i, c := range s.projectConfig.Admin.Integrations {
		if c.ID == id {
			length := len(s.projectConfig.Admin.Integrations)
			s.projectConfig.Admin.Integrations[i] = s.projectConfig.Admin.Integrations[length-1]
			s.projectConfig.Admin.Integrations = s.projectConfig.Admin.Integrations[:length-1]
			break
		}
	}

	// Update the modules and integration manager
	s.modules.Delete(id)
	_ = s.integrationMan.SetConfig(s.projectConfig.Admin.Integrations)

	// Update the stores
	if err := s.store.DeleteProject(ctx, id); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetAdminConfig(ctx, s.projectConfig.Admin); err != nil {
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
	for _, i := range s.projectConfig.Admin.Integrations {
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

	return http.StatusBadRequest, nil, utils.LogError(fmt.Sprintf("Integration (%s) not found", id), "syncman", "get-integrations", nil)
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
	var integrationConfig *config.IntegrationConfig
	for _, i := range s.projectConfig.Admin.Integrations {
		if i.ID == integrationID {
			integrationConfig = i
			break
		}
	}

	if integrationConfig == nil {
		return http.StatusBadRequest, utils.LogError(fmt.Sprintf("Integration (%s) does not exist", integrationID), "syncman", "add-hook", nil)
	}

	// Check if the integration has the permissions to create this hook
	if !integration.HasPermissionForHook(integrationConfig, hookConfig) {
		return http.StatusBadRequest, utils.LogError(fmt.Sprintf("Integration (%s) does not have necessary permissions to create the hook", integrationID), "syncman", "add-hook", nil)
	}

	// Create an empty map if nil
	if integrationConfig.Hooks == nil {
		integrationConfig.Hooks = map[string]*config.IntegrationHook{}
	}

	// Add the hook and store the config
	integrationConfig.Hooks[hookConfig.ID] = hookConfig
	_ = s.integrationMan.SetConfig(s.projectConfig.Admin.Integrations)

	// Store the config in the store
	if err := s.store.SetAdminConfig(ctx, s.projectConfig.Admin); err != nil {
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
	var integrationConfig *config.IntegrationConfig
	for _, i := range s.projectConfig.Admin.Integrations {
		if i.ID == integrationID {
			integrationConfig = i
			break
		}
	}

	// Throw error if hook does not exist
	if integrationConfig == nil {
		return http.StatusBadRequest, utils.LogError(fmt.Sprintf("Integration (%s) does not exist", integrationID), "syncman", "remove-hook", nil)
	}

	// Delete the hook
	delete(integrationConfig.Hooks, hookID)
	_ = s.integrationMan.SetConfig(s.projectConfig.Admin.Integrations)

	// Store the config in the store
	if err := s.store.SetAdminConfig(ctx, s.projectConfig.Admin); err != nil {
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
	var integrationConfig *config.IntegrationConfig
	for _, i := range s.projectConfig.Admin.Integrations {
		if i.ID == integrationID {
			integrationConfig = i
			break
		}
	}

	// Throw error if hook does not exist
	if integrationConfig == nil {
		return http.StatusBadRequest, nil, utils.LogError(fmt.Sprintf("Integration (%s) does not exist", integrationID), "syncman", "get-hook", nil)
	}

	// Return the provided hook if id is present
	if hookID != "*" {
		hook, p := integrationConfig.Hooks[hookID]
		if !p {
			return http.StatusBadRequest, nil, utils.LogError(fmt.Sprintf("Integration hook (%s) does not exist", hookID), "syncman", "get-hook", nil)
		}

		hook.ID = hookID
		hook.IntegrationID = integrationID
		return http.StatusOK, []interface{}{hook}, nil
	}

	// Create an array of hooks
	hooks := make([]interface{}, 0)
	for k, v := range integrationConfig.Hooks {
		v.ID = k
		v.IntegrationID = integrationID
		hooks = append(hooks, v)
	}

	// Return the hooks for that integration
	return http.StatusOK, hooks, nil
}

// GetIntegrationTokens returns the tokens required for an integration
func (s *Manager) GetIntegrationTokens(id, key string) (int, interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Check if an integration by that id exists
	if s.projectConfig.Admin.Integrations == nil {
		return http.StatusNotFound, nil, utils.LogError(fmt.Sprintf("Integration (%s) not found", id), "syncman", "integration-tokens", nil)
	}

	i, p := s.projectConfig.Admin.Integrations.Get(id)
	if !p {
		return http.StatusNotFound, nil, utils.LogError(fmt.Sprintf("Integration (%s) not found", id), "syncman", "integration-tokens", nil)
	}

	// Check if the credentials are right
	if i.Key != key {
		return http.StatusNotFound, nil, utils.LogError(fmt.Sprintf("Invalid credentials provided for integration (%s)", id), "syncman", "integration-tokens", nil)
	}

	// Create the admin token
	adminToken, err := s.adminMan.GetIntegrationToken(id)
	if err != nil {
		return http.StatusInternalServerError, nil, utils.LogError(fmt.Sprintf("Unable to create admin token for integration (%s)", id), "syncman", "integration-tokens", err)
	}

	// Create a token for each project
	projects := map[string]string{}
	for _, p := range s.projectConfig.Projects {
		// Skip if the project is an integration
		if p.IsIntegration && p.ID != id {
			continue
		}

		// Get auth module of that project
		a, err := s.modules.GetAuthModuleForSyncMan(p.ID)
		if err != nil {
			return http.StatusInternalServerError, nil, utils.LogError(fmt.Sprintf("Unable to get auth module of project (%s) for integration (%s)", p.ID, id), "syncman", "integration-tokens", err)
		}

		// Create the token for the project
		projects[p.ID], err = a.GetIntegrationToken(id)
		if err != nil {
			return http.StatusInternalServerError, nil, utils.LogError(fmt.Sprintf("Unable to create token for project (%s) for integration (%s)", p.ID, id), "syncman", "integration-tokens", err)
		}
	}

	// Return the tokens we have created for the integration
	return http.StatusOK, map[string]interface{}{"admin": adminToken, "api": projects}, nil
}
