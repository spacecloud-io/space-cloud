package syncman

import "github.com/spaceuptech/space-cloud/config"

func (s *Manager) SetUserManagement(project, provider string, value *config.AuthStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Auth[provider] = value

	return s.setProject(projectConfig)
}
