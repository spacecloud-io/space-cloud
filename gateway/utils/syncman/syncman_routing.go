package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetProjectRoutes sets a projects routes
func (s *Manager) SetProjectRoutes(ctx context.Context, project string, c config.Routes) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// Update the project's routes
	projectConfig.Modules.Routes = c

	// Apply the config
	s.projects.SetRoutingConfig(project, c)

	// Persist the config
	return s.persistProjectConfig(ctx, projectConfig)
}
