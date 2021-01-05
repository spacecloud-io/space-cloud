package syncman

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetDeleteCollection deletes a collection from the database
func (s *Manager) SetDeleteCollection(ctx context.Context, project, dbAlias, col string, module *crud.Module, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get prepared query provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseSchema, dbAlias, col)
	delete(projectConfig.DatabaseSchemas, resourceID)

	if err := s.modules.SetDatabaseSchemaConfig(ctx, project, projectConfig.DatabaseSchemas); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := module.DeleteTable(ctx, dbAlias, col); err != nil {
		return http.StatusInternalServerError, err
	}

	dbRulesResourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseRule, dbAlias, col, "rule")
	delete(projectConfig.DatabaseRules, dbRulesResourceID)

	if err := s.modules.SetDatabaseRulesConfig(ctx, project, projectConfig.DatabaseRules); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.DeleteResource(ctx, dbRulesResourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDatabaseConnection sets the database connection
func (s *Manager) SetDatabaseConnection(ctx context.Context, project, dbAlias string, v *config.DatabaseConfig, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	v.DbAlias = dbAlias
	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseConfig, dbAlias)
	if projectConfig.DatabaseConfigs == nil {
		projectConfig.DatabaseConfigs = config.DatabaseConfigs{resourceID: v}
	} else {
		projectConfig.DatabaseConfigs[resourceID] = v
	}

	if err := s.modules.SetDatabaseConfig(ctx, project, projectConfig.DatabaseConfigs, projectConfig.DatabaseSchemas, projectConfig.DatabaseRules, projectConfig.DatabasePreparedQueries); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemoveDatabaseConfig removes the database config
func (s *Manager) RemoveDatabaseConfig(ctx context.Context, project, dbAlias string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceIDs := make([]string, 0)

	// update database config
	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseConfig, dbAlias)
	resourceIDs = append(resourceIDs, resourceID)
	delete(projectConfig.DatabaseConfigs, resourceID)

	// delete rules
	for _, databaseRule := range projectConfig.DatabaseRules {
		if databaseRule.DbAlias == dbAlias {
			resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseRule, dbAlias, databaseRule.Table, "rule")
			resourceIDs = append(resourceIDs, resourceID)
			delete(projectConfig.DatabaseRules, resourceID)
		}
	}

	// delete schemas
	for _, databaseSchema := range projectConfig.DatabaseSchemas {
		if databaseSchema.DbAlias == dbAlias {
			resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseSchema, dbAlias, databaseSchema.Table)
			resourceIDs = append(resourceIDs, resourceID)
			delete(projectConfig.DatabaseSchemas, resourceID)
		}
	}

	// delete prepared queries
	for _, preparedQuery := range projectConfig.DatabasePreparedQueries {
		if preparedQuery.DbAlias == dbAlias {
			resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabasePreparedQuery, dbAlias, preparedQuery.ID)
			resourceIDs = append(resourceIDs, resourceID)
			delete(projectConfig.DatabasePreparedQueries, resourceID)
		}
	}

	if err := s.modules.SetDatabaseConfig(ctx, project, projectConfig.DatabaseConfigs, projectConfig.DatabaseSchemas, projectConfig.DatabaseRules, projectConfig.DatabasePreparedQueries); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	// Delete resources from store
	for _, resourceID := range resourceIDs {
		if err := s.store.DeleteResource(ctx, resourceID); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	_ = s.modules.Caching().PurgeCache(ctx, project, &model.CachePurgeRequest{Resource: config.ResourceDatabaseSchema, DbAlias: dbAlias, ID: "*"})

	return http.StatusOK, nil
}

// GetLogicalDatabaseName gets logical database name for provided db alias
func (s *Manager) GetLogicalDatabaseName(ctx context.Context, project, dbAlias string) (string, error) {
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return "", err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseConfig, dbAlias)
	dbConfig, ok := projectConfig.DatabaseConfigs[resourceID]
	if !ok {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get logical database name provided db alias (%s) does not exists", dbAlias), nil, nil)
	}
	return dbConfig.DBName, nil
}

