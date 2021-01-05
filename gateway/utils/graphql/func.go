package graphql

import (
	"context"
	"errors"
	"time"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
	authHelpers "github.com/spaceuptech/space-cloud/gateway/modules/auth/helpers"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) execFuncCall(ctx context.Context, token string, field *ast.Field, store utils.M, cb model.GraphQLCallback) {
	serviceName, _ := graph.getDirectiveName(ctx, field.Directives[0], token, store)

	funcName, err := getFuncName(field)
	if err != nil {
		cb(nil, err)
		return
	}

	timeout, err := getFuncTimeout(ctx, field, store)
	if err != nil {
		cb(nil, err)
		return
	}

	params, err := getFuncParams(ctx, field, store)
	if err != nil {
		cb(nil, err)
		return
	}

	cacheConfig, err := generateCacheOptions(ctx, field.Directives, store)
	if err != nil {
		cb(nil, err)
		return
	}

	actions, reqParams, err := graph.auth.IsFuncCallAuthorised(ctx, graph.project, serviceName, funcName, token, params)
	if err != nil {
		cb(nil, err)
		return
	}
	// Note: there is some inconsistency between REST & GraphQL of remote services, this does not affect the core logic
	// The reqParams object contains only token information such as claims,
	// but in REST api, we use this function to extract more info from request,< reqParams = utils.ExtractRequestParams(r, reqParams, req) >
	// which is not possible in graphql, we can only set the body of req params object
	reqParams.Payload = params

	go func() {
		var ctx2 = ctx
		if timeout != 0 {
			c, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()
			ctx2 = c
		}

		_, result, err := graph.functions.CallWithContext(ctx2, serviceName, funcName, token, reqParams, &model.FunctionsRequest{Params: params, Timeout: timeout, Cache: cacheConfig})
		_ = authHelpers.PostProcessMethod(ctx, graph.aesKey, actions, result)
		cb(result, err)
	}()
}

// func generateFuncCallRequest(field *ast.Field, store utils.M) (*model.FunctionsRequest, error) {
// 	timeout, err := getFuncTimeout(field, store)
// 	if err != nil {
// 		return nil, err
// 	}

// 	params, err := getFuncParams(field, store)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.FunctionsRequest{Params: params, Timeout: timeout}, nil
// }

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

func getFuncTimeout(ctx context.Context, field *ast.Field, store utils.M) (int, error) {
	if len(field.Directives[0].Arguments) > 0 {
		for _, v := range field.Directives[0].Arguments {
			if v.Name.Value == "timeout" {
				val, err := utils.ParseGraphqlValue(v.Value, store)
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
	return 0, nil
}

func getFuncParams(ctx context.Context, field *ast.Field, store utils.M) (map[string]interface{}, error) {
	obj := make(map[string]interface{}, len(field.Arguments))

	for _, v := range field.Arguments {
		val, err := utils.ParseGraphqlValue(v.Value, store)
		if err != nil {
			return nil, err
		}

		obj[v.Name.Value] = val
	}

	return obj, nil
}
