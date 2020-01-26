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

	return s.setProject(ctx, projectConfig)
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

	return s.setProject(ctx, projectConfig)
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

	return s.setProject(ctx, projectConfig)
}

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

func (s *Manager) SetEventingSecurityRules(ctx context.Context, project, evType string, rule *config.Rule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	if len(projectConfig.Modules.Eventing.SecurityRules) != 0 {
		projectConfig.Modules.Eventing.SecurityRules[evType] = &config.Rule{Rule: rule.Rule, Eval: rule.Eval, Type: rule.Type, F1: rule.F1, F2: rule.F2, Clauses: rule.Clauses, DB: rule.DB, Col: rule.Col, Find: rule.Find, URL: rule.URL, Fields: rule.Fields, Field: rule.Field, Value: rule.Value}
	} else {
		projectConfig.Modules.Eventing.SecurityRules = map[string]*config.Rule{
			evType: &config.Rule{Rule: rule.Rule, Eval: rule.Eval, Type: rule.Type, F1: rule.F1, F2: rule.F2, Clauses: rule.Clauses, DB: rule.DB, Col: rule.Col, Find: rule.Find, URL: rule.URL, Fields: rule.Fields, Field: rule.Field, Value: rule.Value},
		}
	}

	return s.setProject(ctx, projectConfig)
}

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
