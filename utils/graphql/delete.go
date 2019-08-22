package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execDeleteRequest(field *ast.Field, token string, store m) (m, error) {
	dbType := field.Directives[0].Name.Value
	col := strings.TrimPrefix(field.Name.Value, "delete_")

	req, err := generateDeleteRequest(field, store)
	if err != nil {
		return nil, err
	}

	status, err := graph.auth.IsDeleteOpAuthorised(graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}

	return m{"status": status}, graph.crud.Delete(context.TODO(), dbType, graph.project, col, req)
}

func generateDeleteRequest(field *ast.Field, store m) (*model.DeleteRequest, error) {
	var err error

	// Create a delete request object
	deleteRequest := model.DeleteRequest{Operation: utils.All}

	deleteRequest.Find, err = extractWhereClause(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &deleteRequest, nil
}
