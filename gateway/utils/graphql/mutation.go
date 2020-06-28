package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) generateAllReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) ([]*model.AllRequest, []interface{}, error) {
	if len(field.Directives) > 0 {
		// Insert query function
		if strings.HasPrefix(field.Name.Value, "insert_") {
			result, returningDocs, err := graph.generateWriteReq(ctx, field, token, store)
			if err != nil {
				return nil, nil, err
			}
			return result, returningDocs, nil
		}

		// Delete query function
		if strings.HasPrefix(field.Name.Value, "delete_") {
			col := strings.TrimPrefix(field.Name.Value, "delete_")

			result, err := graph.genrateDeleteReq(ctx, field, token, store)
			if err != nil {
				return nil, nil, err
			}
			result.Type = string(utils.Delete)
			result.Col = col
			return []*model.AllRequest{result}, nil, nil
		}

		// Update query function
		if strings.HasPrefix(field.Name.Value, "update_") {
			col := strings.TrimPrefix(field.Name.Value, "update_")

			result, err := graph.genrateUpdateReq(ctx, field, token, store)
			if err != nil {
				return nil, nil, err
			}
			result.Type = string(utils.Update)
			result.Col = col
			return []*model.AllRequest{result}, nil, nil

		}
	}
	return nil, nil, fmt.Errorf("target database not provided for field %s", getFieldName(field))
}

func (graph *Module) execAllReq(ctx context.Context, dbAlias, project string, req *model.BatchRequest) (map[string]interface{}, error) {
	if len(req.Requests) == 1 {
		r := req.Requests[0]
		switch r.Type {
		case string(utils.Create):
			t := model.CreateRequest{Operation: r.Operation, Document: r.Document}
			return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Create(ctx, dbAlias, r.Col, &t)

		case string(utils.Delete):

			t := model.DeleteRequest{Operation: r.Operation, Find: r.Find}
			return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Delete(ctx, dbAlias, r.Col, &t)

		case string(utils.Update):

			t := model.UpdateRequest{Operation: r.Operation, Find: r.Find, Update: r.Update}
			return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Update(ctx, dbAlias, r.Col, &t)

		default:
			return map[string]interface{}{"error": "Wrong Operation"}, nil

		}

	}
	return map[string]interface{}{"status": 200, "error": nil}, graph.crud.Batch(ctx, dbAlias, req)
}

func (graph *Module) handleMutation(ctx context.Context, node ast.Node, token string, store utils.M, cb model.GraphQLCallback) {
	op := node.(*ast.OperationDefinition)
	fieldDBMapping := map[string]string{}
	fieldReturningDocsMapping := map[string][]interface{}{}

	reqs := map[string][]*model.AllRequest{}
	queryResults := map[string]map[string]interface{}{}
	results := map[string]interface{}{}

	for _, v := range op.SelectionSet.Selections {

		field := v.(*ast.Field)

		dbAlias, err := graph.GetDBAlias(field)
		if err != nil {
			cb(nil, err)
			return
		}

		r, ok := reqs[dbAlias]
		if !ok {
			r = []*model.AllRequest{}
		}

		// Generate a *model.AllRequest object for this given field
		generatedRequests, returningDocs, err := graph.generateAllReq(ctx, field, token, store)
		if err != nil {
			cb(nil, err)
			return
		}

		// Keep a record of which field maps to which db and which returning docs
		fieldDBMapping[getFieldName(field)] = dbAlias
		fieldReturningDocsMapping[getFieldName(field)] = returningDocs

		// Add the request to the number of requests available for that database
		r = append(r, generatedRequests...)
		reqs[dbAlias] = r
	}

	for dbAlias, reqs := range reqs {
		obj, err := graph.execAllReq(ctx, dbAlias, graph.project, &model.BatchRequest{Requests: reqs})
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
