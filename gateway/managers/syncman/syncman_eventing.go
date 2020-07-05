package syncman

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetEventingRule sets the eventing rules
func (s *Manager) SetEventingRule(ctx context.Context, project, ruleName string, value config.EventingRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = ruleName
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	if projectConfig.Modules.Eventing.Rules == nil {
		projectConfig.Modules.Eventing.Rules = map[string]config.EventingRule{}
	}
	projectConfig.Modules.Eventing.Rules[ruleName] = value

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
	}

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

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetEventingConfig sets the eventing config
func (s *Manager) SetEventingConfig(ctx context.Context, project, dbAlias string, enabled bool) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	_, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok && enabled {
		return fmt.Errorf("unknown db (%s) provided while setting eventing config", dbAlias)
	}
	projectConfig.Modules.Eventing.DBAlias = dbAlias
	projectConfig.Modules.Eventing.Enabled = enabled

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
	}

	if enabled {
		if err := s.applySchemas(ctx, project, dbAlias, projectConfig, config.CrudStub{
			Collections: map[string]*config.TableRule{
				utils.TableEventingLogs:   {Schema: utils.SchemaEventLogs},
				utils.TableInvocationLogs: {Schema: utils.SchemaInvocationLogs},
			},
		}); err != nil {
			return err
		}
	}

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
		projectConfig.Modules.Eventing.Schemas[evType] = config.SchemaObject{Schema: schema, ID: evType}
	} else {
		projectConfig.Modules.Eventing.Schemas = map[string]config.SchemaObject{
			evType: {Schema: schema, ID: evType},
		}
	}

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
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

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetEventingSecurityRules sets the securtiy rule for the given event type
func (s *Manager) SetEventingSecurityRules(ctx context.Context, project, evType string, rule *config.Rule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	rule.ID = evType
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

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
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

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// GetEventingTriggerRules gets trigger rules from config
func (s *Manager) GetEventingTriggerRules(ctx context.Context, project, id string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if id != "*" {
		service, ok := projectConfig.Modules.Eventing.Rules[id]
		if !ok {
			return nil, fmt.Errorf("id (%s) not present in config", id)
		}
		return []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Eventing.Rules {
		services = append(services, value)
	}
	return services, nil
}

// GetEventingSchema gets eventing schema from config
func (s *Manager) GetEventingSchema(ctx context.Context, project, id string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if id != "*" {
		service, ok := projectConfig.Modules.Eventing.Schemas[id]
		if !ok {
			return nil, fmt.Errorf("id (%s) not present in config", id)
		}
		return []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Eventing.Schemas {
		services = append(services, value)
	}
	return services, nil
}

// GetEventingSecurityRules gets eventing security rules from config
func (s *Manager) GetEventingSecurityRules(ctx context.Context, project, id string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if id != "*" {
		service, ok := projectConfig.Modules.Eventing.SecurityRules[id]
		if !ok {
			return nil, fmt.Errorf("id (%s) not present in config", id)
		}
		return []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Eventing.SecurityRules {
		services = append(services, value)
	}
	return services, nil
}
