package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetUserManagement sets the user management
func (s *Manager) SetUserManagement(ctx context.Context, project, provider string, value *config.AuthStub, reqParams model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = provider
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	projectConfig.Modules.Auth[provider] = value

	if err := s.modules.SetUsermanConfig(project, projectConfig.Modules.Auth); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetUserManagement gets user management
func (s *Manager) GetUserManagement(ctx context.Context, project, providerID string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if providerID != "*" {
		auth, ok := projectConfig.Modules.Auth[providerID]
		if !ok {
			return http.StatusBadRequest, nil, fmt.Errorf("providerID (%s) not present in config", providerID)
		}

		return http.StatusOK, []interface{}{auth}, nil
	}

	providers := []interface{}{}
	for _, value := range projectConfig.Modules.Auth {
		providers = append(providers, value)
	}

	return http.StatusOK, providers, nil
}
