package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"errors"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) generateAllReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	if len(field.Directives) > 0 {
		// Insert query function
		if strings.HasPrefix(field.Name.Value, "insert_") {
			col := strings.TrimPrefix(field.Name.Value, "insert_")
			result, err := graph.generateWriteReq(ctx, field, token, store)
			if err != nil {
				return nil, err
			}
			result.Type = string(utils.Create)
			result.Col = col
			return result, nil
		}

		// Delete query function
		if strings.HasPrefix(field.Name.Value, "delete_") {
			col := strings.TrimPrefix(field.Name.Value, "delete_")

			result, err := graph.genrateDeleteReq(ctx, field, token, store)
			if err != nil {
				return nil, err
			}
			result.Type = string(utils.Delete)
			result.Col = col
			return result, nil
		}

		// Update query function
		if strings.HasPrefix(field.Name.Value, "update_") {
			col := strings.TrimPrefix(field.Name.Value, "update_")

			result, err := graph.genrateUpdateReq(ctx, field, token, store)
			if err != nil {
				return nil, err
			}
			result.Type = string(utils.Update)
			result.Col = col
			return result, nil

		}
	}
	return nil, errors.New("No directive present")
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

	reqs := map[string][]model.AllRequest{}
	queryResults := map[string]interface{}{}
	results := map[string]interface{}{}

	for _, v := range op.SelectionSet.Selections {

		field := v.(*ast.Field)

		dbType, err := GetDBType(field)
		if err != nil {
			cb(nil, err)
			return
		}

		fieldDBMapping[getFieldName(field)] = dbType

		r, ok := reqs[dbType]
		if !ok {
			r = []model.AllRequest{}
		}

		singleRequest, err := graph.generateAllReq(ctx, field, token, store)
		if err != nil {
			cb(nil, err)
			return
		}

		r = append(r, *singleRequest)
		reqs[dbType] = r
	}

	for dbType, reqs := range reqs {
		obj, err := graph.execAllReq(ctx, dbType, graph.project, &model.BatchRequest{reqs})
		if err != nil {
			obj["error"] = err.Error()
			obj["status"] = 500
		}

		queryResults[dbType] = obj
	}

	for fieldName, dbType := range fieldDBMapping {
		results[fieldName] = queryResults[dbType]
	}
	// field, ok := op.SelectionSet.Selections[0].(*ast.Field)
	// if !ok {

	// }

	filteredResults := map[string]interface{}{}
	for _, selectionResult := range op.SelectionSet.Selections {
		v, ok := selectionResult.(*ast.Field)
		if !ok {

		}
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

		for _, returnField := range field.SelectionSet.Selections {
			returnFieldName := returnField.(*ast.Field).Name.Value
			value, ok := v[returnFieldName]
			if ok {
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
