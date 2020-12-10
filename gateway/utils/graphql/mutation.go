package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) generateAllReq(ctx context.Context, field *ast.Field, dbAlias, token string, store map[string]interface{}) (model.RequestParams, []*model.AllRequest, []interface{}, error) {
	if len(field.Directives) > 0 {
		// Insert query function
		if strings.HasPrefix(field.Name.Value, "insert_") {
			reqParams, result, returningDocs, err := graph.generateWriteReq(ctx, field, token, store)
			if err != nil {
				return model.RequestParams{}, nil, nil, err
			}
			return reqParams, result, returningDocs, nil
		}

		// Delete query function
		if strings.HasPrefix(field.Name.Value, "delete_") {
			col := strings.TrimPrefix(field.Name.Value, "delete_")

			reqParams, result, err := graph.genrateDeleteReq(ctx, field, token, store)
			if err != nil {
				return model.RequestParams{}, nil, nil, err
			}
			result.Type = string(model.Delete)
			result.Col = col
			result.DBAlias = dbAlias
			return reqParams, []*model.AllRequest{result}, nil, nil
		}

		// Update query function
		if strings.HasPrefix(field.Name.Value, "update_") {
			col := strings.TrimPrefix(field.Name.Value, "update_")

			reqParams, result, err := graph.generateUpdateReq(ctx, field, token, store)
			if err != nil {
				return reqParams, nil, nil, err
			}
			result.Type = string(model.Update)
			result.Col = col
			result.DBAlias = dbAlias
			return reqParams, []*model.AllRequest{result}, nil, nil

		}
	}
	return model.RequestParams{}, nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Target database not provided for field %s", getFieldName(field)), nil, nil)
}

func (graph *Module) execAllReq(ctx context.Context, dbAlias, project string, req *model.BatchRequest, params model.RequestParams) (map[string]interface{}, error) {
	if len(req.Requests) == 1 {
		r := req.Requests[0]
		switch r.Type {
		case string(model.Create):
			t := model.CreateRequest{Operation: r.Operation, Document: r.Document}
			return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Create(ctx, dbAlias, r.Col, &t, params)

		case string(model.Delete):

			t := model.DeleteRequest{Operation: r.Operation, Find: r.Find}
			return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Delete(ctx, dbAlias, r.Col, &t, params)

		case string(model.Update):

			t := model.UpdateRequest{Operation: r.Operation, Find: r.Find, Update: r.Update}
			return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Update(ctx, dbAlias, r.Col, &t, params)

		default:
			return map[string]interface{}{"error": "Wrong Operation"}, nil
		}
	}
	params.Resource = "db-batch"
	return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Batch(ctx, dbAlias, req, params)
}

func (graph *Module) handleMutation(ctx context.Context, node ast.Node, token string, store utils.M, cb model.GraphQLCallback) {
	op := node.(*ast.OperationDefinition)
	fieldDBMapping := map[string]string{}
	fieldReturningDocsMapping := map[string][]interface{}{}

	reqs := map[string][]*model.AllRequest{}
	queryResults := map[string]map[string]interface{}{}
	results := map[string]interface{}{}
	var reqParams model.RequestParams
	// A single mutation query can have same or different types of mutation
	// mutation {
	//		insert_...
	//		update_...
	//		delete_...
	// }
	// range over these different mutation queries
	for _, v := range op.SelectionSet.Selections {

		field := v.(*ast.Field)

		// for query insert_... @db {} -> dbAlias is "db"
		dbAlias, err := graph.GetDBAlias(ctx, field, token, store)
		if err != nil {
			cb(nil, err)
			return
		}

		// Generate a *model.AllRequest object for this given field
		params, generatedRequests, returningDocs, err := graph.generateAllReq(ctx, field, dbAlias, token, store)
		if err != nil {
			cb(nil, err)
			return
		}
		reqParams = params

		// Keep a record of which field maps to which db and which returning docs
		fieldDBMapping[getFieldName(field)] = dbAlias
		fieldReturningDocsMapping[getFieldName(field)] = returningDocs

		// Add the request to the number of requests available for that database
		for _, v := range generatedRequests {
			reqs[v.DBAlias] = append(reqs[v.DBAlias], v)
		}
	}

	for dbAlias, reqs := range reqs {
		obj, err := graph.execAllReq(ctx, dbAlias, graph.project, &model.BatchRequest{Requests: reqs}, reqParams)
		if err != nil {
			obj["error"] = err.Error()
			obj["status"] = 500
		}

		queryResults[dbAlias] = obj
	}

	for fieldName, dbAlias := range fieldDBMapping {
		result := queryResults[dbAlias]
		result["returning"] = fieldReturningDocsMapping[fieldName]
		results[fieldName] = result
	}
	// field, ok := op.SelectionSet.Selections[0].(*ast.Field)
	// if !ok {

	// }

	filteredResults := map[string]interface{}{}
	for _, selectionResult := range op.SelectionSet.Selections {
		v, _ := selectionResult.(*ast.Field)
		filteredResults[getFieldName(v)] = filterResults(v, results)
	}

	cb(filteredResults, nil)
}

func filterResults(field *ast.Field, results map[string]interface{}) map[string]interface{} {

	filteredResults := map[string]interface{}{}
	for _, resultValue := range results {

		v, ok := resultValue.(map[string]interface{})
		if !ok {
			continue
		}
		if field.SelectionSet == nil {
			return filteredResults
		}
		for _, returnFieldTemp := range field.SelectionSet.Selections {
			returnField := returnFieldTemp.(*ast.Field)
			returnFieldName := returnField.Name.Value
			value, ok := v[returnFieldName]
			if ok {
				if returnField.SelectionSet != nil {
					value = Filter(returnField, value)
				}
				filteredResults[returnFieldName] = value
			}
		}

		value, ok := v["__typename"]
		if ok {
			filteredResults["__typename"] = value
		}

	}

	return filteredResults
}

// Filter filers the result based on the provided selection set
func Filter(field *ast.Field, value interface{}) interface{} {
	switch val := value.(type) {
	case map[string]interface{}:
		newMap := map[string]interface{}{}
		for k, v := range val {
			for _, returnFieldTemp := range field.SelectionSet.Selections {
				returnField := returnFieldTemp.(*ast.Field)
				returnFieldName := returnField.Name.Value
				if k == returnFieldName {
					if returnField.SelectionSet != nil {
						v = Filter(returnField, v)
					}
					newMap[k] = v
				}
			}
		}
		return newMap
	case []interface{}:
		newArray := make([]interface{}, len(val))
		for i, v := range val {
			newArray[i] = Filter(field, v)
		}
		return newArray

	default:
		return nil
	}
}