// GetPreparedQuery gets preparedQuery from config
func (s *Manager) GetPreparedQuery(ctx context.Context, project, dbAlias, id string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if dbAlias != "*" {
		if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get prepared query as provided db alias (%s) does not exists", dbAlias), nil, nil)
		}

		if id != "*" {
			resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabasePreparedQuery, dbAlias, id)
			preparedQuery, ok := projectConfig.DatabasePreparedQueries[resourceID]
			if !ok {
				return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Prepared query with id (%s) not present in config", id), nil, nil)
			}
			return http.StatusOK, []interface{}{preparedQuery}, nil
		}
		coll := make([]interface{}, 0)
		for _, value := range projectConfig.DatabasePreparedQueries {
			if value.DbAlias == dbAlias {
				coll = append(coll, value)
			}
		}
		return http.StatusOK, coll, nil
	}
	coll := make([]interface{}, 0)
	for _, dbPreparedQuery := range projectConfig.DatabasePreparedQueries {
		coll = append(coll, dbPreparedQuery)
	}
	return http.StatusOK, coll, nil
}

// SetPreparedQueries sets database preparedqueries
func (s *Manager) SetPreparedQueries(ctx context.Context, project, dbAlias, id string, v *config.DatbasePreparedQuery, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	v.ID = id
	v.DbAlias = dbAlias
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to set prepared query as provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabasePreparedQuery, dbAlias, id)
	if projectConfig.DatabasePreparedQueries == nil {
		projectConfig.DatabasePreparedQueries = config.DatabasePreparedQueries{resourceID: v}
	} else {
		projectConfig.DatabasePreparedQueries[resourceID] = v
	}

	if err := s.modules.SetDatabasePreparedQueryConfig(ctx, project, projectConfig.DatabasePreparedQueries); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set database prepared query config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemovePreparedQueries removes the database PreparedQueries
func (s *Manager) RemovePreparedQueries(ctx context.Context, project, dbAlias, id string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to remove prepared query as provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	// update database reparedQueries
	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabasePreparedQuery, dbAlias, id)
	delete(projectConfig.DatabasePreparedQueries, resourceID)

	if err := s.modules.SetDatabasePreparedQueryConfig(ctx, project, projectConfig.DatabasePreparedQueries); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set database prepared query config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetModifySchema modifies the schema of table
func (s *Manager) SetModifySchema(ctx context.Context, project, dbAlias, col string, v *config.DatabaseSchema, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update schema in config
	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to modify schema provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseSchema, dbAlias, col)
	v.DbAlias = dbAlias
	v.Table = col

	// Modify the schema
	schemaMod, _ := s.modules.GetSchemaModuleForSyncMan(project)
	if err := schemaMod.SchemaModifyAll(ctx, dbAlias, projectConfig.DatabaseConfigs[config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseConfig, dbAlias)].DBName, config.DatabaseSchemas{resourceID: v}); err != nil {
		return http.StatusInternalServerError, err
	}

	if projectConfig.DatabaseSchemas == nil {
		projectConfig.DatabaseSchemas = config.DatabaseSchemas{resourceID: v}
	} else {
		projectConfig.DatabaseSchemas[resourceID] = v
	}

	if err := s.modules.SetDatabaseSchemaConfig(ctx, project, projectConfig.DatabaseSchemas); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetCollectionRules sets the collection rules of the database
func (s *Manager) SetCollectionRules(ctx context.Context, project, dbAlias, col string, v *config.DatabaseRule, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return s.setCollectionRules(ctx, projectConfig, project, dbAlias, col, v)
}

