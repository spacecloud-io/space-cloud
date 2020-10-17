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
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// remove collection from config
	coll, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to set delete collection provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	delete(coll.Collections, col)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := module.DeleteTable(ctx, dbAlias, col); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDatabaseConnection sets the database connection
func (s *Manager) SetDatabaseConnection(ctx context.Context, project, dbAlias string, v config.CrudStub, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	coll, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		projectConfig.Modules.Crud[dbAlias] = &config.CrudStub{Conn: v.Conn, Enabled: v.Enabled, Collections: map[string]*config.TableRule{}, Type: v.Type, DBName: v.DBName, Limit: v.Limit}
	} else {
		coll.Conn = v.Conn
		coll.Enabled = v.Enabled
		coll.Type = v.Type
		coll.DBName = v.DBName
		coll.Limit = v.Limit
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemoveDatabaseConfig removes the database config
func (s *Manager) RemoveDatabaseConfig(ctx context.Context, project, dbAlias string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update database config
	delete(projectConfig.Modules.Crud, dbAlias)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)

	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetLogicalDatabaseName gets logical database name for provided db alias
func (s *Manager) GetLogicalDatabaseName(ctx context.Context, project, dbAlias string) (string, error) {
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return "", err
	}
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get logical database name provided db alias (%s) does not exists", dbAlias), nil, nil)
	}
	return collection.DBName, nil
}

// GetPreparedQuery gets preparedQuery from config
func (s *Manager) GetPreparedQuery(ctx context.Context, project, dbAlias, id string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if dbAlias != "*" {
		databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get prepared query provided db alias (%s) does not exists", dbAlias), nil, nil)
		}

		if id != "*" {
			preparedQuery, ok := databaseConfig.PreparedQueries[id]
			if !ok {
				return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Prepared Queries with id (%s) not present in config", id), nil, nil)
			}
			return http.StatusOK, []interface{}{&preparedQueryResponse{ID: id, DBAlias: dbAlias, SQL: preparedQuery.SQL, Arguments: preparedQuery.Arguments, Rule: preparedQuery.Rule}}, nil
		}
		preparedQuery := databaseConfig.PreparedQueries
		coll := make([]interface{}, 0)
		for key, value := range preparedQuery {
			coll = append(coll, &preparedQueryResponse{ID: key, DBAlias: dbAlias, SQL: value.SQL, Arguments: value.Arguments, Rule: value.Rule})
		}
		return http.StatusOK, coll, nil
	}
	databases := projectConfig.Modules.Crud
	coll := make([]interface{}, 0)
	for alias, dbInfo := range databases {
		for key, value := range dbInfo.PreparedQueries {
			coll = append(coll, &preparedQueryResponse{ID: key, DBAlias: alias, SQL: value.SQL, Arguments: value.Arguments, Rule: value.Rule})
		}
	}
	return http.StatusOK, coll, nil
}

