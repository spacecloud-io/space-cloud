package syncman

import "github.com/spaceuptech/space-cloud/config"

func (s *Manager) SetEventingRule(project *config.Project, ruleName string, value config.EventingRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	project.Modules.Eventing.Rules[ruleName] = value

	return s.setProject(project)
}

func (s *Manager) SetDeleteEventingRule(project *config.Project, ruleName string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(project.Modules.Eventing.Rules, ruleName)

	return s.setProject(project)
}

func (s *Manager) SetEventingStatus(project *config.Project, dbType, col string, enabled bool) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	project.Modules.Eventing.DBType = dbType
	project.Modules.Eventing.Col = col
	project.Modules.Eventing.Enabled = enabled

	return s.setProject(project)
}
