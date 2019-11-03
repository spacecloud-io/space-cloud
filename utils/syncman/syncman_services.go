package syncman

import "github.com/spaceuptech/space-cloud/config"

func (s *Manager) SetService(project, service string, value *config.Service) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Services.Services[service] = value

	return s.setProject(projectConfig)
}

func (s *Manager) SetDeleteService(project, service string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Services.Services, service)

	return s.setProject(projectConfig)
}
