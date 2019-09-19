package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"errors"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) generateAllReq(field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	if len(field.Directives) > 0 {
		// Insert query function
		if strings.HasPrefix(field.Name.Value, "insert_") {
			col := strings.TrimPrefix(field.Name.Value, "insert_")
			result, err := graph.generateWriteReq(field, token, store)
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

			result, err := graph.genrateDeleteReq(field, token, store)
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

			result, err := graph.genrateUpdateReq(field, token, store)
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
			return map[string]interface{}{"status": 200}, graph.crud.Create(context.TODO(), dbType, graph.project, r.Col, &t)

		case string(utils.Delete):

			t := model.DeleteRequest{Operation: r.Operation, Find: r.Find}
			return map[string]interface{}{"status": 200}, graph.crud.Delete(context.TODO(), dbType, graph.project, r.Col, &t)

		case string(utils.Update):

			t := model.UpdateRequest{Operation: r.Operation, Find: r.Find, Update: r.Update}
			return map[string]interface{}{"status": 200}, graph.crud.Update(context.TODO(), dbType, graph.project, r.Col, &t)

		default:
			return map[string]interface{}{"error": "Wrong Operation"}, nil

		}

	}
	return map[string]interface{}{"status": 200}, graph.crud.Batch(context.TODO(), dbType, graph.project, req)
}

func (graph *Module) handleMutation(node ast.Node, token string, store utils.M, cb callback) {
	op := node.(*ast.OperationDefinition)
	fieldDBMapping := map[string]string{}

	reqs := map[string][]model.AllRequest{}
	queryResults := map[string]interface{}{}
	results := map[string]interface{}{}

	for _, v := range op.SelectionSet.Selections {

		field := v.(*ast.Field)

		dbType := GetDBType(field)

		fieldDBMapping[getFieldName(field)] = dbType

		r, ok := reqs[dbType]
		if !ok {
			r = []model.AllRequest{}
		}

		singleRequest, err := graph.generateAllReq(field, token, store)
		if err != nil {
			cb(nil, err)
			return
		}

		r = append(r, *singleRequest)
		reqs[dbType] = r
	}

	for dbType, reqs := range reqs {
		obj, err := graph.execAllReq(context.TODO(), dbType, graph.project, &model.BatchRequest{reqs})
		if err != nil {
			obj["error"] = err.Error()
			obj["status"] = 500
		}

		queryResults[dbType] = obj
	}

	for fieldName, dbType := range fieldDBMapping {
		results[fieldName] = queryResults[dbType]
	}

	cb(results, nil)
	return
}
