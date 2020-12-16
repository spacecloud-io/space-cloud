package graphql_test

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

var prepareQueryTestCases = []tests{
	{
		name: "Prepared query : Simple query",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"custom_sql"},
				paramsReturned: []interface{}{"custom_sql", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"custom_sql", "insert1"},
				paramsReturned: []interface{}{true},
			},
			{
				method:         "ExecPreparedQuery",
				args:           []interface{}{mock.Anything, "custom_sql", "insert1", &model.PreparedQueryRequest{Params: map[string]interface{}{"id": "1", "name": "ash"}}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{}, map[string]interface{}{}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsPreparedQueryAuthorised",
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
								insert1(
									id : "1",
									name : "ash"
								) @custom_sql {
									status
									error
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert1": []interface{}{}},
	},
	{
		name: "Prepared query : Simple query error query not authorized",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"custom_sql"},
				paramsReturned: []interface{}{"custom_sql", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"custom_sql", "insert1"},
				paramsReturned: []interface{}{true},
			},
			{
				method:         "ExecPreparedQuery",
				args:           []interface{}{mock.Anything, "custom_sql", "insert1", &model.PreparedQueryRequest{Params: map[string]interface{}{"id": "1", "name": "ash"}}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{}, map[string]interface{}{}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsPreparedQueryAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("query not authorized")},
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
								insert1(
									id : "1",
									name : "ash"
								) @custom_sql {
									status
									error
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
		name: "Prepared query : Simple query error incorrect arguments",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"custom_sql"},
				paramsReturned: []interface{}{"custom_sql", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"custom_sql", "insert1"},
				paramsReturned: []interface{}{true},
			},
			{
				method:         "ExecPreparedQuery",
				args:           []interface{}{mock.Anything, "custom_sql", "insert1", &model.PreparedQueryRequest{Params: map[string]interface{}{"id": "1", "name": "ash"}}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{}, map[string]interface{}{}, nil},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsPreparedQueryAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("query not authorized")},
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
								insert1(
									id : "1",
									name : $data
								) @custom_sql {
									status
									error
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
