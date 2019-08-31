package graphql

import (
	"errors"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execFuncCall(field *ast.Field, store utils.M) (interface{}, error) {
	serviceName := getDBType(field)
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

func generateFuncCallRequest(field *ast.Field, store utils.M) (*model.FunctionsRequest, error) {
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

func getFuncTimeout(field *ast.Field, store utils.M) (int, error) {
	if len(field.Directives[0].Arguments) > 0 {
		for _, v := range field.Directives[0].Arguments {
			if v.Name.Value == "func" {
				val, err := ParseValue(v.Value, store)
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

func getFuncParams(field *ast.Field, store utils.M) (utils.M, error) {
	obj := make(utils.M, len(field.Arguments))

	for _, v := range field.Arguments {
		val, err := ParseValue(v.Value, store)
		if err != nil {
			return nil, err
		}

		obj[v.Name.Value] = val
	}

	return obj, nil
}
