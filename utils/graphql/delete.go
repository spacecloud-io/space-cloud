package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execDeleteRequest(field *ast.Field, token string, store utils.M) (map[string]interface{}, error) {
	dbType := GetDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "delete_")

	req, err := generateDeleteRequest(field, store)
	if err != nil {
		return nil, err
	}

	status, err := graph.auth.IsDeleteOpAuthorised(graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}

	return utils.M{"status": status}, graph.crud.Delete(context.TODO(), dbType, graph.project, col, req)
}

func (graph *Module) genrateDeleteReq(field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	dbType := GetDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "delete_")

	req, err := generateDeleteRequest(field, store)
	if err != nil {
		return nil, err
	}

	_, err = graph.auth.IsDeleteOpAuthorised(graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}
	return generateDeleteAllRequest(req), nil

}

func generateDeleteAllRequest(req *model.DeleteRequest) *model.AllRequest {
	return &model.AllRequest{Operation: req.Operation, Find: req.Find}
}

func generateDeleteRequest(field *ast.Field, store utils.M) (*model.DeleteRequest, error) {
	var err error

	// Create a delete request object
	deleteRequest := model.DeleteRequest{Operation: utils.All}

	deleteRequest.Find, err = ExtractWhereClause(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &deleteRequest, nil
}
