package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetUserManagement sets the user management
func (s *Manager) SetUserManagement(ctx context.Context, project, provider string, value *config.AuthStub, reqParams model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, reqParams)
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

	value.ID = provider
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceAuthProvider, provider)
	if projectConfig.Auths == nil {
		projectConfig.Auths = config.Auths{resourceID: value}
	} else {
		projectConfig.Auths[resourceID] = value
	}

	if err := s.modules.SetUsermanConfig(ctx, project, projectConfig.Auths); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetResource(ctx, resourceID, value); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetUserManagement gets user management
func (s *Manager) GetUserManagement(ctx context.Context, project, providerID string, params model.RequestParams) (int, []interface{}, error) {
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

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if providerID != "*" {
		auth, ok := projectConfig.Auths[config.GenerateResourceID(s.clusterID, project, config.ResourceAuthProvider, providerID)]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("provider with id (%s) does not exist in user management config", providerID), nil, nil)
		}

		return http.StatusOK, []interface{}{auth}, nil
	}

	providers := []interface{}{}
	for _, value := range projectConfig.Auths {
		providers = append(providers, value)
	}

	return http.StatusOK, providers, nil
}

// DeleteUserManagement deletes the user management
func (s *Manager) DeleteUserManagement(ctx context.Context, project, provider string, reqParams model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceAuthProvider, provider)

	delete(projectConfig.Auths, resourceID)

	if err := s.modules.SetUsermanConfig(ctx, project, projectConfig.Auths); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