func (s *Manager) setCollectionRules(ctx context.Context, projectConfig *config.Project, project, dbAlias, col string, v *config.DatabaseRule) (int, error) {
	// update collection rules & is realtime in config
	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to set collection/table rules as provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseRule, dbAlias, col, "rule")
	v.Table = col
	v.DbAlias = dbAlias
	if projectConfig.DatabaseRules == nil {
		projectConfig.DatabaseRules = config.DatabaseRules{resourceID: v}
	} else {
		projectConfig.DatabaseRules[resourceID] = v
	}

	if err := s.modules.SetDatabaseRulesConfig(ctx, project, projectConfig.DatabaseRules); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set database rule config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteCollectionRules deletes the collection rules of the database
func (s *Manager) DeleteCollectionRules(ctx context.Context, project, dbAlias, col string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete collection/table rules as provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseRule, dbAlias, col, "rule")
	_, ok := projectConfig.DatabaseRules[resourceID]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete collection rules as provided table or collection (%s) does not exists", col), nil, nil)
	}

	delete(projectConfig.DatabaseRules, resourceID)

	if err := s.modules.SetDatabaseRulesConfig(ctx, project, projectConfig.DatabaseRules); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set database rules config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetReloadSchema reloads of the schema
func (s *Manager) SetReloadSchema(ctx context.Context, dbAlias, project string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}
	// Get the schema module
	schemaMod, _ := s.modules.GetSchemaModuleForSyncMan(project)

	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, errors.New("specified database not present in config")
	}

	for _, dbSchema := range projectConfig.DatabaseSchemas {
		if dbSchema.Table == "default" || dbSchema.DbAlias != dbAlias {
			continue
		}
		result, err := schemaMod.SchemaInspection(ctx, dbAlias, projectConfig.DatabaseConfigs[config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseConfig, dbAlias)].DBName, dbSchema.Table)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// set new schema in config & return in response body
		dbSchema.Schema = result
	}

	if err := s.modules.SetDatabaseSchemaConfig(ctx, project, projectConfig.DatabaseSchemas); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	for resourceID, v := range projectConfig.DatabaseSchemas {
		if err := s.store.SetResource(ctx, resourceID, v); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

// SetSchemaInspection inspects the schema
func (s *Manager) SetSchemaInspection(ctx context.Context, project, dbAlias, col, schema string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update schema in config
	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, errors.New("specified database not present in config")
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseSchema, dbAlias, col)
	v := &config.DatabaseSchema{Table: col, DbAlias: dbAlias, Schema: schema}
	if projectConfig.DatabaseSchemas == nil {
		projectConfig.DatabaseSchemas = config.DatabaseSchemas{resourceID: v}
	} else {
		projectConfig.DatabaseSchemas[resourceID] = v
	}
	if err := s.modules.SetDatabaseSchemaConfig(ctx, project, projectConfig.DatabaseSchemas); err != nil {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting crud config", err, nil)
	}

	if err := s.store.SetResource(ctx, resourceID, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemoveCollection removed the collection from the database collection schema in config
func (s *Manager) RemoveCollection(ctx context.Context, project, dbAlias, col string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update schema in config
	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return http.StatusBadRequest, errors.New("specified database not present in config")
	}

	dbSchemaResourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseSchema, dbAlias, col)
	delete(projectConfig.DatabaseSchemas, dbSchemaResourceID)

	if err := s.modules.SetDatabaseSchemaConfig(ctx, project, projectConfig.DatabaseSchemas); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting crud config", err, nil)
	}

	dbRulesResourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseRule, dbAlias, col, "rule")
	delete(projectConfig.DatabaseRules, dbRulesResourceID)

	if err := s.modules.SetDatabaseRulesConfig(ctx, project, projectConfig.DatabaseRules); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.store.DeleteResource(ctx, dbSchemaResourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.DeleteResource(ctx, dbRulesResourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetModifyAllSchema modifies schema of all tables
func (s *Manager) SetModifyAllSchema(ctx context.Context, dbAlias, project string, v config.CrudStub, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := s.applySchemas(ctx, project, dbAlias, projectConfig, v); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *Manager) applySchemas(ctx context.Context, project, dbAlias string, projectConfig *config.Project, v config.CrudStub) error {

	// update schema in config
	if _, p := s.checkIfDbAliasExists(projectConfig.DatabaseConfigs, dbAlias); !p {
		return errors.New("specified database not present in config")
	}

	if projectConfig.DatabaseSchemas == nil {
		projectConfig.DatabaseSchemas = make(config.DatabaseSchemas)
	}

	dbSchemas := make(config.DatabaseSchemas)
	for colName, colValue := range v.Collections {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseSchema, dbAlias, colName)
		dbSchemas[resourceID] = &config.DatabaseSchema{Table: colName, DbAlias: dbAlias, Schema: colValue.Schema}
		projectConfig.DatabaseSchemas[resourceID] = &config.DatabaseSchema{Table: colName, DbAlias: dbAlias, Schema: colValue.Schema}
	}

	schemaEventing, err := s.modules.GetSchemaModuleForSyncMan(project)
	if err != nil {
		return err
	}

	if err := schemaEventing.SchemaModifyAll(ctx, dbAlias, v.DBName, dbSchemas); err != nil {
		return err
	}

	if err := s.modules.SetDatabaseSchemaConfig(ctx, project, projectConfig.DatabaseSchemas); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	for resourceID, v := range dbSchemas {
		if err := s.store.SetResource(ctx, resourceID, v); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting crud config", err, nil)
		}
	}

	return nil
}

