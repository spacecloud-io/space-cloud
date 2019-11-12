package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execUpdateRequest(ctx context.Context, field *ast.Field, token string, store utils.M) (map[string]interface{}, error) {
	dbType, err := GetDBType(field)
	if err != nil {
		return nil, err
	}

	col := strings.TrimPrefix(field.Name.Value, "update_")
	req, err := generateUpdateRequest(field, store)
	if err != nil {
		return nil, err
	}

	status, err := graph.auth.IsUpdateOpAuthorised(ctx, graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}

	return utils.M{"status": status}, graph.crud.Update(ctx, dbType, graph.project, col, req)
}

func (graph *Module) genrateUpdateReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	dbType, err := GetDBType(field)
	if err != nil {
		return nil, err
	}
	col := strings.TrimPrefix(field.Name.Value, "update_")
	req, err := generateUpdateRequest(field, store)
	if err != nil {
		return nil, err
	}

	_, err = graph.auth.IsUpdateOpAuthorised(ctx, graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}
	return generateUpdateAllRequest(req), nil
}

func generateUpdateAllRequest(req *model.UpdateRequest) *model.AllRequest {
	return &model.AllRequest{Operation: req.Operation, Find: req.Find, Update: req.Update}
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

func generateUpdateRequest(field *ast.Field, store utils.M) (*model.UpdateRequest, error) {
	var err error
	var updateRequest model.UpdateRequest

	updateRequest.Operation, err = extractUpdateOperation(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	updateRequest.Find, err = ExtractWhereClause(field.Arguments, store)
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
