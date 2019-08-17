package graphql

import (
	"errors"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
)

func (graph *Module) execFuncCall(field *ast.Field, store m) (interface{}, error) {
	serviceName := field.Directives[0].Name.Value
	funcName, err := getFuncName(field)
	if err != nil {
		return nil, err
	}

	timeout, err := getFuncTimeout(field, store)
	if err != nil {
		return nil, err
	}

	params, err := getFuncParams(field, store)
	if err != nil {
		return nil, err
	}

	claims, err := graph.auth.IsFuncCallAuthorised(graph.project, serviceName, funcName, "", params)
	if err != nil {
		return nil, err
	}

	return graph.functions.Call(serviceName, funcName, claims, params, timeout)
}

func generateFuncCallRequest(field *ast.Field, store m) (*model.FunctionsRequest, error) {
	timeout, err := getFuncTimeout(field, store)
	if err != nil {
		return nil, err
	}

	params, err := getFuncParams(field, store)
	if err != nil {
		return nil, err
	}

	return &model.FunctionsRequest{Params: params, Timeout: timeout}, nil
}

func getFuncName(field *ast.Field) (string, error) {
	if len(field.Directives[0].Arguments) > 0 {
		for _, v := range field.Directives[0].Arguments {
			if v.Name.Value == "func" {
				col, ok := v.Value.GetValue().(string)
				if !ok {
					return "", errors.New("Invalid value for collection: " + string(v.Value.GetLoc().Source.Body)[v.Value.GetLoc().Start:v.Value.GetLoc().End])
				}
				return col, nil
			}
		}
	}
	return field.Name.Value, nil
}

func getFuncTimeout(field *ast.Field, store m) (int, error) {
	if len(field.Directives[0].Arguments) > 0 {
		for _, v := range field.Directives[0].Arguments {
			if v.Name.Value == "func" {
				val, err := parseValue(v.Value, store)
				if err != nil {
					return 0, err
				}

				timeout, ok := val.(int)
				if !ok {
					return 0, errors.New("Invalid value for collection: " + string(v.Value.GetLoc().Source.Body)[v.Value.GetLoc().Start:v.Value.GetLoc().End])
				}
				return timeout, nil
			}
		}
	}
	return 5, nil
}

func getFuncParams(field *ast.Field, store m) (m, error) {
	obj := make(m, len(field.Arguments))

	for _, v := range field.Arguments {
		val, err := parseValue(v.Value, store)
		if err != nil {
			return nil, err
		}

		obj[v.Name.Value] = val
	}

	return obj, nil
}
