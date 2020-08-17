package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetEventingRule sets the eventing rules
func (s *Manager) SetEventingRule(ctx context.Context, project, ruleName string, value config.EventingRule, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = ruleName
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if projectConfig.Modules.Eventing.Rules == nil {
		projectConfig.Modules.Eventing.Rules = map[string]config.EventingRule{}
	}
	projectConfig.Modules.Eventing.Rules[ruleName] = value

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteEventingRule deletes an eventing rule
func (s *Manager) SetDeleteEventingRule(ctx context.Context, project, ruleName string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}
	delete(projectConfig.Modules.Eventing.Rules, ruleName)

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetEventingConfig sets the eventing config
func (s *Manager) SetEventingConfig(ctx context.Context, project, dbAlias string, enabled bool, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}
	_, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok && enabled {
		return http.StatusBadRequest, fmt.Errorf("unknown db (%s) provided while setting eventing config", dbAlias)
	}

	projectConfig.Modules.Eventing.DBAlias = dbAlias
	projectConfig.Modules.Eventing.Enabled = enabled

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if enabled {
		if err := s.applySchemas(ctx, project, dbAlias, projectConfig, config.CrudStub{
			Collections: map[string]*config.TableRule{
				utils.TableEventingLogs:   {Schema: utils.SchemaEventLogs, Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}},
				utils.TableInvocationLogs: {Schema: utils.SchemaInvocationLogs, Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}},
			},
		}); err != nil {
			return http.StatusInternalServerError, err
		}
		status, err := s.setCollectionRules(ctx, projectConfig, project, dbAlias, utils.TableEventingLogs, &config.TableRule{Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}})
		if err != nil {
			return status, err
		}
		status, err = s.setCollectionRules(ctx, projectConfig, project, dbAlias, utils.TableInvocationLogs, &config.TableRule{Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}})
		if err != nil {
			return status, err
		}
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetEventingConfig returns the eventing config
func (s *Manager) GetEventingConfig(ctx context.Context, project string, params model.RequestParams) (int, interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	eventing := projectConfig.Modules.Eventing
	return http.StatusOK, config.Eventing{DBAlias: eventing.DBAlias, Enabled: eventing.Enabled}, nil
}

// SetEventingSchema sets the schema for the given event type
func (s *Manager) SetEventingSchema(ctx context.Context, project string, evType string, schema string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
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
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteEventingSchema deletes the schema for the given event type
func (s *Manager) SetDeleteEventingSchema(ctx context.Context, project string, evType string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}
	delete(projectConfig.Modules.Eventing.Schemas, evType)

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetEventingSecurityRules sets the securtiy rule for the given event type
func (s *Manager) SetEventingSecurityRules(ctx context.Context, project, evType string, rule *config.Rule, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	rule.ID = evType
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
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
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteEventingSecurityRules deletes the security rule for the given event type
func (s *Manager) SetDeleteEventingSecurityRules(ctx context.Context, project, evType string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	delete(projectConfig.Modules.Eventing.SecurityRules, evType)

	if err := s.modules.SetEventingConfig(project, &projectConfig.Modules.Eventing); err != nil {
		logrus.Errorf("error setting eventing config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetEventingTriggerRules gets trigger rules from config
func (s *Manager) GetEventingTriggerRules(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		service, ok := projectConfig.Modules.Eventing.Rules[id]
		if !ok {
			return http.StatusBadRequest, nil, fmt.Errorf("id (%s) not present in config", id)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Eventing.Rules {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}

// GetEventingSchema gets eventing schema from config
func (s *Manager) GetEventingSchema(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		service, ok := projectConfig.Modules.Eventing.Schemas[id]
		if !ok {
			return http.StatusBadRequest, nil, fmt.Errorf("id (%s) not present in config", id)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Eventing.Schemas {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}

// GetEventingSecurityRules gets eventing security rules from config
func (s *Manager) GetEventingSecurityRules(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		service, ok := projectConfig.Modules.Eventing.SecurityRules[id]
		if !ok {
			return http.StatusBadRequest, nil, fmt.Errorf("id (%s) not present in config", id)
		}

		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Eventing.SecurityRules {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}
