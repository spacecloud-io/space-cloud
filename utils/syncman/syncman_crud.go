package syncman

import (
	"context"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth/schema"
)

func (s *Manager) SetDeleteCollection(projectConfig *config.Project, dbType, col string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// remove collection from config
	coll := projectConfig.Modules.Crud[dbType]
	delete(coll.Collections, col)

	return s.setProject(projectConfig)
}

func (s *Manager) SetDatabaseConnection(projectConfig *config.Project, dbType string, connection string, enabled bool) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// update database config
	coll := projectConfig.Modules.Crud[dbType]
	coll.Conn = connection
	coll.Enabled = enabled

	return s.setProject(projectConfig)
}

func (s *Manager) SetModifySchema(projectConfig *config.Project, dbType, col, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// update schema in config
	collection := projectConfig.Modules.Crud[dbType]
	temp, ok := collection.Collections[col]
	// if collection doesn't exist then add to config
	if !ok {
		collection.Collections[col] = &config.TableRule{} // TODO: rule field here is null
	}
	temp.Schema = schema

	return s.setProject(projectConfig)
}

func (s *Manager) SetCollectionRules(projectConfig *config.Project, dbType, col string, v *config.TableRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// update collection rules & is realtime in config
	collection, ok := projectConfig.Modules.Crud[dbType].Collections[col]
	if !ok {
		projectConfig.Modules.Crud[dbType].Collections[col] = v
	} else {
		collection.IsRealTimeEnabled = v.IsRealTimeEnabled
		collection.Rules = v.Rules
	}
	return s.setProject(projectConfig)
}

func (s *Manager) SetReloadSchema(ctx context.Context, projectConfig *config.Project, dbType, col, project string, schemaArg *schema.Schema) (map[string]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	collectionConfig := projectConfig.Modules.Crud[dbType]
	colResult := map[string]interface{}{}
	for colName, colValue := range collectionConfig.Collections {
		result, err := schemaArg.SchemaInspection(ctx, dbType, project, col)
		if err != nil {
			return nil, err
		}
		// set new schema in config & return in response body
		colValue.Schema = result
		colResult[colName] = result
	}

	return colResult, s.setProject(projectConfig)
}

func (s *Manager) SetSchemaInspection(projectConfig *config.Project, dbType, col, schema string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// update schema in config
	coll := projectConfig.Modules.Crud[dbType]
	coll.Collections[col].Schema = schema

	return s.setProject(projectConfig)
}

func (s *Manager) SetModifyAllSchema(ctx context.Context, projectConfig *config.Project, dbType, project string, schemaArg *schema.Schema, v config.CrudStub) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// update schema in config
	collection := projectConfig.Modules.Crud[dbType]

	for colName, colValue := range v.Collections {
		parsedColValue, err := schemaArg.Inspector(ctx, dbType, project, colName)
		if err != nil {
			return err
		}
		schemaArg.SchemaJoin(ctx, parsedColValue, dbType, colName, project, v)
		temp, ok := collection.Collections[colName]
		// if collection doesn't exist then add to config
		if !ok {
			collection.Collections[colName] = &config.TableRule{} // TODO: rule field here is null
		}
		temp.Schema = colValue.Schema
	}

	return s.setProject(projectConfig)
}
