package schema

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"
)

// GetSchemaForDB gets schema of specified database & collection
// If * is provided for database or collection. It will get all the databases and collection
func (s *Schema) GetSchemaForDB(ctx context.Context, dbAlias, col, format string) ([]interface{}, error) {
	s.lock.RLock()
	defer s.lock.RLock()

	alreadyAddedTables := map[string]bool{}
	schemaResponse := make([]interface{}, 0)
	if dbAlias != "*" && col != "*" {
		db, ok := s.config[dbAlias]
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Provided database doesn't exists (%s)", dbAlias), nil, nil)
		}
		if err := s.getSchemaResponse(ctx, format, dbAlias, col, true, alreadyAddedTables, db.Collections, &schemaResponse); err != nil {
			return nil, err
		}
	} else if dbAlias != "*" {
		collections := s.config[dbAlias].Collections
		for key := range collections {
			if err := s.getSchemaResponse(ctx, format, dbAlias, key, false, alreadyAddedTables, collections, &schemaResponse); err != nil {
				return nil, err
			}
		}
	} else {
		for dbName, dbInfo := range s.config {
			for key := range dbInfo.Collections {
				if err := s.getSchemaResponse(ctx, format, dbName, key, false, alreadyAddedTables, dbInfo.Collections, &schemaResponse); err != nil {
					return nil, err
				}
			}
		}
	}
	return schemaResponse, nil
}
