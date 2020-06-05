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
				args:           []interface{}{mock.Anything, "arithmetic", "adder", "", model.TokenClaims{}, map[string]interface{}{"num1": 10, "num2": 20}},
				paramsReturned: []interface{}{map[string]interface{}{"sum": 30}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsFuncCallAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.TokenClaims{}, nil},
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
		wantErr:    false,
		wantResult: map[string]interface{}{"adder": map[string]interface{}{"sum": 30}},
	},
}
