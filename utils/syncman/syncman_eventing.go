package syncman

import "github.com/spaceuptech/space-cloud/config"

func (s *Manager) SetEventingRule(project, ruleName string, value config.EventingRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Eventing.Rules[ruleName] = value

	return s.setProject(projectConfig)
}

func (s *Manager) SetDeleteEventingRule(project, ruleName string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Eventing.Rules, ruleName)

	return s.setProject(projectConfig)
}

func (s *Manager) SetEventingConfig(project, dbType, col string, enabled bool) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Eventing.DBType = dbType
	projectConfig.Modules.Eventing.Col = col
	projectConfig.Modules.Eventing.Enabled = enabled

	return s.setProject(projectConfig)
}
