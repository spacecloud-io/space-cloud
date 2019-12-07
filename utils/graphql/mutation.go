package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) generateAllReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) ([]model.AllRequest, []interface{}, error) {
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
			return []model.AllRequest{*result}, nil, nil
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
			return []model.AllRequest{*result}, nil, nil

		}
	}
	return nil, nil, fmt.Errorf("target database not provided for field %s", getFieldName(field))
}

func (graph *Module) execAllReq(ctx context.Context, dbType, project string, req *model.BatchRequest) (map[string]interface{}, error) {
	if len(req.Requests) == 1 {
		r := req.Requests[0]
		switch r.Type {
		case string(utils.Create):
			t := model.CreateRequest{Operation: r.Operation, Document: r.Document}
			return map[string]interface{}{"status": 200}, graph.crud.Create(ctx, dbType, graph.project, r.Col, &t)

		case string(utils.Delete):

			t := model.DeleteRequest{Operation: r.Operation, Find: r.Find}
			return map[string]interface{}{"status": 200}, graph.crud.Delete(ctx, dbType, graph.project, r.Col, &t)

		case string(utils.Update):

			t := model.UpdateRequest{Operation: r.Operation, Find: r.Find, Update: r.Update}
			return map[string]interface{}{"status": 200}, graph.crud.Update(ctx, dbType, graph.project, r.Col, &t)

		default:
			return map[string]interface{}{"error": "Wrong Operation"}, nil

		}

	}
	return map[string]interface{}{"status": 200}, graph.crud.Batch(ctx, dbType, graph.project, req)
}

func (graph *Module) handleMutation(ctx context.Context, node ast.Node, token string, store utils.M, cb callback) {
	op := node.(*ast.OperationDefinition)
	fieldDBMapping := map[string]string{}
	fieldReturningDocsMapping := map[string][]interface{}{}

	reqs := map[string][]model.AllRequest{}
	queryResults := map[string]map[string]interface{}{}
	results := map[string]interface{}{}

	for _, v := range op.SelectionSet.Selections {

		field := v.(*ast.Field)

		dbType, err := GetDBType(field)
		if err != nil {
			cb(nil, err)
			return
		}

		r, ok := reqs[dbType]
		if !ok {
			r = []model.AllRequest{}
		}

		// Generate a *model.AllRequest object for this given field
		generatedRequests, returningDocs, err := graph.generateAllReq(ctx, field, token, store)
		if err != nil {
			cb(nil, err)
			return
		}

		// Keep a record of which field maps to which db and which returning docs
		fieldDBMapping[getFieldName(field)] = dbType
		fieldReturningDocsMapping[getFieldName(field)] = returningDocs

		// Add the request to the number of requests available for that database
		r = append(r, generatedRequests...)
		reqs[dbType] = r
	}

	for dbType, reqs := range reqs {
		obj, err := graph.execAllReq(ctx, dbType, graph.project, &model.BatchRequest{Requests: reqs})
		if err != nil {
			obj["error"] = err.Error()
			obj["status"] = 500
		}

		queryResults[dbType] = obj
	}

	for fieldName, dbType := range fieldDBMapping {
		result := queryResults[dbType]
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
	return
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
					value = filter(returnField, value)
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

func filter(field *ast.Field, value interface{}) interface{} {
	switch val := value.(type) {
	case map[string]interface{}:
		newMap := map[string]interface{}{}
		for k, v := range val {
			for _, returnFieldTemp := range field.SelectionSet.Selections {
				returnField := returnFieldTemp.(*ast.Field)
				returnFieldName := returnField.Name.Value
				if k == returnFieldName {
					if returnField.SelectionSet != nil {
						v = filter(returnField, v)
					}
					newMap[k] = v
				}
			}
		}
		return newMap
	case []interface{}:
		newArray := make([]interface{}, len(val))
		for i, v := range val {
			newArray[i] = filter(field, v)
		}
		return newArray

	default:
		return nil
	}
}
