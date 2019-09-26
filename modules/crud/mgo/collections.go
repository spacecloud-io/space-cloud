package mgo

import (
	"context"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetCollections returns collection / tables name of specified database
func (m *Mongo) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {

	collections, err := m.client.Database(project).ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	dbCols := make([]utils.DatabaseCollections, len(collections))
	for i, col := range collections {
		dbCols[i] = utils.DatabaseCollections{TableName: col}
	}

	return dbCols, nil
}
