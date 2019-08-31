package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execWriteRequest(field *ast.Field, token string, store utils.M) (map[string]interface{}, error) {
	dbType := getDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "insert_")

	req, err := generateCreateRequest(field, store)
	if err != nil {
		return nil, err
	}

	t := model.CreateRequest{Operation: req.Operation, Document: req.Document}

	status, err := graph.auth.IsCreateOpAuthorised(graph.project, dbType, col, token, &t)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"status": status}, graph.crud.Create(context.TODO(), dbType, graph.project, col, &t)
}

func (graph *Module) generateWriteReq(field *ast.Field, token string, store map[string]interface{}) (*model.AllRequest, error) {
	dbType := getDBType(field)
	col := strings.TrimPrefix(field.Name.Value, "insert_")

	req, err := generateCreateRequest(field, store)
	if err != nil {
		return nil, err
	}

	t := model.CreateRequest{Operation: req.Operation, Document: req.Document}

	_, err = graph.auth.IsCreateOpAuthorised(graph.project, dbType, col, token, &t)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func generateCreateRequest(field *ast.Field, store utils.M) (*model.AllRequest, error) {
	// Create a create request object
	req := model.AllRequest{Operation: utils.All}

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
