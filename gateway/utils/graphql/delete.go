package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) genrateDeleteReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) (model.RequestParams, *model.AllRequest, error) {
	dbAlias, err := graph.GetDBAlias(ctx, field)
	if err != nil {
		return model.RequestParams{}, nil, err
	}
	col := strings.TrimPrefix(field.Name.Value, "delete_")

	req, err := generateDeleteRequest(ctx, field, store)
	if err != nil {
		return model.RequestParams{}, nil, err
	}

	_, err = graph.auth.IsDeleteOpAuthorised(ctx, graph.project, dbAlias, col, token, req)
	if err != nil {
		return model.RequestParams{}, nil, err
	}
	return model.RequestParams{}, generateDeleteAllRequest(req), nil

}

func generateDeleteAllRequest(req *model.DeleteRequest) *model.AllRequest {
	return &model.AllRequest{Operation: req.Operation, Find: req.Find}
}

func generateDeleteRequest(ctx context.Context, field *ast.Field, store utils.M) (*model.DeleteRequest, error) {
	var err error

	// Create a delete request object
	deleteRequest := model.DeleteRequest{Operation: utils.All}

	deleteRequest.Find, err = ExtractWhereClause(ctx, field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &deleteRequest, nil
}
