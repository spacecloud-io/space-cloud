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

// SetSecurityFunction sets the security function
func (s *Manager) SetSecurityFunction(ctx context.Context, project, id string, securityFunction *config.SecurityFunction, reqParams model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	securityFunction.ID = id
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceSecurityFunction, id)
	if projectConfig.SecurityFunctions == nil {
		projectConfig.SecurityFunctions = config.SecurityFunctions{resourceID: securityFunction}
	} else {
		projectConfig.SecurityFunctions[resourceID] = securityFunction
	}

	if err := s.modules.SetSecurityFunctionConfig(ctx, project, projectConfig.SecurityFunctions); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetResource(ctx, resourceID, securityFunction); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetSecurityFunction gets security function
func (s *Manager) GetSecurityFunction(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		securityFunction, ok := projectConfig.SecurityFunctions[config.GenerateResourceID(s.clusterID, project, config.ResourceSecurityFunction, id)]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("security function with name (%s) does not exist in auth config", id), nil, nil)
		}
		return http.StatusOK, []interface{}{securityFunction}, nil
	}

	securityFunctions := make([]interface{}, 0)
	for _, value := range projectConfig.SecurityFunctions {
		securityFunctions = append(securityFunctions, value)
	}

	return http.StatusOK, securityFunctions, nil
}

// DeleteSecurityFunction deletes the security function
func (s *Manager) DeleteSecurityFunction(ctx context.Context, project, id string, reqParams model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceSecurityFunction, id)

	delete(projectConfig.SecurityFunctions, resourceID)

	if err := s.modules.SetSecurityFunctionConfig(ctx, project, projectConfig.SecurityFunctions); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
