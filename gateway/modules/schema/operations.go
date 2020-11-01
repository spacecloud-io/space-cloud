package schema

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// GetSchemaForDB gets schema of specified database & collection
// If * is provided for database or collection. It will get all the databases and collection
func (s *Schema) GetSchemaForDB(ctx context.Context, dbAlias, col, format string) ([]interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	alreadyAddedTables := map[string]bool{}
	schemaResponse := make([]interface{}, 0)
	if dbAlias != "*" && col != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, s.project, config.ResourceDatabaseSchema, dbAlias, col)
		_, ok := s.dbSchemas[resourceID]
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("collection (%s) not present in config for dbAlias (%s) )", dbAlias, col), nil, nil)
		}
		if err := s.getSchemaResponse(ctx, format, dbAlias, col, true, alreadyAddedTables, &schemaResponse); err != nil {
			return nil, err
		}
	} else if dbAlias != "*" {
		for _, dbSchema := range s.dbSchemas {
			if err := s.getSchemaResponse(ctx, format, dbAlias, dbSchema.Table, false, alreadyAddedTables, &schemaResponse); err != nil {
				return nil, err
			}
		}
	} else {
		for _, dbSchema := range s.dbSchemas {
			if err := s.getSchemaResponse(ctx, format, dbSchema.DbAlias, dbSchema.Table, false, alreadyAddedTables, &schemaResponse); err != nil {
				return nil, err
			}
		}
	}
	return schemaResponse, nil
}
