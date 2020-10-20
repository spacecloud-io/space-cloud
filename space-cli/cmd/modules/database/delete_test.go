package database

import (
	"errors"
	"net/http"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func Test_deleteDBRules(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "a"
	// surveyNoMatchReturnValue stores the values returned from the survey when prefix is not matched
	surveyNoMatchReturnValue := "b"
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		project string
		dbAlias string
		prefix  string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get db rules",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "users"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one table name but unable to delete db rules",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "users"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-users": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/users/rules", map[string]string{"dbAlias": "dbAlias", "col": "users"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one table name and rules deleted successfully",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "users"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-users": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/users/rules", map[string]string{"dbAlias": "dbAlias", "col": "users"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "prefix matches multiple table names but unable to survey table name",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
								map[string]interface{}{
									"dbAlias-address": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"age", "address"}, Default: []string{"age", "address"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "age"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple table names but unable to delete db rules",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
								map[string]interface{}{
									"dbAlias-address": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/age/rules", map[string]string{"dbAlias": "dbAlias", "col": "age"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"age", "address"}, Default: []string{"age", "address"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "age"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple table names and db rules successfully deleted",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
								map[string]interface{}{
									"dbAlias-address": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/age/rules", map[string]string{"dbAlias": "dbAlias", "col": "age"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"age", "address"}, Default: []string{"age", "address"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "age"},
				},
			},
		},
		{
			name: "prefix does not match any table names and unable to survey table name",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
								map[string]interface{}{
									"dbAlias-address": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"age", "address"}, Default: []string{"age", "address"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "age"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any table names but unable to delete rules",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
								map[string]interface{}{
									"dbAlias-address": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/age/rules", map[string]string{"dbAlias": "dbAlias", "col": "age"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"age", "address"}, Default: []string{"age", "address"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "age"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any table names and rules successfully deleted",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
								map[string]interface{}{
									"dbAlias-address": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/age/rules", map[string]string{"dbAlias": "dbAlias", "col": "age"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"age", "address"}, Default: []string{"age", "address"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "age"},
				},
			},
		},
		{
			name: "prefix does not match any table names of len 1 but unable to delete table name",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/age/rules", map[string]string{"dbAlias": "dbAlias", "col": "age"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any table names of len 1 and table name is succesfully deleted",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/collections/rules", map[string]string{"dbAlias": "dbAlias", "col": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"dbAlias-age": map[string]interface{}{
										"isRealtimeEnabled": false,
										"rules": map[string]interface{}{
											"id":   "ruleID",
											"rule": "rule",
										},
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/collections/age/rules", map[string]string{"dbAlias": "dbAlias", "col": "age"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.MocketAuthProviders{}
			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockTransport
			input.Survey = &mockSurvey

			if err := deleteDBRules(tt.args.project, tt.args.dbAlias, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteDBRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_deleteDBConfigs(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "m"
	// surveyNoMatchReturnValue stores the values returned from the survey when prefix is not matched
	surveyNoMatchReturnValue := "b"
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project string
		prefix  string
	}
	tests := []struct {
		name              string
		args              args
		surveyMockArgs    []mockArgs
		transportMockArgs []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get db config",
			args: args{project: "myproject", prefix: "mongo"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one alias but unable to delete config",
			args: args{project: "myproject", prefix: "mongo"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/mongo/config/database-config", map[string]string{"dbAlias": "mongo"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one alias and config deleted sucessfully",
			args: args{project: "myproject", prefix: "mongo"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/mongo/config/database-config", map[string]string{"dbAlias": "mongo"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "prefix matches multiple aliases but unable to survey alias",
			args: args{project: "myproject", prefix: "m"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
								map[string]interface{}{
									"mysql": map[string]interface{}{
										"enabled": false,
										"conn":    "root:my-secret-pw@tcp(localhost:3306)/",
										"type":    "mysql",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"mongo", "mysql"}, Default: []string{"mongo", "mysql"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "mongo"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple aliases but unable to delete db rules",
			args: args{project: "myproject", prefix: "m"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
								map[string]interface{}{
									"mysql": map[string]interface{}{
										"enabled": false,
										"conn":    "root:my-secret-pw@tcp(localhost:3306)/",
										"type":    "mysql",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/mongo/config/database-config", map[string]string{"dbAlias": "mongo"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"mongo", "mysql"}, Default: []string{"mongo", "mysql"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "mongo"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple aliases and db config successfully deleted",
			args: args{project: "myproject", prefix: "m"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
								map[string]interface{}{
									"mysql": map[string]interface{}{
										"enabled": false,
										"conn":    "root:my-secret-pw@tcp(localhost:3306)/",
										"type":    "mysql",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/mongo/config/database-config", map[string]string{"dbAlias": "mongo"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"mongo", "mysql"}, Default: []string{"mongo", "mysql"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "mongo"},
				},
			},
		},
		{
			name: "prefix does not match any aliases and unable to survey alias",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
								map[string]interface{}{
									"mysql": map[string]interface{}{
										"enabled": false,
										"conn":    "root:my-secret-pw@tcp(localhost:3306)/",
										"type":    "mysql",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"mongo", "mysql"}, Default: []string{"mongo", "mysql"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "mongo"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any aliases but unable to config",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
								map[string]interface{}{
									"mysql": map[string]interface{}{
										"enabled": false,
										"conn":    "root:my-secret-pw@tcp(localhost:3306)/",
										"type":    "mysql",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/mongo/config/database-config", map[string]string{"dbAlias": "mongo"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"mongo", "mysql"}, Default: []string{"mongo", "mysql"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "mongo"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any aliases and config deleted successfully",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/config", map[string]string{"dbAlias": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"mongo": map[string]interface{}{
										"enabled": false,
										"conn":    "mongodb://localhost:27017",
										"type":    "mongo",
										"name":    "dbName",
									},
								},
								map[string]interface{}{
									"mysql": map[string]interface{}{
										"enabled": false,
										"conn":    "root:my-secret-pw@tcp(localhost:3306)/",
										"type":    "mysql",
										"name":    "dbName",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/mongo/config/database-config", map[string]string{"dbAlias": "mongo"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"mongo", "mysql"}, Default: []string{"mongo", "mysql"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "mongo"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.MocketAuthProviders{}
			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockTransport
			input.Survey = &mockSurvey

			if err := deleteDBConfigs(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteDBConfigs() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockSurvey.AssertExpectations(t)
			mockTransport.AssertExpectations(t)
		})
	}
}

func Test_deleteDBPreparedQuery(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "p"
	// surveyNoMatchReturnValue stores the values returned from the survey when prefix is not matched
	surveyNoMatchReturnValue := "b"
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		project string
		dbAlias string
		prefix  string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get db prepared queries",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "prep1"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one prepared query but unable to delete prepared query",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "prep1"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/prepared-queries/prep1", map[string]string{"dbAlias": "dbAlias", "id": "prep1"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one prepared query and prepared query deleted successfully",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "prep1"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/prepared-queries/prep1", map[string]string{"dbAlias": "dbAlias", "id": "prep1"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "prefix matches multiple prepared queries but unable to survey prepared query",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "p"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
								map[string]interface{}{
									"id":  "prep2",
									"db":  "dbAlias",
									"sql": "select * from age",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"prep1", "prep2"}, Default: []string{"prep1", "prep2"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "prep1"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple prepared queries but unable to delete prepared query",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "p"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
								map[string]interface{}{
									"id":  "prep2",
									"db":  "dbAlias",
									"sql": "select * from age",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/prepared-queries/prep1", map[string]string{"dbAlias": "dbAlias", "id": "prep1"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"prep1", "prep2"}, Default: []string{"prep1", "prep2"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "prep1"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple prepared queries and prepared query deleted successfully",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "p"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
								map[string]interface{}{
									"id":  "prep2",
									"db":  "dbAlias",
									"sql": "select * from age",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/prepared-queries/prep1", map[string]string{"dbAlias": "dbAlias", "id": "prep1"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"prep1", "prep2"}, Default: []string{"prep1", "prep2"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "prep1"},
				},
			},
		},
		{
			name: "prefix does not match any prepared queries and unable to survey prepared query",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
								map[string]interface{}{
									"id":  "prep2",
									"db":  "dbAlias",
									"sql": "select * from age",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"prep1", "prep2"}, Default: []string{"prep1", "prep2"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "prep1"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any prepared queries but unable to delete prepared query",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
								map[string]interface{}{
									"id":  "prep2",
									"db":  "dbAlias",
									"sql": "select * from age",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/prepared-queries/prep1", map[string]string{"dbAlias": "dbAlias", "id": "prep1"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"prep1", "prep2"}, Default: []string{"prep1", "prep2"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "prep1"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any prepared queries and prepared query successfully deleted",
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "prep1",
									"db":  "dbAlias",
									"sql": "select * from users",
								},
								map[string]interface{}{
									"id":  "prep2",
									"db":  "dbAlias",
									"sql": "select * from age",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/database/dbAlias/prepared-queries/prep1", map[string]string{"dbAlias": "dbAlias", "id": "prep1"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 200,
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"prep1", "prep2"}, Default: []string{"prep1", "prep2"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "prep1"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.MocketAuthProviders{}
			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockTransport
			input.Survey = &mockSurvey

			if err := deleteDBPreparedQuery(tt.args.project, tt.args.dbAlias, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteDBRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
