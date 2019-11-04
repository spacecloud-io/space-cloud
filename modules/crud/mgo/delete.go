package mgo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Delete removes the document(s) from the database which match the condition
func (m *Mongo) Delete(ctx context.Context, project, col string, req *model.DeleteRequest) (int64, error) {
	collection := m.client.Database(project).Collection(col)

	switch req.Operation {
	case utils.One:
		_, err := collection.DeleteOne(ctx, req.Find)
		if err != nil {
			return 0, err
		}

		return 1, nil

	case utils.All:
		res, err := collection.DeleteMany(ctx, req.Find)
		if err != nil {
			return 0, err
		}

		return res.DeletedCount, nil

	default:
		return 0, errors.New("Invalid operation")
	}
}

// DeleteCollection removes a collection from database`
func (m Mongo) DeleteCollection(ctx context.Context, project, col string) error {
	return m.client.Database(project).Collection(col, &options.CollectionOptions{}).Drop(ctx)
}
