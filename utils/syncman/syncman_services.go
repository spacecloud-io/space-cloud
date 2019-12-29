package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/config"
)

func (s *Manager) SetService(ctx context.Context, project, service string, value *config.Service) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Services.Services[service] = value

	// Set the services config
	if err := s.projects.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		return err
	}

	// Persist the config
	return s.persistProjectConfig(ctx, projectConfig)
}

func (s *Manager) SetDeleteService(ctx context.Context, project, service string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Services.Services, service)

	// Set the services config
	if err := s.projects.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		return err
	}

	// Persist the config
	return s.persistProjectConfig(ctx, projectConfig)
}
