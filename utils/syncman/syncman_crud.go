package syncman

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/schema"
)

func (s *Manager) SetDeleteCollection(project, dbType, col string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// remove collection from config
	coll, ok := projectConfig.Modules.Crud[dbType]
	if !ok {
		return errors.New("specified database not present in config")
	}
	delete(coll.Collections, col)

	return s.setProject(projectConfig)
}

func (s *Manager) SetDatabaseConnection(project, dbType string, connection string, enabled bool) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update database config
	coll, ok := projectConfig.Modules.Crud[dbType]
	if !ok {
		projectConfig.Modules.Crud[dbType] = &config.CrudStub{Conn: connection, Enabled: enabled, Collections: map[string]*config.TableRule{}}
	} else {
		coll.Conn = connection
		coll.Enabled = enabled
	}

	return s.setProject(projectConfig)
}

func (s *Manager) SetModifySchema(project, dbType, col, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbType]
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

	return s.setProject(projectConfig)
}

func (s *Manager) SetCollectionRules(project, dbType, col string, v *config.TableRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	// update collection rules & is realtime in config
	databaseConfig, ok := projectConfig.Modules.Crud[dbType]
	if !ok {
		return errors.New("specified database not present in config")
	}
	collection, ok := databaseConfig.Collections[col]
	if !ok {
		databaseConfig.Collections = map[string]*config.TableRule{col: v}
	} else {
		collection.IsRealTimeEnabled = v.IsRealTimeEnabled
		collection.Rules = v.Rules
	}
	return s.setProject(projectConfig)
}

func (s *Manager) SetReloadSchema(ctx context.Context, dbType, project string, schemaArg *schema.Schema) (map[string]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}

	collectionConfig, ok := projectConfig.Modules.Crud[dbType]
	if !ok {
		return nil, errors.New("specified database not present in config")
	}
	colResult := map[string]interface{}{}
	for colName, colValue := range collectionConfig.Collections {
		if colName == "default" {
			continue
		}
		result, err := schemaArg.SchemaInspection(ctx, dbType, project, colName)
		if err != nil {
			return nil, err
		}

		// set new schema in config & return in response body
		colValue.Schema = result
		colResult[colName] = result
	}

	return colResult, s.setProject(projectConfig)
}

func (s *Manager) SetSchemaInspection(project, dbType, col, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbType]
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

	return s.setProject(projectConfig)
}

func (s *Manager) SetModifyAllSchema(ctx context.Context, dbType, project string, schemaArg *schema.Schema, v config.CrudStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// update schema in config
	collection, ok := projectConfig.Modules.Crud[dbType]
	if !ok {
		return errors.New("specified database not present in config")
	}

	if err := schemaArg.SchemaModifyAll(ctx, dbType, project, v.Collections); err != nil {
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

	return s.setProject(projectConfig)
}