// GetDatabaseConfig gets database config
func (s *Manager) GetDatabaseConfig(ctx context.Context, project, dbAlias string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if dbAlias != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseConfig, dbAlias)
		dbConfig, ok := projectConfig.DatabaseConfigs[resourceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("specified dbAlias (%s) not present in config", dbAlias), nil, nil)
		}
		return http.StatusOK, []interface{}{dbConfig}, nil
	}

	services := make([]interface{}, 0)
	for _, value := range projectConfig.DatabaseConfigs {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}

// GetCollectionRules gets collection rules
func (s *Manager) GetCollectionRules(ctx context.Context, project, dbAlias, col string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if dbAlias != "*" && col != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceDatabaseRule, dbAlias, col, "rule")
		collectionInfo, ok := projectConfig.DatabaseRules[resourceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Specified collection/table (%s) not present in config of dbAlias (%s)", col, dbAlias), nil, nil)
		}
		return http.StatusOK, []interface{}{collectionInfo}, nil
	} else if dbAlias != "*" {
		rules := make([]interface{}, 0)
		for _, value := range projectConfig.DatabaseRules {
			if value.DbAlias == dbAlias {
				rules = append(rules, value)
			}
		}
		return http.StatusOK, rules, nil
	}
	result := make([]interface{}, 0)
	for _, dbRule := range projectConfig.DatabaseRules {
		result = append(result, dbRule)
	}
	return http.StatusOK, result, nil
}

// GetSchemas gets schemas from config
func (s *Manager) GetSchemas(ctx context.Context, project, dbAlias, col, format string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	a, _ := s.modules.GetSchemaModuleForSyncMan(project)
	arr, err := a.GetSchemaForDB(ctx, dbAlias, col, format)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, arr, nil
}

type result struct {
	Result []*secret `json:"result,omitempty"`
}

type secret struct {
	Data map[string]string `json:"data,omitempty"`
}

// GetSecrets gets secrets from runner
// This function should be called only from setConfig method of any module
func (s *Manager) GetSecrets(project, secretName, key string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate internal access token
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to get internal access token", err, map[string]interface{}{})
	}

	// makes http request to get secrets from runner
	var vPtr result
	url := fmt.Sprintf("http://%s/v1/runner/%s/secrets?id=%s", s.runnerAddr, project, secretName)
	if err := s.MakeHTTPRequest(ctx, "GET", url, token, "", map[string]interface{}{}, &vPtr); err != nil {
		return "", err
	}

	return vPtr.Result[0].Data[key], nil
}
