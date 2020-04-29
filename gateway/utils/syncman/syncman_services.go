package syncman

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetService adds a remote service
func (s *Manager) SetService(ctx context.Context, project, service string, value *config.Service) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = service
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	if projectConfig.Modules.Services.Services == nil {
		projectConfig.Modules.Services.Services = config.Services{}
	}
	projectConfig.Modules.Services.Services[service] = value

	if err := s.modules.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		logrus.Errorf("error setting services config - %s", err.Error())
		return err
	}

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

	if err := s.modules.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		logrus.Errorf("error setting services config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// GetServices gets a remotes service
func (s *Manager) GetServices(ctx context.Context, project, serviceID string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if serviceID != "" {
		service, ok := projectConfig.Modules.Services.Services[serviceID]
		if !ok {
			return nil, fmt.Errorf("serviceID (%s) not present in config", serviceID)
		}
		return []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Services.Services {
		services = append(services, value)
	}
	return services, nil
}
