package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execWriteRequest(field *ast.Field, token string, store m) (m, error) {
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

	return m{"status": status}, graph.crud.Create(context.TODO(), dbType, graph.project, col, req)
}

func generateCreateRequest(field *ast.Field, store m) (*model.CreateRequest, error) {
	// Create a create request object
	req := model.CreateRequest{Operation: utils.All}

	var err error
	req.Document, err = extractDocs(field.Arguments, store)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func extractDocs(args []*ast.Argument, store m) ([]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "docs":
			temp, err := parseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			return temp.([]interface{}), nil
		}
	}

	return []interface{}{}, nil
}
