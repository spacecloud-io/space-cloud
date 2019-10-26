package syncman

import "github.com/spaceuptech/space-cloud/config"

func (s *Manager) SetUserManagement(project *config.Project, provider string, value *config.AuthStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	project.Modules.Auth[provider] = value

	return s.setProject(project)
}
