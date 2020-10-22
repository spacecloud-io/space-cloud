package graphql_test

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

var distinct = "type"
var number int64 = 5

var queryTestCases = []tests{
	{
		name: "Query: Simple Query with templated directive",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db_t1"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db_t1", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db_t1"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db_t1", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db_t1", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "ParseToken",
				args:           []interface{}{mock.Anything, mock.Anything},
				paramsReturned: []interface{}{map[string]interface{}{"tenant": "t1"}, nil},
			},
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons @template(value: "db_{{.auth.tenant}}") {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}},
	},
	{
		name:           "Query: Simple Query with invalid templated directive",
		crudMockArgs:   []mockArgs{},
		schemaMockArgs: []mockArgs{},
		authMockArgs:   []mockArgs{},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons @template {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr: true,
	},
	{
		name: "Query: Simple Query with templated directive in variables",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db_t1"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db_t1", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db_t1"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db_t1", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db_t1", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "ParseToken",
				args:           []interface{}{mock.Anything, mock.Anything},
				paramsReturned: []interface{}{map[string]interface{}{"role": "t1"}, nil},
			},
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons @template(value: $tmpl) {
									id
									name
									power_level
								}
							}`,
				Variables: map[string]interface{}{"tmpl": "db_{{.auth.role}}"},
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}},
	},
	{
		name: "Query: Simple Query",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}},
	},
	{
		name: "Query: Simple Query error read request not authorized",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("request not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons @db {
									id
									name
									power_level
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
		name: "Query: Simple Query error get collection is incorrect",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("request not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons @db(col : $data) {
									id
									name
									power_level
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
		name: "Query: Simple Query error read error incorrect where clause provided",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "bulbasaur", "power_level": 60}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, errors.New("request not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : $data
									) @db {
									id
									name
									power_level
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
		name: "Query: Using where clause with equality operator (skipping _eq in where)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": 100},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : 100
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}},
	},

	{
		name: "Query: Using where clause with equality operator (skipping _eq in where)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": 100},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : 100
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}},
	},
	{
		name: "Query: aggregation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$eq": 100,
						}},
					Aggregate:   map[string][]string{"sum": {"power_level"}},
					GroupBy:     []interface{}{"power_level"},
					Operation:   utils.All,
					Options:     &model.ReadOptions{},
					IsBatch:     false,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"aggregate": map[string]interface{}{"sum": map[string]interface{}{"power_level": 100}}}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									group : [power_level],
									where : {
										power_level : {
											_eq : 100
										}
									}
								) @db {
									aggregate {
										sum {
											power_level
										}
									}
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"aggregate": map[string]interface{}{"sum": map[string]interface{}{"power_level": 100}}}}},
	},
	{
		name: "Query: aggregation invalid type provided for group by ",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$eq": 100,
						}},
					Aggregate:   map[string][]string{"sum": {"power_level"}},
					GroupBy:     []interface{}{"power_level"},
					Operation:   utils.All,
					Options:     &model.ReadOptions{},
					IsBatch:     false,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"aggregate": map[string]interface{}{"sum": map[string]interface{}{"power_level": 100}}}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									group : $data,
									where : {
										power_level : {
											_eq : 100
										}
									}
								) @db {
									aggregate {
										sum {
											power_level
										}
									}
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
		name: "Query: Using where clause with not equality operator (_ne)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$ne": 100,
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "4", "name": "snorlax", "power_level": 30}, map[string]interface{}{"id": "5", "name": "jigglypuff", "power_level": 40}, map[string]interface{}{"id": "5", "name": "squirtle", "power_level": 50}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_ne : 100
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "4", "name": "snorlax", "power_level": 30}, map[string]interface{}{"id": "5", "name": "jigglypuff", "power_level": 40}, map[string]interface{}{"id": "5", "name": "squirtle", "power_level": 50}}},
	},
	{
		name: "Query: Using where clause with comparision operator greater than (_gt)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$gt": 50,
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_gt : 50
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}},
	},
	{
		name: "Query: Using where clause with comparision operator greater than equal to (_gte)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$gte": 50,
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}, map[string]interface{}{"id": "2", "name": "ditto", "power_level": 50}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_gte : 50
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}, map[string]interface{}{"id": "2", "name": "ditto", "power_level": 50}}},
	},
	{
		name: "Query: Using where clause with comparision operator less than (_lt)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$lt": 50,
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "4", "name": "snorlax", "power_level": 30}, map[string]interface{}{"id": "5", "name": "jigglypuff", "power_level": 40}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_lt : 50
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "4", "name": "snorlax", "power_level": 30}, map[string]interface{}{"id": "5", "name": "jigglypuff", "power_level": 40}}},
	},
	{
		name: "Query: Using where clause with comparision operator less than equal to (_lte)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$lte": 50,
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "4", "name": "snorlax", "power_level": 30}, map[string]interface{}{"id": "5", "name": "jigglypuff", "power_level": 40}, map[string]interface{}{"id": "5", "name": "squirtle", "power_level": 50}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_lte : 50
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "4", "name": "snorlax", "power_level": 30}, map[string]interface{}{"id": "5", "name": "jigglypuff", "power_level": 40}, map[string]interface{}{"id": "5", "name": "squirtle", "power_level": 50}}},
	},
	{
		name: "Query: Using where clause with search operator (_regex)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"description": map[string]interface{}{
							"$regex": "(?i)strong",
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										description : {
											_regex: "(?i)strong"
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}},
	},
	{
		name: "Query: Using where clause with list based operator (_in)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$in": []interface{}{100, 50},
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}, map[string]interface{}{"id": "5", "name": "squirtle", "power_level": 50}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_in: [100, 50]
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}, map[string]interface{}{"id": "5", "name": "squirtle", "power_level": 50}}},
	},
	{
		name: "Query: Using where clause with list based operator not in (_nin)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"power_level": map[string]interface{}{
							"$nin": []interface{}{50},
						}},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										power_level : {
											_nin: [50]
										}
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}},
	},
	// {
	// 	name: "Query: Using where clause with JSON operator contains (_contains)",
	// 	crudMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method:         "IsPreparedQueryPresent",
	// 			args:           []interface{}{"db", "pokemons"},
	// 			paramsReturned: []interface{}{false},
	// 		},
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method: "Read",
	// 			args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
	//			Find: map[string]interface{}{
	// 					"stats": map[string]interface{}{
	// 						"$contains": "$jsonFilter",
	// 					}},
	// 				Aggregate: map[string][]string{},
	// 				GroupBy:   []interface{}{},
	// 				Operation: utils.All,
	// 				Options: &model.ReadOptions{
	// 					Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
	// 				},
	// 				IsBatch: true,
	// 			}},
	// 			paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}, nil},
	// 		},
	// 	},
	// 	schemaMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetSchema",
	// 			args:           []interface{}{"db", "pokemons"},
	// 			paramsReturned: []interface{}{model.Fields{}, true},
	// 		},
	// 	},
	// 	authMockArgs: []mockArgs{
	// 		{
	// 			method:         "IsReadOpAuthorised",
	// 			args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
	// 		},
	// 		{
	// 			method:         "PostProcessMethod",
	// 			args:           []interface{}{mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{nil},
	// 		},
	// 	},
	// 	args: args{
	// 		req: &model.GraphQLRequest{
	// 			OperationName: "query",
	// 			Query: `query {
	// 						pokemons(
	// 							where : {
	// 								stats : {
	// 									_contains: $jsonFilter
	// 								}
	// 							}
	// 						) @db {
	// 							id
	// 							name
	// 							power_level
	// 						}
	// 					}`,
	// 			Variables: map[string]interface{}{
	// 				"$jsonFilter": map[string]interface{}{
	// 					"combat_power": 500,
	// 				},
	// 			},
	// 		},
	// 		token: "",
	// 	},
	// 	wantErr:    false,
	// 	wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}, map[string]interface{}{"id": "2", "name": "charmander", "power_level": 100}}},
	// },
	{
		name: "Query: Multiple filters in where clause with and(default) operator",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"caught_on": map[string]interface{}{
							"$gte": "2019-06-01",
							"$lte": "2019-09-15",
						},
						"id": 1,
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										caught_on : {
											_gte: "2019-06-01",
											_lte: "2019-09-15"
										},
										id : 1
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}},
	},
	// TODO: check the query
	{
		name: "Query: Multiple filters in where clause with and operator (_and) explicit",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"$and": []interface{}{
							map[string]interface{}{
								"caught_on": map[string]interface{}{
									"$gte": "2019-06-01",
									"$lte": "2019-09-15",
								}},
							map[string]interface{}{
								"id": 1,
							}},
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										_and: [
										{
											caught_on : {
												_gte: "2019-06-01",
												_lte: "2019-09-15"
											}
										},
										{
											id : 1
										}
										]
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}},
	},
	{
		name: "Query: Multiple filters in where clause with or operator (_or)",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"$or": []interface{}{
							map[string]interface{}{"type": "fire"},
							map[string]interface{}{"is_legendary": true}},
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									where : {
										_or: [
											{type: "fire"},
											{is_legendary: true}
										]
									}
								) @db {
									id
									name
									power_level
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}},
	},
	// {
	// 	name: "Query: Filter with nested queries",
	// 	crudMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method:         "IsPreparedQueryPresent",
	// 			args:           []interface{}{"db", "pokemons"},
	// 			paramsReturned: []interface{}{false},
	// 		},
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method: "Read",
	// 			args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
	//			Find: map[string]interface{}{
	// 					"$or": []interface{}{
	// 						map[string]interface{}{"type": "fire"},
	// 						map[string]interface{}{"is_legendary": true}},
	// 				},
	// 				Aggregate: map[string][]string{},
	// 				GroupBy:   []interface{}{},
	// 				Operation: utils.All,
	// 				Options: &model.ReadOptions{
	// 					Select: map[string]int32{"id": 1, "name": 1, "power_level": 1},
	// 				},
	// 				IsBatch: true,
	// 			}},
	// 			paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}, nil},
	// 		},
	// 	},
	// 	schemaMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetSchema",
	// 			args:           []interface{}{"db", "pokemons"},
	// 			paramsReturned: []interface{}{model.Fields{}, true},
	// 		},
	// 	},
	// 	authMockArgs: []mockArgs{
	// 		{
	// 			method:         "IsReadOpAuthorised",
	// 			args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
	// 		},
	// 		{
	// 			method:         "PostProcessMethod",
	// 			args:           []interface{}{mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{nil},
	// 		},
	// 	},
	// 	args: args{
	// 		req: &model.GraphQLRequest{
	// 			OperationName: "query",
	// 			Query: `query {
	// 						 trainers(
	// 							where: {joined_on: "2019-09-15"}
	// 						) @postgres {
	// 							_id
	// 							name
	// 							caught_pokemons(
	// 								where: {
	// 									trainer_id: "trainers.id"
	// 									type: "Fire"
	// 								}
	// 							) @postgres {
	// 								_id
	// 								name
	// 							}
	// 						}
	// 					}`,
	// 			Variables: nil,
	// 		},
	// 		token: "",
	// 	},
	// 	wantErr:    false,
	// 	wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "power_level": 100}}},
	// },
	{
		name: "Query: Sorting simple queries",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"id": 1, "name": 1},
						Sort:       []string{"name"},
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers(
									sort : ["name"]
								) @db {
									id
									name
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},
	{
		name: "Query: Sorting by multiple fields",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "caught_pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"id": 1, "name": 1, "caught_on": 1},
						Sort:       []string{"name", "-caught_on"},
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"caught_pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash", "caught_on": "2019-06-01"}, map[string]interface{}{"id": "2", "name": "james", "caught_on": "2019-06-01"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "caught_pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								caught_pokemons(
									sort : ["name","-caught_on"]
								) @db {
									id
									name
									caught_on
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"caught_pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "ash", "caught_on": "2019-06-01"}, map[string]interface{}{"id": "2", "name": "james", "caught_on": "2019-06-01"}}},
	},
	// {
	// 	name: "Query: Sorting nested fields",
	// 	crudMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method:         "IsPreparedQueryPresent",
	// 			args:           []interface{}{"db", "caught_pokemons"},
	// 			paramsReturned: []interface{}{false},
	// 		},
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method: "Read",
	// 			args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.ReadRequest{
	//			Find:      map[string]interface{}{},
	// 				Aggregate: map[string][]string{},
	// 				GroupBy:   []interface{}{},
	// 				Operation: utils.All,
	// 				Options: &model.ReadOptions{
	// 					Select:     map[string]int32{"id": 1, "name": 1, "caught_on": 1},
	// 					Sort:       []string{"name", "-caught_on"},
	// 					HasOptions: true,
	// 				},
	// 				IsBatch: true,
	// 			}},
	// 			paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash", "caught_on": "2019-06-01"}, map[string]interface{}{"id": "2", "name": "james", "caught_on": "2019-06-01"}}, nil},
	// 		},
	// 	},
	// 	schemaMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetSchema",
	// 			args:           []interface{}{"db", "caught_pokemons"},
	// 			paramsReturned: []interface{}{model.Fields{}, true},
	// 		},
	// 	},
	// 	authMockArgs: []mockArgs{
	// 		{
	// 			method:         "IsReadOpAuthorised",
	// 			args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
	// 		},
	// 		{
	// 			method:         "PostProcessMethod",
	// 			args:           []interface{}{mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{nil},
	// 		},
	// 	},
	// 	args: args{
	// 		req: &model.GraphQLRequest{
	// 			OperationName: "query",
	// 			Query: `query {
	// 						caught_pokemons(
	// 							sort : ["name","-caught_on"]
	// 						) @db {
	// 							id
	// 							name
	// 							caught_on
	// 						}
	// 					}`,
	// 			Variables: nil,
	// 		},
	// 		token: "",
	// 	},
	// 	wantErr:    false,
	// 	wantResult: map[string]interface{}{"caught_pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "ash", "caught_on": "2019-06-01"}, map[string]interface{}{"id": "2", "name": "james", "caught_on": "2019-06-01"}}},
	// },
	{
		name: "Query: Distinct fields",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.Distinct,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"type": 1},
						Distinct:   &distinct,
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"type": "fire"}, map[string]interface{}{"type": "water"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								pokemons(
									distinct : "type"
								) @db {
									type
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"type": "fire"}, map[string]interface{}{"type": "water"}}},
	},
	{
		name: "Query: Pagination limit operator",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"id": 1, "name": 1},
						Limit:      &number,
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers(
									limit : 5
								) @db {
									id
									name
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},
	{
		name: "Query: Pagination skip operator",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"id": 1, "name": 1},
						Skip:       &number,
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers(
									skip : 5
								) @db {
									id
									name
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},

	{
		name: "Query: Pagination skip & limit operator",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"id": 1, "name": 1},
						Skip:       &number,
						Limit:      &number,
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers(
									skip : 5
									limit : 5
								) @db {
									id
									name
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},
	// {
	// 	name: "Query: Pagination skip & limit on nested queries",
	// 	crudMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method:         "IsPreparedQueryPresent",
	// 			args:           []interface{}{"db", "trainers"},
	// 			paramsReturned: []interface{}{false},
	// 		},
	// 		{
	// 			method:         "GetDBType",
	// 			args:           []interface{}{"db"},
	// 			paramsReturned: []interface{}{"postgres", nil},
	// 		},
	// 		{
	// 			method: "Read",
	// 			args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
	//			Find:      map[string]interface{}{},
	// 				Aggregate: map[string][]string{},
	// 				GroupBy:   []interface{}{},
	// 				Operation: utils.All,
	// 				Options: &model.ReadOptions{
	// 					Select:     map[string]int32{"id": 1, "name": 1},
	// 					Skip:       &number,
	// 					HasOptions: true,
	// 				},
	// 				IsBatch: true,
	// 			}},
	// 			paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
	// 		},
	// 	},
	// 	schemaMockArgs: []mockArgs{
	// 		{
	// 			method:         "GetSchema",
	// 			args:           []interface{}{"db", "trainers"},
	// 			paramsReturned: []interface{}{model.Fields{}, true},
	// 		},
	// 	},
	// 	authMockArgs: []mockArgs{
	// 		{
	// 			method:         "IsReadOpAuthorised",
	// 			args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
	// 		},
	// 		{
	// 			method:         "PostProcessMethod",
	// 			args:           []interface{}{mock.Anything, mock.Anything},
	// 			paramsReturned: []interface{}{nil},
	// 		},
	// 	},
	// 	args: args{
	// 		req: &model.GraphQLRequest{
	// 			OperationName: "query",
	// 			Query: `query {
	// 						trainers(
	// 							skip : 5
	// 						) @db {
	// 							id
	// 							name
	// 						}
	// 					}`,
	// 			Variables: nil,
	// 		},
	// 		token: "",
	// 	},
	// 	wantErr:    false,
	// 	wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	// },
	{
		name: "Query: Multiple Operations",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{"type": "water"},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select:     map[string]int32{"id": 1, "name": 1},
						Skip:       &number,
						Sort:       []string{"name"},
						HasOptions: true,
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers(
									where : {
										type : "water"
									}
									sort : ["name"]
									skip : 5
								) @db {
									id
									name
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},
	{
		name: "Query: Multiple queries in a single graphql request",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"type": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"type": "water", "name": "bulbasur"}, map[string]interface{}{"type": "fire", "name": "charmander"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers @db {
									id
									name
								}
								pokemons @db {
									name
									type
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"type": "water", "name": "bulbasur"}, map[string]interface{}{"type": "fire", "name": "charmander"}}, "trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},
	{
		name: "Query: Multiple queries in a single graphql request",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"type": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"type": "water", "name": "bulbasur"}, map[string]interface{}{"type": "fire", "name": "charmander"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers @db {
									id
									name
								}
								pokemons @db {
									name
									type
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"pokemons": []interface{}{map[string]interface{}{"type": "water", "name": "bulbasur"}, map[string]interface{}{"type": "fire", "name": "charmander"}}, "trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}},
	},
	{
		name: "Query: Same database joins",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Find:        map[string]interface{}{"trainer_id": "1"},
					Operation:   utils.All,
					Options:     &model.ReadOptions{},
					GroupBy:     []interface{}{},
					Aggregate:   map[string][]string{},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}, nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Find:        map[string]interface{}{"trainer_id": "2"},
					Operation:   utils.All,
					GroupBy:     []interface{}{},
					Options:     &model.ReadOptions{},
					Aggregate:   map[string][]string{},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", DBType: "db", From: "id", To: "trainer_id"}}}, true},
			},
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "trainer_id": &model.FieldType{FieldName: "trainer_id", IsFieldTypeRequired: true, Kind: model.TypeID, IsForeign: true, JointTable: &model.TableProperties{Table: "trainers", To: "id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers @db {
									id
									name
									pokemons {
										id
										name
									}
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}}, map[string]interface{}{"id": "2", "name": "james", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}}}},
	},
	{
		name: "Query: Performing joins on the fly",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"trainer_id": "2",
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}, nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"trainer_id": "1",
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "trainer_id": &model.FieldType{FieldName: "trainer_id", IsFieldTypeRequired: true, Kind: model.TypeID, IsForeign: true, JointTable: &model.TableProperties{Table: "trainers", To: "id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers @db {
									id
									name
									pokemons (
										where : {
											trainer_id : "trainers.id"
										}
									) @db {
										id
										name
									}
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}}, map[string]interface{}{"id": "2", "name": "james", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}}}},
	},
	{
		name: "Query: Cross database joins",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{false},
			},
			{
				method:         "GetDBType",
				args:           []interface{}{"mg"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "db", "trainers", &model.ReadRequest{
					Extras:    map[string]interface{}{},
					Find:      map[string]interface{}{},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"trainers": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}, nil},
			},
			{
				method:         "IsPreparedQueryPresent",
				args:           []interface{}{"mg", "pokemons"},
				paramsReturned: []interface{}{false},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "mg", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"trainer_id": "2",
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}, nil},
			},
			{
				method: "Read",
				args: []interface{}{mock.Anything, "mg", "pokemons", &model.ReadRequest{
					Extras: map[string]interface{}{},
					Find: map[string]interface{}{
						"trainer_id": "1",
					},
					Aggregate: map[string][]string{},
					GroupBy:   []interface{}{},
					Operation: utils.All,
					Options: &model.ReadOptions{
						Select: map[string]int32{"id": 1, "name": 1},
					},
					IsBatch:     true,
					PostProcess: map[string]*model.PostProcess{"pokemons": &model.PostProcess{}},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{[]interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}, nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
			{
				method:         "GetSchema",
				args:           []interface{}{"mg", "pokemons"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "trainer_id": &model.FieldType{FieldName: "trainer_id", IsFieldTypeRequired: true, Kind: model.TypeID, IsForeign: true, JointTable: &model.TableProperties{Table: "trainers", To: "id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsReadOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{&model.PostProcess{}, model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `query {
								trainers @db {
									id
									name
									pokemons (
										where : {
											trainer_id : "trainers.id"
										}
									) @mg {
										id
										name
									}
								}
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"trainers": []interface{}{map[string]interface{}{"id": "1", "name": "ash", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}}, map[string]interface{}{"id": "2", "name": "james", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "squirtle"}, map[string]interface{}{"id": "2", "name": "pikachu"}}}}},
	}}

var mutationTestCases = []tests{
	{
		name: "Mutation: Insert single object with templated directed",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db_t1"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db_t1", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db_t1", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "ParseToken",
				args:           []interface{}{mock.Anything, mock.Anything},
				paramsReturned: []interface{}{map[string]interface{}{"tenant": "t1"}, nil},
			},
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: [
								      {id: "1", name: "ash"}
								    ]
								  ) @template(value: "db_{{.auth.tenant}}") {
								    status
								    error
								    returning {
								      id
								      name
								    }
								  }
								}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert_trainers": map[string]interface{}{"error": nil, "status": 200, "returning": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}}}},
	},
	{
		name: "Mutation: Insert single object",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: [
								      {id: "1", name: "ash"}
								    ]
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								    }
								  }
								}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert_trainers": map[string]interface{}{"error": nil, "status": 200, "returning": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}}}},
	},
	{
		name: "Mutation: Insert single object error improper query",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation 
								  insert_trainers(
								    docs: [
								      {id: "1", name: "ash"}
								    ]
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								    }
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
		name: "Mutation: Insert single object with default & created at & updated at directives",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash", "age": 19}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"age": &model.FieldType{FieldName: "age", IsFieldTypeRequired: true, Kind: model.TypeInteger, IsDefault: true, Default: 19}, "id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: [
								      {id: "1", name: "ash"}
								    ]
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								    }
								  }
								}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert_trainers": map[string]interface{}{"error": nil, "status": 200, "returning": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}}}},
	},

	{
		name: "Mutation: Insert single object error request is not authorized",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, errors.New("request is not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: [
								      {id: "1", name: "ash"}
								    ]
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								    }
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
		name: "Mutation: Insert single object error invalid doc type provided",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, errors.New("request is not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: {}
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								    }
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
		name: "Mutation: Insert multiple objects",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Create",
				args: []interface{}{mock.Anything, "db", "trainers", &model.CreateRequest{
					Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: [
								      {id: "1", name: "ash"},
								      {id: "2", name: "james"}
								    ]
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								    }
								  }
								}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert_trainers": map[string]interface{}{"error": nil, "status": 200, "returning": []interface{}{map[string]interface{}{"id": "1", "name": "ash"}, map[string]interface{}{"id": "2", "name": "james"}}}},
	},
	{
		name: "Mutation: Insert an object along with its related objects through relationships",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Batch",
				args: []interface{}{mock.Anything, "db", &model.BatchRequest{
					Requests: []*model.AllRequest{
						{
							DBAlias:   "db",
							Col:       "trainers",
							Document:  []interface{}{map[string]interface{}{"id": "1", "name": "ash"}},
							Operation: string(utils.All),
							Type:      string(model.Create),
						},
						{
							DBAlias:   "db",
							Col:       "pokemons",
							Document:  []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "combat_power": 200, "trainer_id": "1"}, map[string]interface{}{"id": "2", "name": "charmender", "combat_power": 150, "trainer_id": "1"}},
							Operation: string(utils.All),
							Type:      string(model.Create),
						},
					},
				}, model.RequestParams{Resource: "db-batch"}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "trainers"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "pokemons": &model.FieldType{IsList: true, Kind: model.TypeObject, IsLinked: true, LinkedTable: &model.TableProperties{Table: "pokemons", From: "id", To: "trainer_id", DBType: "db"}}}, true},
			},
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "pokemons"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}, "combat_power": &model.FieldType{FieldName: "combat_power", Kind: model.TypeInteger}, "trainer_id": &model.FieldType{FieldName: "trainer_id", IsFieldTypeRequired: true, Kind: model.TypeID, IsForeign: true, JointTable: &model.TableProperties{Table: "trainers", To: "id"}}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_trainers(
								    docs: [
								      {
								        id: "1", 
								        name: "ash", 
								        pokemons: [
								          {
								            id: "1",
								            name: "pikachu",
								            combat_power: 200
								          },
								          {
								            id: "2",
								            name: "charmender",
								            combat_power: 150
								          }
								        ]
								      }
								    ]
								  ) @db {
								    status
								    error
								    returning {
								      id
								      name
								      city
								      pokemons {
								        id
								        name
								        combat_power
								      }
								    }
								  }
								}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert_trainers": map[string]interface{}{"error": nil, "status": 200, "returning": []interface{}{map[string]interface{}{"id": "1", "name": "ash", "pokemons": []interface{}{map[string]interface{}{"id": "1", "name": "pikachu", "combat_power": 200}, map[string]interface{}{"id": "2", "name": "charmender", "combat_power": 150}}}}}},
	},
	{
		name: "Mutation: Updating single object set operation incorrect value for set opeation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$set": "$data",
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, errors.New("query not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: {
										id: {
											_eq: "1"
											}
										},
								    set: $data
								  ) @db {
									error
								    status
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
		name: "Mutation: Updating single object set operation error update operation not authorized",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "my cool pikachu",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, errors.New("query not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: {
										id: {
											_eq: "1"
											}
										},
								    set: {
										name: "my cool pikachu"
									}
								  ) @db {
									error
								    status
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
		name: "Mutation: Updating single object set operation error incorrect value for where clause",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "my cool pikachu",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: $data,
								    set: {
										name: "my cool pikachu"
									}
								  ) @db {
									error
								    status
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
		name: "Mutation: Updating single object set operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "my cool pikachu",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: {
										id: {
											_eq: "1"
											}
										},
								    set: {
										name: "my cool pikachu"
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object increment operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$inc": map[string]interface{}{
							"combat_power": 50,
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: {
										id: {
											_eq: "1"
											}
										},
								    inc: {
										combat_power : 50
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object increment operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$inc": map[string]interface{}{
							"combat_power": 50,
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: {
										id: {
											_eq: "1"
											}
										},
								    inc: {
										combat_power : 50
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object multiply operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$mul": map[string]interface{}{
							"combat_power": 2,
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where: {
										id: {
											_eq: "1"
											}
										},
								    mul: {
										combat_power : 2
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object min operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$min": map[string]interface{}{
							"lowest_score": 2,
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    min : {
										lowest_score : 2
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object max operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"postgres", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$max": map[string]interface{}{
							"highest_score": 2,
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    max : {
										highest_score : 2
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object push operation only for mongodb",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$push": map[string]interface{}{
							"attacks": "thunderbolt",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    push : {
										attacks : "thunderbolt"
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object unset operation only for mongodb",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$unset": map[string]interface{}{
							"is_favourite": "",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    unset : {
										is_favourite : ""
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Updating single object rename operation only for mongodb",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "caught_pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$rename": map[string]interface{}{
							"is_favourite": "favourite",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_caught_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    rename : {
										is_favourite : "favourite"
									}
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
}
var upsertTestCases = []tests{
	{
		name: "Mutation: Upsert operation error invalid type provided for op",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "pikachu",
							"type": "electric",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    set : {
										name : "pikachu",
										type : "electric"
									}
									op : $data
								  ) @db {
									error
								    status
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
		name: "Mutation: Upsert operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Update",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "pikachu",
							"type": "electric",
						},
					},
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsUpdateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  update_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								    set : {
										name : "pikachu",
										type : "electric"
									}
									op : upsert
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"update_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
}

var deleteTestCases = []tests{
	{
		name: "Mutation: Delete operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Delete",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.DeleteRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsDeleteOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  delete_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								  ) @db {
									error
								    status
								  }
							}`,
				Variables: nil,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"delete_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
	{
		name: "Mutation: Delete operation error not authorized",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Delete",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.DeleteRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsDeleteOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, errors.New("operation not authorized")},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  delete_pokemons(
								    where : {
										id: {
											_eq: "1"
											}
										},
								  ) @db {
									error
								    status
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
		name: "Mutation: Delete operation error couldn't generate delete request as where clause is invalid",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Delete",
				args: []interface{}{mock.Anything, "db", "pokemons", &model.DeleteRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					Operation: utils.All,
				}, model.RequestParams{}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{},
		authMockArgs: []mockArgs{
			{
				method:         "IsDeleteOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  delete_pokemons(
								    where : 1
								  ) @db {
									error
								    status
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

var transactionTestCases = []tests{
	{
		name: "Mutation: transaction operation",
		crudMockArgs: []mockArgs{
			{
				method:         "GetDBType",
				args:           []interface{}{"db"},
				paramsReturned: []interface{}{"mongo", nil},
			},
			{
				method: "Batch",
				args: []interface{}{mock.Anything, "db", &model.BatchRequest{
					Requests: []*model.AllRequest{
						{
							DBAlias:   "db",
							Col:       "caught_pokemons",
							Document:  []interface{}{map[string]interface{}{"id": "5", "name": "chalizard"}, map[string]interface{}{"id": "6", "name": "squirtle"}},
							Operation: utils.All,
							Type:      string(model.Create),
						},
						{
							DBAlias: "db",
							Col:     "caught_pokemons",
							Find: map[string]interface{}{
								"id": "4",
							},
							Operation: utils.All,
							Type:      string(model.Delete),
						},
					},
				}, model.RequestParams{Resource: "db-batch"}},
				paramsReturned: []interface{}{nil},
			},
		},
		schemaMockArgs: []mockArgs{
			{
				method:         "GetSchema",
				args:           []interface{}{"db", "caught_pokemons"},
				paramsReturned: []interface{}{model.Fields{"id": &model.FieldType{FieldName: "id", IsFieldTypeRequired: true, IsPrimary: true, Kind: model.TypeID}, "name": &model.FieldType{FieldName: "name", Kind: model.TypeString}}, true},
			},
		},
		authMockArgs: []mockArgs{
			{
				method:         "IsCreateOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
			{
				method:         "IsDeleteOpAuthorised",
				args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
				paramsReturned: []interface{}{model.RequestParams{}, nil},
			},
		},
		args: args{
			req: &model.GraphQLRequest{
				OperationName: "query",
				Query: `mutation {
								  insert_caught_pokemons(
								    docs : [
										{ id : "5", name : "chalizard"},
										{ id : "6", name : "squirtle"},
										],
								  ) @db {
									error
								    status
								  }

                                  delete_caught_pokemons(
                                  	where : {
                                  		id : "4"
									}
                                  	) @db {
                                  		error 
                                  		status
                                  }
							}`,
			},
			token: "",
		},
		wantErr:    false,
		wantResult: map[string]interface{}{"insert_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}, "delete_caught_pokemons": map[string]interface{}{"error": nil, "status": 200}},
	},
}
