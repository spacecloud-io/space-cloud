package mgo

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetCollections returns collection / tables name of specified database
func (m *Mongo) GetCollections(ctx context.Context) ([]utils.DatabaseCollections, error) {

	collections, err := m.client.Database(m.dbName).ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	dbCols := make([]utils.DatabaseCollections, len(collections))
	for i, col := range collections {
		dbCols[i] = utils.DatabaseCollections{TableName: col}
	}

	return dbCols, nil
}
