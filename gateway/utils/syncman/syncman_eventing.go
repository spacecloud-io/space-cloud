package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetEventingRule sets the eventing rules
func (s *Manager) SetEventingRule(ctx context.Context, project, ruleName string, value config.EventingRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	projectConfig.Modules.Eventing.Rules[ruleName] = value

	return s.setProject(ctx, projectConfig)
}

// SetDeleteEventingRule deletes an eventing rule
func (s *Manager) SetDeleteEventingRule(ctx context.Context, project, ruleName string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Eventing.Rules, ruleName)

	return s.setProject(ctx, projectConfig)
}

// SetEventingConfig sets the eventing config
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

	return s.setProject(ctx, projectConfig)
}

// SetEventingSchema sets the schema for the given event type
func (s *Manager) SetEventingSchema(ctx context.Context, project string, evType string, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	if len(projectConfig.Modules.Eventing.Schemas) != 0 {
		projectConfig.Modules.Eventing.Schemas[evType] = config.SchemaObject{Schema: schema}
	} else {
		projectConfig.Modules.Eventing.Schemas = map[string]config.SchemaObject{
			evType: config.SchemaObject{Schema: schema},
		}
	}

	return s.setProject(ctx, projectConfig)
}

// SetDeleteEventingSchema deletes the schema for the given event type
func (s *Manager) SetDeleteEventingSchema(ctx context.Context, project string, evType string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Eventing.Schemas, evType)

	return s.setProject(ctx, projectConfig)
}

// SetEventingSecurityRules sets the securtiy rule for the given event type
func (s *Manager) SetEventingSecurityRules(ctx context.Context, project, evType string, rule *config.Rule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	if len(projectConfig.Modules.Eventing.SecurityRules) != 0 {
		projectConfig.Modules.Eventing.SecurityRules[evType] = rule
	} else {
		projectConfig.Modules.Eventing.SecurityRules = map[string]*config.Rule{
			evType: rule,
		}
	}

	return s.setProject(ctx, projectConfig)
}

// SetDeleteEventingSecurityRules deletes the security rule for the given event type
func (s *Manager) SetDeleteEventingSecurityRules(ctx context.Context, project, evType string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	delete(projectConfig.Modules.Eventing.SecurityRules, evType)

	return s.setProject(ctx, projectConfig)
}
