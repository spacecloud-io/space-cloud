package syncman

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/utils"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
)

// SetDeleteCollection deletes a collection from the database
func (s *Manager) SetDeleteCollection(ctx context.Context, project, dbAlias, col string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// remove collection from config
	coll, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return errors.New("specified database not present in config")
	}
	delete(coll.Collections, col)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetDatabaseConnection sets the database connection
func (s *Manager) SetDatabaseConnection(ctx context.Context, project, dbAlias string, v config.CrudStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	coll, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		projectConfig.Modules.Crud[dbAlias] = &config.CrudStub{Conn: v.Conn, Enabled: v.Enabled, Collections: map[string]*config.TableRule{}, Type: v.Type, DBName: v.DBName}
	} else {
		coll.Conn = v.Conn
		coll.Enabled = v.Enabled
		coll.Type = v.Type
		// coll.Name = v.Name// TODO CHECK IF THIS IS REQUIRED
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// RemoveDatabaseConfig removes the database config
func (s *Manager) RemoveDatabaseConfig(ctx context.Context, project, dbAlias string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update database config
	delete(projectConfig.Modules.Crud, dbAlias)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// GetLogicalDatabaseName gets logical database name for provided db alias
func (s *Manager) GetLogicalDatabaseName(ctx context.Context, project, dbAlias string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return "", err
	}
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return "", errors.New("specified database not present in config")
	}
	return collection.DBName, nil
}

// GetPreparedQuery gets preparedQuery from config
func (s *Manager) GetPreparedQuery(ctx context.Context, project, dbAlias, id string) ([]interface{}, error) {
	// Acquire a lock
	type response struct {
		ID        string   `json:"id"`
		SQL       string   `json:"sql"`
		Arguments []string `json:"arguments" yaml:"arguments"`
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}

	if dbAlias != "" {
		databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
		if !ok {
			return nil, fmt.Errorf("specified database (%s) not present in config", dbAlias)
		}

		if id != "" {
			preparedQuery, ok := databaseConfig.PreparedQueries[id]
			if !ok {
				return nil, fmt.Errorf("Prepared Queries for id (%s) not present in config for dbAlias (%s) )", id, dbAlias)
			}
			return []interface{}{&response{ID: id, SQL: preparedQuery.SQL, Arguments: preparedQuery.Arguments}}, nil
		}
		preparedQuery := databaseConfig.PreparedQueries
		var coll []interface{} = make([]interface{}, 0)
		for key, value := range preparedQuery {
			coll = append(coll, &response{ID: key, SQL: value.SQL, Arguments: value.Arguments})
		}
		return coll, nil
	}
	databases := projectConfig.Modules.Crud
	var coll []interface{} = make([]interface{}, 0)
	for _, dbInfo := range databases {
		for key, value := range dbInfo.PreparedQueries {
			coll = append(coll, &response{ID: key, SQL: value.SQL, Arguments: value.Arguments})
		}
	}
	return coll, nil
}

// SetPreparedQueries sets database preparedqueries
func (s *Manager) SetPreparedQueries(ctx context.Context, project, dbAlias, id string, v *config.PreparedQuery) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	v.ID = id
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update database PreparedQueries
	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return fmt.Errorf("specified database (%s) not present in config", dbAlias)
	}

	if databaseConfig.PreparedQueries == nil {
		databaseConfig.PreparedQueries = make(map[string]*config.PreparedQuery, 1)
	}
	databaseConfig.PreparedQueries[id] = v

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// RemovePreparedQueries removes the database PreparedQueries
func (s *Manager) RemovePreparedQueries(ctx context.Context, project, dbAlias, id string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return fmt.Errorf("specified database (%s) not present in config", dbAlias)
	}

	// update database reparedQueries
	delete(databaseConfig.PreparedQueries, id)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetModifySchema modifies the schema of table
func (s *Manager) SetModifySchema(ctx context.Context, project, dbAlias, col, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return errors.New("specified database not present in config")
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
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetCollectionRules sets the collection rules of the database
func (s *Manager) SetCollectionRules(ctx context.Context, project, dbAlias, col string, v *config.TableRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	// update collection rules & is realtime in config
	databaseConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return errors.New("specified database not present in config")
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
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetReloadSchema reloads of the schema
func (s *Manager) SetReloadSchema(ctx context.Context, dbAlias, project string, schemaArg *schema.Schema) (map[string]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}

	collectionConfig, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return nil, errors.New("specified database not present in config")
	}
	colResult := map[string]interface{}{}
	for colName, colValue := range collectionConfig.Collections {
		if colName == "default" {
			continue
		}
		result, err := schemaArg.SchemaInspection(ctx, dbAlias, collectionConfig.DBName, colName)
		if err != nil {
			return nil, err
		}

		// set new schema in config & return in response body
		colValue.Schema = result
		colResult[colName] = result
	}

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return nil, err
	}

	return colResult, s.setProject(ctx, projectConfig)
}