// SetPreparedQueries sets database preparedqueries
func (s *Manager) SetPreparedQueries(ctx context.Context, project, dbAlias, id string, v *config.PreparedQuery, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	v.ID = id
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update database PreparedQueries
	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to set prepared query provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	if databaseConfig.PreparedQueries == nil {
		databaseConfig.PreparedQueries = make(map[string]*config.PreparedQuery, 1)
	}
	databaseConfig.PreparedQueries[id] = v

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemovePreparedQueries removes the database PreparedQueries
func (s *Manager) RemovePreparedQueries(ctx context.Context, project, dbAlias, id string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to remove prepared query provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	// update database reparedQueries
	delete(databaseConfig.PreparedQueries, id)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetModifySchema modifies the schema of table
func (s *Manager) SetModifySchema(ctx context.Context, project, dbAlias, col string, v *config.TableRule, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to modify schema provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	// Modify the schema
	schemaMod := s.modules.GetSchemaModuleForSyncMan()
	if err := schemaMod.SchemaModifyAll(ctx, dbAlias, collection.DBName, map[string]*config.TableRule{col: v}); err != nil {
		return http.StatusInternalServerError, err
	}

	if collection.Collections == nil {
		collection.Collections = map[string]*config.TableRule{}
	}
	temp, ok := collection.Collections[col]
	if !ok {
		collection.Collections[col] = &config.TableRule{Schema: v.Schema, Rules: map[string]*config.Rule{}}
	} else {
		temp.Schema = v.Schema
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetCollectionRules sets the collection rules of the database
func (s *Manager) SetCollectionRules(ctx context.Context, project, dbAlias, col string, v *config.TableRule, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return s.setCollectionRules(ctx, projectConfig, project, dbAlias, col, v)
}

func (s *Manager) setCollectionRules(ctx context.Context, projectConfig *config.Project, project, dbAlias, col string, v *config.TableRule) (int, error) {
	// update collection rules & is realtime in config
	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to set collection rules provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	collection, ok := databaseConfig.Collections[col]
	if !ok {
		if databaseConfig.Collections == nil {
			databaseConfig.Collections = map[string]*config.TableRule{col: v}
		} else {
			databaseConfig.Collections[col] = &config.TableRule{IsRealTimeEnabled: v.IsRealTimeEnabled, Rules: v.Rules}
		}
	} else {
		collection.IsRealTimeEnabled = v.IsRealTimeEnabled
		collection.Rules = v.Rules
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteCollectionRules deletes the collection rules of the database
func (s *Manager) DeleteCollectionRules(ctx context.Context, project, dbAlias, col string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete collection rules, provided db alias (%s) does not exists", dbAlias), nil, nil)
	}

	collection, ok := databaseConfig.Collections[col]
	if !ok {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete collection rules, provided collection (%s) does not exists", col), nil, nil)
	}

	collection.Rules = make(map[string]*config.Rule)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetReloadSchema reloads of the schema
func (s *Manager) SetReloadSchema(ctx context.Context, dbAlias, project string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the schema module
	schemaMod := s.modules.GetSchemaModuleForSyncMan()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	collectionConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, errors.New("specified database not present in config")
	}

	colResult := map[string]interface{}{}
	for colName, colValue := range collectionConfig.Collections {
		if colName == "default" {
			continue
		}
		result, err := schemaMod.SchemaInspection(ctx, dbAlias, collectionConfig.DBName, colName)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// set new schema in config & return in response body
		colValue.Schema = result
		colResult[colName] = result
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetSchemaInspection inspects the schema
func (s *Manager) SetSchemaInspection(ctx context.Context, project, dbAlias, col, schema string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, errors.New("specified database not present in config")
	}

	if collection.Collections == nil {
		collection.Collections = map[string]*config.TableRule{}
	}
	temp, ok := collection.Collections[col]
	if !ok {
		collection.Collections[col] = &config.TableRule{Schema: schema, Rules: map[string]*config.Rule{}}
	} else {
		temp.Schema = schema
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusBadRequest, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// RemoveCollection removed the collection from the database collection schema in config
func (s *Manager) RemoveCollection(ctx context.Context, project, dbAlias, col string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return http.StatusBadRequest, errors.New("specified database not present in config")
	}

	if collection.Collections == nil {
		return http.StatusOK, nil
	}

	delete(collection.Collections, col)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting crud config", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetModifyAllSchema modifies schema of all tables
func (s *Manager) SetModifyAllSchema(ctx context.Context, dbAlias, project string, v config.CrudStub, params model.RequestParams) (int, error) {
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

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *Manager) applySchemas(ctx context.Context, project, dbAlias string, projectConfig *config.Project, v config.CrudStub) error {

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return errors.New("specified database not present in config")
	}

	if err := s.modules.GetSchemaModuleForSyncMan().SchemaModifyAll(ctx, dbAlias, collection.DBName, v.Collections); err != nil {
		return err
	}

	for colName, colValue := range v.Collections {
		temp, ok := collection.Collections[colName]
		// if collection doesn't exist then add to config
		if collection.Collections == nil {
			collection.Collections = map[string]*config.TableRule{}
		}
		if !ok {
			collection.Collections[colName] = &config.TableRule{Schema: colValue.Schema, Rules: map[string]*config.Rule{}}
		} else {
			temp.Schema = colValue.Schema
		}
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting crud config", err, nil)
	}

	return nil
}

// GetDatabaseConfig gets database config
func (s *Manager) GetDatabaseConfig(ctx context.Context, project, dbAlias string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if dbAlias != "*" {
		dbConfig, ok := projectConfig.Modules.Crud[dbAlias]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("specified dbAlias (%s) not present in config", dbAlias), nil, nil)
		}
		return http.StatusOK, []interface{}{config.Crud{dbAlias: {Enabled: dbConfig.Enabled, Conn: dbConfig.Conn, Type: dbConfig.Type, DBName: dbConfig.DBName}}}, nil
	}

	services := []interface{}{}
	for key, value := range projectConfig.Modules.Crud {
		services = append(services, config.Crud{key: {Enabled: value.Enabled, Conn: value.Conn, Type: value.Type, DBName: value.DBName}})
	}
	return http.StatusOK, services, nil
}

// GetCollectionRules gets collection rules
func (s *Manager) GetCollectionRules(ctx context.Context, project, dbAlias, col string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if dbAlias != "*" && col != "*" {
		collectionInfo, ok := projectConfig.Modules.Crud[dbAlias].Collections[col]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("specified collection (%s) not present in config for dbAlias (%s) )", dbAlias, col), nil, nil)
		}
		return http.StatusOK, []interface{}{map[string]*dbRulesResponse{fmt.Sprintf("%s-%s", dbAlias, col): {IsRealTimeEnabled: collectionInfo.IsRealTimeEnabled, Rules: collectionInfo.Rules}}}, nil
	} else if dbAlias != "*" {
		collections := projectConfig.Modules.Crud[dbAlias].Collections
		coll := map[string]*dbRulesResponse{}
		for key, value := range collections {
			coll[fmt.Sprintf("%s-%s", dbAlias, key)] = &dbRulesResponse{IsRealTimeEnabled: value.IsRealTimeEnabled, Rules: value.Rules}
		}
		return http.StatusOK, []interface{}{coll}, nil
	}
	databases := projectConfig.Modules.Crud
	coll := map[string]*dbRulesResponse{}
	for dbName, dbInfo := range databases {
		for key, value := range dbInfo.Collections {
			coll[fmt.Sprintf("%s-%s", dbName, key)] = &dbRulesResponse{IsRealTimeEnabled: value.IsRealTimeEnabled, Rules: value.Rules}
		}
	}
	return http.StatusOK, []interface{}{coll}, nil
}

// GetSchemas gets schemas from config
func (s *Manager) GetSchemas(ctx context.Context, project, dbAlias, col, format string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	a := s.modules.GetSchemaModuleForSyncMan()
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
