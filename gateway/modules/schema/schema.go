package schema

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	schemaHelpers "github.com/spaceuptech/space-cloud/gateway/modules/schema/helpers"
)

// Schema data stucture for schema package
type Schema struct {
	lock      sync.RWMutex
	SchemaDoc model.Type
	crud      model.CrudSchemaInterface
	project   string
	dbSchemas config.DatabaseSchemas
	clusterID string
}

// Init creates a new instance of the schema object
func Init(clusterID string, crud model.CrudSchemaInterface) *Schema {
	return &Schema{clusterID: clusterID, SchemaDoc: model.Type{}, crud: crud}
}

// SetDatabaseSchema modifies the tables according to the schema on save
func (s *Schema) SetDatabaseSchema(c config.DatabaseSchemas, project string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.dbSchemas = c
	s.project = project
	if err := s.parseSchema(c); err != nil {
		return err
	}

	return nil
}

// GetSchema function gets schema
func (s *Schema) GetSchema(dbAlias, col string) (model.Fields, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	dbSchema, p := s.SchemaDoc[dbAlias]
	if !p {
		return nil, false
	}

	colSchema, p := dbSchema[col]
	if !p {
		return nil, false
	}

	fields := make(model.Fields, len(colSchema))
	for k, v := range colSchema {
		fields[k] = v
	}

	return fields, true
}

// parseSchema Initializes Schema field in Module struct
func (s *Schema) parseSchema(crud config.DatabaseSchemas) error {
	schema, err := schemaHelpers.Parser(crud)
	if err != nil {
		return err
	}
	s.SchemaDoc = schema
	return nil
}
