package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetService adds a remote service
func (s *Manager) SetService(ctx context.Context, project, service string, value *config.Service) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Services.Services[service] = value

	s.modules.SetServicesConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

	return s.setProject(ctx, projectConfig)
}

// DeleteService deletes a remotes service
func (s *Manager) DeleteService(ctx context.Context, project, service string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Services.Services, service)

	s.modules.SetServicesConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

	return s.setProject(ctx, projectConfig)
}
