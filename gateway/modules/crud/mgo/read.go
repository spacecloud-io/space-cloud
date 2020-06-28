package mgo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Read querys document(s) from the database
func (m *Mongo) Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, error) {
	collection := m.client.Database(m.dbName).Collection(col)

	switch req.Operation {
	case utils.Count:
		countOptions := options.Count()

		count, err := collection.CountDocuments(ctx, req.Find, countOptions)
		if err != nil {
			return 0, nil, err
		}

		return count, count, nil

	case utils.Distinct:
		distinct := req.Options.Distinct
		if distinct == nil {
			return 0, nil, utils.ErrInvalidParams
		}

		result, err := collection.Distinct(ctx, *distinct, req.Find)
		if err != nil {
			return 0, nil, err
		}

		// convert result []string to []map[string]interface
		finalResult := []interface{}{}
		for _, value := range result {
			doc := map[string]interface{}{}
			doc[*distinct] = value
			finalResult = append(finalResult, doc)
		}

		return int64(len(result)), finalResult, nil

	case utils.All:
		findOptions := options.Find()

		if req.Options != nil {
			if req.Options.Select != nil {
				findOptions = findOptions.SetProjection(req.Options.Select)
			}

			if req.Options.Skip != nil {
				findOptions = findOptions.SetSkip(*req.Options.Skip)
			}

			if req.Options.Limit != nil {
				findOptions = findOptions.SetLimit(*req.Options.Limit)
			}

			if req.Options.Sort != nil {
				findOptions = findOptions.SetSort(generateSortOptions(req.Options.Sort))
			}
		}

		results := []interface{}{}
		cur, err := collection.Find(ctx, req.Find, findOptions)
		if err != nil {
			return 0, nil, err
		}
		defer func() { _ = cur.Close(ctx) }()

		var count int64
		// Finding multiple documents returns a cursor
		// Iterating through the cursor allows us to decode documents one at a time
		for cur.Next(ctx) {
			// Increment the counter
			count++

			// Read the document
			var doc map[string]interface{}
			err := cur.Decode(&doc)
			if err != nil {
				return 0, nil, err
			}

			results = append(results, doc)
		}

		if err := cur.Err(); err != nil {
			return 0, nil, err
		}

		return count, results, nil

	case utils.One:
		findOneOptions := options.FindOne()

		if req.Options != nil {
			if req.Options.Select != nil {
				findOneOptions = findOneOptions.SetProjection(req.Options.Select)
			}

			if req.Options.Skip != nil {
				findOneOptions = findOneOptions.SetSkip(*req.Options.Skip)
			}

			if req.Options.Sort != nil {
				findOneOptions = findOneOptions.SetSort(generateSortOptions(req.Options.Sort))
			}
		}

		var res map[string]interface{}
		err := collection.FindOne(ctx, req.Find, findOneOptions).Decode(&res)
		if err != nil {
			return 0, nil, err
		}

		return 1, res, nil

	default:
		return 0, nil, utils.ErrInvalidParams
	}
}

func generateSortOptions(array []string) bson.D {
	sort := bson.D{}
	for _, value := range array {
		if strings.HasPrefix(value, "-") {
			sort = append(sort, primitive.E{Key: strings.TrimPrefix(value, "-"), Value: -1})
		} else {
			sort = append(sort, primitive.E{Key: value, Value: 1})
		}
	}

	return sort
}
