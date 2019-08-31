package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execDeleteRequest(field *ast.Field, token string, store utils.M) (map[string]interface{}, error) {
	dbType := getDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "delete_")

	req, err := generateDeleteRequest(field, store)
	if err != nil {
		return nil, err
	}

	t := model.DeleteRequest{Operation: req.Operation, Find: req.Find}

	status, err := graph.auth.IsDeleteOpAuthorised(graph.project, dbType, col, token, &t)
	if err != nil {
		return nil, err
	}

	return utils.M{"status": status}, graph.crud.Delete(context.TODO(), dbType, graph.project, col, &t)
}

func (graph *Module) genrateDeleteReq(field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	dbType := field.Directives[0].Name.Value
	col := strings.TrimPrefix(field.Name.Value, "delete_")

	req, err := generateDeleteRequest(field, store)
	if err != nil {
		return nil, err
	}
	t := model.DeleteRequest{Operation: req.Operation, Find: req.Find}

	_, err = graph.auth.IsDeleteOpAuthorised(graph.project, dbType, col, token, &t)
	if err != nil {
		return nil, err
	}
	return req, nil

}

func generateDeleteRequest(field *ast.Field, store utils.M) (*model.AllRequest, error) {
	var err error

	// Create a delete request object
	deleteRequest := model.AllRequest{Operation: utils.All}

	deleteRequest.Find, err = extractWhereClause(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &deleteRequest, nil
}