// SetSchemaInspection inspects the schema
func (s *Manager) SetSchemaInspection(ctx context.Context, project, dbAlias, col, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return errors.New("specified database not present in config")
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
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// RemoveSchemaInspection removed the collection from the database collection schema in config
func (s *Manager) RemoveSchemaInspection(ctx context.Context, project, dbAlias, col string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		return errors.New("specified database not present in config")
	}

	if collection.Collections == nil {
		return nil
	}

	delete(collection.Collections, col)

	if err := s.modules.SetCrudConfig(project, projectConfig.Modules.Crud); err != nil {
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetModifyAllSchema modifies schema of all tables
func (s *Manager) SetModifyAllSchema(ctx context.Context, dbAlias, project string, v config.CrudStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	if err := s.applySchemas(ctx, project, dbAlias, projectConfig, v); err != nil {
		return err
	}

	return s.setProject(ctx, projectConfig)
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
		logrus.Errorf("error setting crud config - %s", err.Error())
		return err
	}

	return nil
}

// GetDatabaseConfig gets database config
func (s *Manager) GetDatabaseConfig(ctx context.Context, project, dbAlias string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if dbAlias != "" {
		dbConfig, ok := projectConfig.Modules.Crud[dbAlias]
		if !ok {
			return nil, fmt.Errorf("specified dbAlias (%s) not present in config", dbAlias)
		}
		return []interface{}{config.Crud{dbAlias: {Enabled: dbConfig.Enabled, Conn: dbConfig.Conn, Type: dbConfig.Type}}}, nil
	}

	services := []interface{}{}
	for key, value := range projectConfig.Modules.Crud {
		services = append(services, config.Crud{key: {Enabled: value.Enabled, Conn: value.Conn, Type: value.Type}})
	}
	return services, nil
}

// GetCollectionRules gets collection rules
func (s *Manager) GetCollectionRules(ctx context.Context, project, dbAlias, col string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	type response struct {
		IsRealTimeEnabled bool                    `json:"isRealtimeEnabled"`
		Rules             map[string]*config.Rule `json:"rules"`
	}
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if dbAlias != "" && col != "" {
		collectionInfo, ok := projectConfig.Modules.Crud[dbAlias].Collections[col]
		if !ok {
			return nil, fmt.Errorf("specified collection (%s) not present in config for dbAlias (%s) )", dbAlias, col)
		}
		return []interface{}{map[string]*response{fmt.Sprintf("%s-%s", dbAlias, col): {IsRealTimeEnabled: collectionInfo.IsRealTimeEnabled, Rules: collectionInfo.Rules}}}, nil
	} else if dbAlias != "" {
		collections := projectConfig.Modules.Crud[dbAlias].Collections
		coll := map[string]*response{}
		for key, value := range collections {
			coll[fmt.Sprintf("%s-%s", dbAlias, key)] = &response{IsRealTimeEnabled: value.IsRealTimeEnabled, Rules: value.Rules}
		}
		return []interface{}{coll}, nil
	}
	databases := projectConfig.Modules.Crud
	coll := map[string]*response{}
	for dbName, dbInfo := range databases {
		for key, value := range dbInfo.Collections {
			coll[fmt.Sprintf("%s-%s", dbName, key)] = &response{IsRealTimeEnabled: value.IsRealTimeEnabled, Rules: value.Rules}
		}
	}
	return []interface{}{coll}, nil
}

// GetSchemas gets schemas from config
func (s *Manager) GetSchemas(ctx context.Context, project, dbAlias, col string) ([]interface{}, error) {
	// Acquire a lock
	type response struct {
		Schema string `json:"schema"`
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if dbAlias != "" && col != "" {
		collectionInfo, ok := projectConfig.Modules.Crud[dbAlias].Collections[col]
		if !ok {
			return nil, fmt.Errorf("collection (%s) not present in config for dbAlias (%s) )", dbAlias, col)
		}
		return []interface{}{map[string]*response{fmt.Sprintf("%s-%s", dbAlias, col): {Schema: collectionInfo.Schema}}}, nil
	} else if dbAlias != "" {
		collections := projectConfig.Modules.Crud[dbAlias].Collections
		coll := map[string]*response{}
		for key, value := range collections {
			coll[fmt.Sprintf("%s-%s", dbAlias, key)] = &response{Schema: value.Schema}
		}
		return []interface{}{coll}, nil
	}
	databases := projectConfig.Modules.Crud
	coll := map[string]*response{}
	for dbName, dbInfo := range databases {
		for key, value := range dbInfo.Collections {
			coll[fmt.Sprintf("%s-%s", dbName, key)] = &response{Schema: value.Schema}
		}
	}
	return []interface{}{coll}, nil
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
		return "", utils.LogError("cannot get internal access token", "syncman", "GetSecrets", err)
	}

	// makes http request to get secrets from runner
	var vPtr result
	url := fmt.Sprintf("http://%s/v1/runner/%s/secrets?id=%s", s.runnerAddr, project, secretName)
	if err := s.MakeHTTPRequest(ctx, "GET", url, token, "", map[string]interface{}{}, &vPtr); err != nil {
		return "", err
	}

	return vPtr.Result[0].Data[key], nil
}
