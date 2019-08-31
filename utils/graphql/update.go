package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execUpdateRequest(field *ast.Field, token string, store utils.M) (map[string]interface{}, error) {
	dbType := getDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "update_")
	req, err := generateUpdateRequest(field, store)
	if err != nil {
		return nil, err
	}

	t := model.UpdateRequest{Operation: req.Operation, Find: req.Find, Update: req.Update}

	status, err := graph.auth.IsUpdateOpAuthorised(graph.project, dbType, col, token, &t)
	if err != nil {
		return nil, err
	}

	return utils.M{"status": status}, graph.crud.Update(context.TODO(), dbType, graph.project, col, &t)
}

func (graph *Module) genrateUpdateReq(field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	dbType := getDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "update_")
	req, err := generateUpdateRequest(field, store)
	if err != nil {
		return nil, err
	}

	t := model.UpdateRequest{Operation: req.Operation, Find: req.Find, Update: req.Update}

	_, err = graph.auth.IsUpdateOpAuthorised(graph.project, dbType, col, token, &t)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func extractUpdateOperation(args []*ast.Argument, store utils.M) (string, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "op":
			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return "", err
			}
			if temp.(string) == "upsert" {
				return utils.Upsert, nil
			}

			return utils.All, nil
		}
	}
	return utils.All, nil
}

func generateUpdateRequest(field *ast.Field, store utils.M) (*model.AllRequest, error) {
	var err error
	var updateRequest model.AllRequest

	updateRequest.Operation, err = extractUpdateOperation(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	updateRequest.Find, err = extractWhereClause(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	updateRequest.Update, err = extractUpdateArgs(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &updateRequest, nil
}

func extractUpdateArgs(args []*ast.Argument, store utils.M) (utils.M, error) {
	t := map[string]interface{}{}
	for _, v := range args {
		switch v.Name.Value {
		case "set", "inc", "mul", "max", "min", "currentTimestamp", "currentDate", "push", "rename", "remove":
			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			t["$"+v.Name.Value] = temp
		}
	}
	return t, nil
}
