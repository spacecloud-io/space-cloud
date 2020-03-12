package syncman

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
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

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

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

	// update database config
	coll, ok := projectConfig.Modules.Crud[dbAlias]
	if !ok {
		projectConfig.Modules.Crud[dbAlias] = &config.CrudStub{Conn: v.Conn, Enabled: v.Enabled, Collections: map[string]*config.TableRule{}, Type: v.Type}
	} else {
		coll.Conn = v.Conn
		coll.Enabled = v.Enabled
		coll.Type = v.Type
	}

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

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

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

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
	temp, ok := collection.Collections[col]
	// if collection doesn't exist then add to config
	if !ok {
		collection.Collections[col] = &config.TableRule{Schema: schema, Rules: map[string]*config.Rule{}} // TODO: rule field here is null
	} else {
		temp.Schema = schema
	}

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

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

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

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
		result, err := schemaArg.SchemaInspection(ctx, dbAlias, project, colName)
		if err != nil {
			return nil, err
		}

		// set new schema in config & return in response body
		colValue.Schema = result
		colResult[colName] = result
	}

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

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

	temp, ok := collection.Collections[col]
	// if collection doesn't exist then add to config
	if !ok {
		collection.Collections[col] = &config.TableRule{Schema: schema, Rules: map[string]*config.Rule{}} // TODO: rule field here is null
	} else {
		temp.Schema = schema
	}

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

	return s.setProject(ctx, projectConfig)
}

// SetModifyAllSchema modifies schema of all tables
func (s *Manager) SetModifyAllSchema(ctx context.Context, dbAlias, project string, schemaArg model.SchemaEventingInterface, v config.CrudStub) error {
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

	if err := schemaArg.SchemaModifyAll(ctx, dbAlias, project, v.Collections); err != nil {
		return err
	}

	for colName, colValue := range v.Collections {
		temp, ok := collection.Collections[colName]
		// if collection doesn't exist then add to config
		if !ok {
			collection.Collections[colName] = &config.TableRule{Schema: colValue.Schema} // TODO: rule field here is null
		} else {
			temp.Schema = colValue.Schema
		}
	}

	s.modules.SetCrudConfig(project, projectConfig.Secret, projectConfig.AESkey, projectConfig.Modules.Crud, projectConfig.Modules.FileStore, projectConfig.Modules.Services, &projectConfig.Modules.Eventing)

	return s.setProject(ctx, projectConfig)
}
