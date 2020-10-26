package mgo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Read queries document(s) from the database
func (m *Mongo) Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, error) {
	if req.Options != nil && len(req.Options.Join) > 0 {
		return 0, nil, errors.New("cannot perform joins in mongo db")
	}
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
		if len(req.Aggregate) > 0 {
			if len(req.Find) > 0 {
				pipeline = append(pipeline, bson.M{"$match": req.Find})
			}
			sortFields := make([]string, 0)
			functionsMap := make(bson.M)
			for function, colArray := range req.Aggregate {
				for _, column := range colArray {
					asColumnName := getAggregateAsColumnName(function, column)
					switch function {
					case "sum":
						getGroupByStageFunctionsMap(functionsMap, asColumnName, function, getAggregateColumnName(column))
					case "min":
						getGroupByStageFunctionsMap(functionsMap, asColumnName, function, getAggregateColumnName(column))
					case "max":
						getGroupByStageFunctionsMap(functionsMap, asColumnName, function, getAggregateColumnName(column))
					case "avg":
						getGroupByStageFunctionsMap(functionsMap, asColumnName, function, getAggregateColumnName(column))
					case "count":
						getGroupByStageFunctionsMap(functionsMap, asColumnName, function, "*")
					default:
						return 0, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf(`Unknown aggregate funcion %s`, function), nil, map[string]interface{}{})
					}
					for _, field := range req.Options.Sort {
						if sortValue := generateSortFields(field, column, asColumnName); sortValue != "" {
							sortFields = append(sortFields, sortValue)
						}
					}
				}
			}
			groupStage, sortArr := createGroupByStage(functionsMap, req.GroupBy, req.Options.Sort)
			sortFields = append(sortFields, sortArr...)
			pipeline = append(pipeline, groupStage)
			if req.Options != nil {
				pipeline = append(pipeline, getOptionStage(req.Options, sortFields)...)
			}
		}

		var cur *mongo.Cursor
		var err error
		results := []interface{}{}

		if len(req.Aggregate) > 0 {
			helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Mongo aggregate", map[string]interface{}{"pipeline": pipeline})
			cur, err = collection.Aggregate(ctx, pipeline)
		} else {
			helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Mongo query", map[string]interface{}{"find": req.Find, "options": findOptions})
			cur, err = collection.Find(ctx, req.Find, findOptions)
		}
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

			if len(req.Aggregate) > 0 {
				getNestedObject(doc)
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

func generateSortFields(sortColumn, currentColumn, newColumnName string) string {
	isDescending := false
	if strings.HasPrefix(sortColumn, "-") {
		isDescending = true
		sortColumn = strings.TrimPrefix(sortColumn, "-")
	}
	if sortColumn == currentColumn {
		if isDescending {
			return "-" + newColumnName
		}
		return newColumnName
	}
	return ""
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

func getGroupByStageFunctionsMap(functionsMap bson.M, asColumnName, function, column string) {
	if column != "*" {
		functionsMap[asColumnName] = bson.M{
			fmt.Sprintf("$%s", function): fmt.Sprintf("$%s", column),
		}
	} else {
		functionsMap[asColumnName] = bson.M{
			"$sum": 1,
		}
	}
}

func createGroupByStage(functionsMap bson.M, groupBy []interface{}, sort []string) (bson.M, []string) {
	groupByMap := make(map[string]interface{})
	groupStage := bson.M{
		"$group": bson.M{"_id": bson.M{}},
	}
	sortArr := make([]string, 0)
	if len(groupBy) > 0 {
		for _, val := range groupBy {
			key := fmt.Sprintf("%v", val)
			value := fmt.Sprintf("$%v", val)
			groupByMap[key] = value
			for _, sortKey := range sort {
				if sortValue := generateSortFields(sortKey, key, "_id."+key); sortValue != "" {
					sortArr = append(sortArr, sortValue)
				}
			}
		}
		groupStage["$group"].(bson.M)["_id"] = groupByMap
	}
	for key, value := range functionsMap {
		groupStage["$group"].(bson.M)[key] = value
	}
	return groupStage, sortArr
}

func generateAggregateSortOptions(array []string) bson.M {
	sort := bson.M{}
	for _, value := range array {
		if strings.HasPrefix(value, "-") {
			sort[strings.TrimPrefix(value, "-")] = -1
		} else {
			sort[value] = 1
		}
	}

	return sort
}

func getOptionStage(options *model.ReadOptions, sortFields []string) []bson.M {
	var optionStage []bson.M

	if options.Skip != nil {
		// NOTE: we are sorting the result before skip operation to give a consistent $skip result
		optionStage = append(optionStage, bson.M{"$sort": generateAggregateSortOptions([]string{"_id"})})
		optionStage = append(optionStage, bson.M{"$skip": options.Skip})
	}
	if options.Limit != nil {
		optionStage = append(optionStage, bson.M{"$sort": generateAggregateSortOptions([]string{"_id"})})
		optionStage = append(optionStage, bson.M{"$limit": options.Limit})
	}
	if options.Sort != nil {
		optionStage = append(optionStage, bson.M{"$sort": generateAggregateSortOptions(sortFields)})
	}

	return optionStage
}

func getAggregateColumnName(column string) string {
	return strings.Split(column, ":")[0]
}

func getAggregateAsColumnName(function, column string) string {
	format := "nested"
	arr := strings.Split(column, ":")
	if len(arr) == 2 && arr[1] == "table" {
		format = "table"
		column = arr[0]
	}

	return fmt.Sprintf("%s___%s___%s___%s", utils.GraphQLAggregate, format, function, strings.Join(strings.Split(column, "."), "__"))
}

func splitAggregateAsColumnName(asColumnName string) (format, functionName, columnName string, isAggregateColumn bool) {
	v := strings.Split(asColumnName, "___")
	if len(v) != 4 || !strings.HasPrefix(asColumnName, utils.GraphQLAggregate) {
		return "", "", "", false
	}
	return v[1], v[2], v[3], true
}

func getNestedObject(doc map[string]interface{}) {
	resultObj := make(map[string]interface{})
	for asColumnName, value := range doc {
		format, functionName, columnName, isAggregateColumn := splitAggregateAsColumnName(asColumnName)
		if isAggregateColumn {
			delete(doc, asColumnName)

			if format == "table" {
				doc[columnName] = value
				continue
			}

			funcValue, ok := resultObj[functionName]
			if !ok {
				// NOTE: This case occurs for count function with no column name (using * operator instead)
				if columnName == "" {
					resultObj[functionName] = value
				} else {
					resultObj[functionName] = map[string]interface{}{columnName: value}
				}
				continue
			}
			funcValue.(map[string]interface{})[columnName] = value
		}
		groupDoc, ok := doc["_id"]
		if ok {
			for key, value := range groupDoc.(map[string]interface{}) {
				doc[key] = value
			}
		}
	}
	if len(resultObj) > 0 {
		doc[utils.GraphQLAggregate] = resultObj
	}
	delete(doc, "_id")
}
