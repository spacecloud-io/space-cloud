package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Update updates the document(s) which match the condition provided.
func (m *Mongo) Update(ctx context.Context, project, col string, req *model.UpdateRequest) (int64, error) {
	collection := m.client.Database(project).Collection(col)

	switch req.Operation {
	case utils.One:
		_, err := collection.UpdateOne(ctx, req.Find, req.Update)
		if err != nil {
			return 0, err
		}

		return 1, nil

	case utils.All:
		res, err := collection.UpdateMany(ctx, req.Find, req.Update)
		if err != nil {
			return 0, err
		}

		return res.MatchedCount, nil

	case utils.Upsert:
		doUpsert := true
		res, err := collection.UpdateOne(ctx, req.Find, req.Update, &options.UpdateOptions{Upsert: &doUpsert})
		if err != nil {
			return 0, err
		}

		return res.MatchedCount + res.UpsertedCount, nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
