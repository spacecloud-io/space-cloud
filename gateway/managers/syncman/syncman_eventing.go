package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetEventingRule sets the eventing rules
func (s *Manager) SetEventingRule(ctx context.Context, project, ruleName string, value *config.EventingTrigger, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = ruleName
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingTrigger, ruleName)
	if projectConfig.EventingTriggers == nil {
		projectConfig.EventingTriggers = config.EventingTriggers{resourceID: value}
	} else {
		projectConfig.EventingTriggers[resourceID] = value
	}

	if err := s.modules.SetEventingTriggerConfig(ctx, projectConfig.EventingTriggers); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, value); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteEventingRule deletes an eventing rule
func (s *Manager) SetDeleteEventingRule(ctx context.Context, project, ruleName string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingTrigger, ruleName)
	delete(projectConfig.EventingTriggers, resourceID)

	if err := s.modules.SetEventingTriggerConfig(ctx, projectConfig.EventingTriggers); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetEventingConfig sets the eventing config
func (s *Manager) SetEventingConfig(ctx context.Context, project, dbAlias string, enabled bool, dbTableInclusionMap map[string][]string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	dbConfig, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias)
	if !p && enabled {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unknown db alias (%s) provided while setting eventing config", dbAlias), nil, nil)
	}

	projectConfig.EventingConfig.DBAlias = dbAlias
	projectConfig.EventingConfig.Enabled = enabled
	projectConfig.EventingConfig.DBTablesInclusionMap = dbTableInclusionMap

	if err := s.modules.SetEventingConfig(ctx, project, projectConfig.EventingConfig, projectConfig.EventingRules, projectConfig.EventingSchemas, projectConfig.EventingTriggers); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if enabled {
		if err := s.applySchemas(ctx, project, dbAlias, projectConfig, config.CrudStub{
			Collections: map[string]*config.TableRule{
				utils.TableEventingLogs:   {Schema: utils.SchemaEventLogs, Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}},
				utils.TableInvocationLogs: {Schema: utils.SchemaInvocationLogs, Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}},
			},
			DBName: dbConfig.DBName,
		}); err != nil {
			return http.StatusInternalServerError, err
		}
		status, err := s.setCollectionRules(ctx, projectConfig, project, dbAlias, utils.TableEventingLogs, &config.DatabaseRule{Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}})
		if err != nil {
			return status, err
		}
		status, err = s.setCollectionRules(ctx, projectConfig, project, dbAlias, utils.TableInvocationLogs, &config.DatabaseRule{Rules: map[string]*config.Rule{"create": {Rule: "deny"}, "read": {Rule: "deny"}, "update": {Rule: "deny"}, "delete": {Rule: "deny"}}})
		if err != nil {
			return status, err
		}
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingConfig, "eventing")
	if err := s.store.SetResource(ctx, resourceID, projectConfig.EventingConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetEventingConfig returns the eventing config
func (s *Manager) GetEventingConfig(ctx context.Context, project string, params model.RequestParams) (int, interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, projectConfig.EventingConfig, nil
}

// SetEventingSchema sets the schema for the given event type
func (s *Manager) SetEventingSchema(ctx context.Context, project string, evType string, schema string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingSchema, evType)
	v := &config.EventingSchema{Schema: schema, ID: evType}
	if projectConfig.EventingSchemas == nil {
		projectConfig.EventingSchemas = config.EventingSchemas{resourceID: v}
	} else {
		projectConfig.EventingSchemas[resourceID] = v
	}

	if err := s.modules.SetEventingSchemaConfig(ctx, projectConfig.EventingSchemas); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteEventingSchema deletes the schema for the given event type
func (s *Manager) SetDeleteEventingSchema(ctx context.Context, project string, evType string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}
	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingSchema, evType)
	delete(projectConfig.EventingSchemas, resourceID)

	if err := s.modules.SetEventingSchemaConfig(ctx, projectConfig.EventingSchemas); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
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
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingRule, evType)
	if projectConfig.EventingRules == nil {
		projectConfig.EventingRules = config.EventingRules{resourceID: rule}
	} else {
		projectConfig.EventingRules[resourceID] = rule
	}

	if err := s.modules.SetEventingRuleConfig(ctx, projectConfig.EventingRules); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, rule); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteEventingSecurityRules deletes the security rule for the given event type
func (s *Manager) SetDeleteEventingSecurityRules(ctx context.Context, project, evType string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingRule, evType)
	delete(projectConfig.EventingRules, resourceID)

	if err := s.modules.SetEventingRuleConfig(ctx, projectConfig.EventingRules); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting eventing config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetEventingTriggerRules gets trigger rules from config
func (s *Manager) GetEventingTriggerRules(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingTrigger, id)
		service, ok := projectConfig.EventingTriggers[resourceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Trigger rule (%s) does not exists for eventing config", id), nil, nil)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.EventingTriggers {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}

// GetEventingSchema gets eventing schema from config
func (s *Manager) GetEventingSchema(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingSchema, id)
		service, ok := projectConfig.EventingSchemas[resourceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Schema (%s) does not exists in eventing config", id), nil, nil)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.EventingSchemas {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}

// GetEventingSecurityRules gets eventing security rules from config
func (s *Manager) GetEventingSecurityRules(ctx context.Context, project, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if id != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceEventingRule, id)
		service, ok := projectConfig.EventingRules[resourceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Security rule (%s) does not exists for eventing config", id), nil, nil)
		}

		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.EventingRules {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}
