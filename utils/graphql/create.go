package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execWriteRequest(ctx context.Context, field *ast.Field, token string, store utils.M) (map[string]interface{}, error) {
	dbType, err := GetDBType(field)
	if err != nil {
		return nil, err
	}

	col := strings.TrimPrefix(field.Name.Value, "insert_")

	req, err := generateCreateRequest(field, store)
	if err != nil {
		return nil, err
	}

	status, err := graph.auth.IsCreateOpAuthorised(ctx, graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"status": status}, graph.crud.Create(ctx, dbType, graph.project, col, req)
}

func (graph *Module) generateWriteReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	dbType, err := GetDBType(field)
	if err != nil {
		return nil, err
	}

	col := strings.TrimPrefix(field.Name.Value, "insert_")

	req, err := generateCreateRequest(field, store)
	if err != nil {
		return nil, err
	}

	_, err = graph.auth.IsCreateOpAuthorised(ctx, graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}
	return generateCreateAllRequest(req), nil
}

func generateCreateAllRequest(req *model.CreateRequest) *model.AllRequest {
	return &model.AllRequest{Operation: req.Operation, Document: req.Document}
}

func generateCreateRequest(field *ast.Field, store utils.M) (*model.CreateRequest, error) {
	// Create a create request object
	req := model.CreateRequest{Operation: utils.All}

	var err error
	req.Document, err = extractDocs(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func extractDocs(args []*ast.Argument, store utils.M) ([]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "docs":
			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			return temp.([]interface{}), nil
		}
	}

	return []interface{}{}, nil
}
