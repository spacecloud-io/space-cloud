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

	// Set the user man config
	if err := s.projects.SetUserManConfig(project, projectConfig.Modules.Auth); err != nil {
		return err
	}

	// Persist the config
	return s.persistProjectConfig(projectConfig)
}
