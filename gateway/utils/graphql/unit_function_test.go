package graphql_test

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

var functionTestCases = []tests{
	{
		name: "Function: Querying static endpoints",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"arithmetic"},
				paramsReturned: []interface{}{"", errors.New("invalid db alias provided")},
			},
		},
		functionMockArgs: []mockArgs{
			{
				method:         "CallWithContext",
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", mock.Anything, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
			{
				method:         "PostProcessMethod",
				args:           []interface{}{mock.Anything, mock.Anything},
				paramsReturned: []interface{}{nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								adder(
									num1 : 10,
									num2 : 20,
								) @arithmetic(timeout:10,func : "adder") {
									sum
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"adder": map[string]interface{}{"sum": 30}},
	},
	{
		name: "Function: Querying static endpoints error invalid type for directive arguments",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"arithmetic"},
				paramsReturned: []interface{}{"", errors.New("invalid db alias provided")},
			},
		},
		functionMockArgs: []mockArgs{
			{
				method:         "CallWithContext",
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", mock.Anything, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
			{
				method:         "PostProcessMethod",
				args:           []interface{}{mock.Anything, mock.Anything},
				paramsReturned: []interface{}{nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								adder(
									num1 : 10,
									num2 : 20,
								) @arithmetic(timeout:10,func : $data) {
									sum
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    true,
		wantResult: nil,
	},
	{
		name: "Function: Querying static endpoints error invalid type for directive timeout arguments ",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"arithmetic"},
				paramsReturned: []interface{}{"", errors.New("invalid db alias provided")},
			},
		},
		functionMockArgs: []mockArgs{
			{
				method:         "CallWithContext",
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", mock.Anything, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
			{
				method:         "PostProcessMethod",
				args:           []interface{}{mock.Anything, mock.Anything},
				paramsReturned: []interface{}{nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								adder(
									num1 : 10,
									num2 : 20,
								) @arithmetic(timeout: $data ,func : "adder") {
									sum
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    true,
		wantResult: nil,
	},

	{
		name: "Function: Querying static endpoints error function call not authorized",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"arithmetic"},
				paramsReturned: []interface{}{"", errors.New("invalid db alias provided")},
			},
		},
		functionMockArgs: []mockArgs{
			{
				method:         "CallWithContext",
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", mock.Anything, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("function call not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								adder(
									num1 : 10,
									num2 : 20,
								) @arithmetic {
									sum
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    true,
		wantResult: nil,
	},
	{
		name: "Function: Querying static endpoints error invalid function params",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"arithmetic"},
				paramsReturned: []interface{}{"", errors.New("invalid db alias provided")},
			},
		},
		functionMockArgs: []mockArgs{
			{
				method:         "CallWithContext",
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", mock.Anything, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("function call not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								adder(
									num1 : $data,
									num2 : 20,
								) @arithmetic {
									sum
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    true,
		wantResult: nil,
	},
	{
		name: "Function: Querying static endpoints error invalid function params",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"arithmetic"},
				paramsReturned: []interface{}{"", errors.New("invalid db alias provided")},
			},
		},
		functionMockArgs: []mockArgs{
			{
				method:         "CallWithContext",
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", mock.Anything, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("function call not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								adder(
									num1 : 10,
									num2 : 20,
								) @arithmetic {
									sum
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    true,
		wantResult: nil,
	},
}
