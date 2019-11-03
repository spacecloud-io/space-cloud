package mgo

import (
	"context"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Create inserts a document (or multiple when op is "all") into the database
func (m *Mongo) Create(ctx context.Context, project, col string, req *model.CreateRequest) (int64, error) {
	// Create a collection object
	collection := m.client.Database(project).Collection(col)

	switch req.Operation {
	case utils.One:
		// Insert single document
		_, err := collection.InsertOne(ctx, req.Document)
		if err != nil {
			return 0, err
		}

		return 1, nil

	case utils.All:
		// Insert multiple documents
		objs, ok := req.Document.([]interface{})
		if !ok {
			return 0, utils.ErrInvalidParams
		}

		res, err := collection.InsertMany(ctx, objs)
		if err != nil {
			return 0, err
		}

		return int64(len(res.InsertedIDs)), nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
