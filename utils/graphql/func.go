package graphql

import (
	"context"
	"errors"
	"time"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execFuncCall(ctx context.Context, token string, field *ast.Field, store utils.M, cb callback) {
	serviceName, err := GetDBType(field)
	if err != nil {
		cb(nil, err)
		return
	}

	funcName, err := getFuncName(field)
	if err != nil {
		cb(nil, err)
		return
	}

	timeout, err := getFuncTimeout(field, store)
	if err != nil {
		cb(nil, err)
		return
	}

	params, err := getFuncParams(field, store)
	if err != nil {
		cb(nil, err)
		return
	}

	if _, err := graph.auth.IsFuncCallAuthorised(ctx, graph.project, serviceName, funcName, token, params); err != nil {
		cb(nil, err)
		return
	}

	go func() {
		ctx2, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()

		result, err := graph.functions.CallWithContext(ctx2, serviceName, funcName, token, params)
		cb(result, err)
		return
	}()
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

func getFuncParams(field *ast.Field, store utils.M) (map[string]interface{}, error) {
	obj := make(map[string]interface{}, len(field.Arguments))

	for _, v := range field.Arguments {
		val, err := ParseValue(v.Value, store)
		if err != nil {
			return nil, err
		}

		obj[v.Name.Value] = val
	}

	return obj, nil
}
