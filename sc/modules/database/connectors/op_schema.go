package connectors

import (
	"context"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/database/connectors/schema"
	"github.com/spaceuptech/helpers"
)

// InspectCollectionSchema generates a schema object based on decription of a table
func (m *Module) InspectCollectionSchema(ctx context.Context, col string) (model.CollectionSchemas, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return nil, err
	}

	// Get the description of the table from the database
	fields, indexes, err := m.connector.DescribeTable(ctx, col)
	if err != nil {
		return nil, err
	}
	return schema.ParseCollectionDescription(m.dbConfig.Type, col, fields, indexes, m.schemaDoc)
}

// ApplyCollectionSchema creates or alters tables of a sql database.
func (m *Module) ApplyCollectionSchema(ctx context.Context, tableName string, newSchema model.CollectionSchemas) error {
	// Return gracefully if db type is mongo
	if m.dbConfig.Type == string(model.Mongo) || m.dbConfig.Type == string(model.EmbeddedDB) {
		return nil
	}

	// Load the current schema
	currentSchema, err := m.InspectCollectionSchema(ctx, tableName)
	if err != nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Schema Inspector Error", map[string]interface{}{"error": err.Error()})
	}

	// Prepare creation queries to run as batch
	queries, err := schema.PrepareCreationQueries(ctx, m.dbConfig.Type, tableName, m.dbConfig.DBName, newSchema, currentSchema, m.ApplyCollectionSchema)
	if err != nil {
		return err
	}

	return m.RawBatch(ctx, queries)
}
