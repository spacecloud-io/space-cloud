package mgo

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetCollections returns collection / tables name of specified database
func (m *Mongo) GetCollections(ctx context.Context) ([]utils.DatabaseCollections, error) {

	collections, err := m.getClient().Database(m.dbName).ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to query database to get tables in database (%s)", m.dbName), err, nil)
	}

	dbCols := make([]utils.DatabaseCollections, len(collections))
	for i, col := range collections {
		dbCols[i] = utils.DatabaseCollections{TableName: col}
	}

	return dbCols, nil
}
