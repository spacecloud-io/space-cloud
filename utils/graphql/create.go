package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

<<<<<<< HEAD
func (graph *Module) execWriteRequest(field *ast.Field, token string, store m) (m, error) {
=======
func (graph *Module) execWriteRequest(field *ast.Field, store utils.M) (utils.M, error) {
>>>>>>> 9e6cacee503bece605f7e123f7ca4f25c1005c5b
	dbType := field.Directives[0].Name.Value
	col := strings.TrimPrefix(field.Name.Value, "insert_")

	req, err := generateCreateRequest(field, store)
	if err != nil {
		return nil, err
	}

	status, err := graph.auth.IsCreateOpAuthorised(graph.project, dbType, col, token, req)
	if err != nil {
		return nil, err
	}

	return utils.M{"status": status}, graph.crud.Create(context.TODO(), dbType, graph.project, col, req)
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
