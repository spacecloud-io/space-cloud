package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func (s *Manager) SetEventingRule(ctx context.Context, project, ruleName string, value config.EventingRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Eventing.Rules[ruleName] = value

	// Set the eventing config
	if err := s.projects.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		return err
	}

	// Persist the config
	return s.persistProjectConfig(ctx, projectConfig)
}

func (s *Manager) SetDeleteEventingRule(ctx context.Context, project, ruleName string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Eventing.Rules, ruleName)

	// Set the eventing config
	if err := s.projects.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		return err
	}

	// Persist the config
	return s.persistProjectConfig(ctx, projectConfig)
}

func (s *Manager) SetEventingConfig(ctx context.Context, project, dbType, col string, enabled bool) error {
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

	// Set the eventing config
	if err := s.projects.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		return err
	}

	// Persist the config
	return s.persistProjectConfig(ctx, projectConfig)
}
