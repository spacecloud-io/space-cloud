package mgo

import (
	"context"
	"fmt"
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

		pipeline := make([]bson.M, 0)
		for function, colArray := range req.Aggregate {
			for _, column := range colArray {
				asColumnName := generateAggregateAsColumnName(function, column)
				switch function {
				case "sum":
					matchStage := getMatchStage(req.Find)
					if matchStage != nil {
						pipeline = append(pipeline, matchStage)
					}
					groupStage := getGroupByStage(req.GroupBy, asColumnName, column)
					pipeline = append(pipeline, groupStage)
				default:
					return 0, nil, utils.LogError(fmt.Sprintf(`Unknown aggregate funcion "%s"`, function), "mgo", "Read", nil)
				}
			}
		}

		results := []interface{}{}
		cur, err := collection.Aggregate(ctx, pipeline)
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

			resultObj := make(map[string]interface{})
			for key, value := range doc {
				v := strings.Split(key, "__")
				if len(v) != 3 || !strings.HasPrefix(key, utils.GraphQLAggregate) {
					resultObj[v[0]] = value
					continue
				}
				resultObj[v[0]] = map[string]interface{}{v[1]: map[string]interface{}{v[2]: value}}
			}

			results = append(results, resultObj)
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

func getMatchStage(find map[string]interface{}) bson.M {
	if len(find) > 0 {
		matchStage := bson.M{
			"$match": find,
		}
		return matchStage
	}
	return nil
}

func getGroupByStage(groupBy []interface{}, asColumnName, column string) bson.M {
	if len(groupBy) > 0 {
		groupStage := bson.M{}
		groupByMap := make(map[string]interface{})
		for _, val := range groupBy {
			groupByMap[fmt.Sprintf("%v", val)] = fmt.Sprintf("$%v", val)
		}
		groupStage = bson.M{
			"$group": bson.M{
				"_id": groupByMap,
				asColumnName: bson.M{
					"$sum": fmt.Sprintf("$%s", column),
				},
			},
		}
		return groupStage
	}
	return bson.M{"$group": bson.M{"_id": bson.M{}}}
}

func generateAggregateAsColumnName(function, column string) string {
	return fmt.Sprintf("%s__%s__%s", utils.GraphQLAggregate, function, column)
}
