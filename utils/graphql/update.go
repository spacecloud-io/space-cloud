package graphql

import (
	"context"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execUpdateRequest(field *ast.Field, store m) (m, error) {
	dbType := field.Directives[0].Name.Value
	col := strings.TrimPrefix(field.Name.Value, "update_")

	req, err := generateUpdateRequest(field, store)
	if err != nil {
		return nil, err
	}
	status, err := graph.auth.IsUpdateOpAuthorised(graph.project, dbType, col, "", req)
	if err != nil {
		return nil, err
	}

	return m{"status": status}, graph.crud.Update(context.TODO(), dbType, graph.project, col, req)
}

func generateUpdateRequest(field *ast.Field, store m) (*model.UpdateRequest, error) {
	var err error
	var updateRequest model.UpdateRequest

	//---	 upsert identification is wrong
	if strings.HasPrefix(field.Name.Value, "upsert") {
		updateRequest = model.UpdateRequest{Operation: utils.Upsert}
	} else {
		updateRequest = model.UpdateRequest{Operation: utils.All}
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

func extractWhereClause(args []*ast.Argument, store m) (m, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "where":
			temp, err := parseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			return temp.(m), nil
		}
	}
	return m{}, nil
}

func extractUpdateArgs(args []*ast.Argument, store m) (m, error) {
	var t map[string]interface{}
	for _, v := range args {
		switch v.Name.Value {
		case "set", "inc", "mul", "max", "min", "currentTimestamp", "currentDate":
			temp, err := parseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			t["$"+v.Name.Value] = temp
		}
	}
	return t, nil
}
