package syncman

import "github.com/spaceuptech/space-cloud/config"

func (s *Manager) SetService(project *config.Project, service string, value *config.Service) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	project.Modules.Services.Services[service] = value

	return s.setProject(project)
}

func (s *Manager) SetDeleteService(project *config.Project, service string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(project.Modules.Services.Services, service)

	return s.setProject(project)
}
