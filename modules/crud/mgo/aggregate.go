package mgo

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Aggregate performs a mongo db pipeline aggregation
func (m *Mongo) Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error) {
	collection := m.client.Database(project).Collection(col)

	switch req.Operation {
	case utils.One:
		var result map[string]interface{}

		cur, err := collection.Aggregate(ctx, req.Pipeline)
		if err != nil {
			return nil, err
		}
		defer cur.Close(ctx)

		if !cur.Next(ctx) {
			return nil, errors.New("No result found")
		}

		err = cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		return result, nil

	case utils.All:
		results := []interface{}{}

		cur, err := collection.Aggregate(ctx, req.Pipeline)
		defer cur.Close(ctx)
		if err != nil {
			return nil, err
		}

		for cur.Next(ctx) {
			var doc map[string]interface{}
			err := cur.Decode(&doc)
			if err != nil {
				return nil, err
			}

			results = append(results, doc)
		}

		if err := cur.Err(); err != nil {
			return nil, err
		}

		return results, nil

	default:
		return nil, utils.ErrInvalidParams
	}
}
