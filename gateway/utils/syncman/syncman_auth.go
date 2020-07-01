package syncman

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetUserManagement sets the user management
func (s *Manager) SetUserManagement(ctx context.Context, project, provider string, value *config.AuthStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = provider
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Auth[provider] = value

	s.modules.SetUsermanConfig(project, projectConfig.Modules.Auth)

	return s.setProject(ctx, projectConfig)
}

// GetUserManagement gets user management
func (s *Manager) GetUserManagement(ctx context.Context, project, providerID string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if providerID != "*" {
		auth, ok := projectConfig.Modules.Auth[providerID]
		if !ok {
			return nil, fmt.Errorf("providerID (%s) not present in config", providerID)
		}
		return []interface{}{auth}, nil
	}

	providers := []interface{}{}
	for _, value := range projectConfig.Modules.Auth {
		providers = append(providers, value)
	}
	return providers, nil
}
