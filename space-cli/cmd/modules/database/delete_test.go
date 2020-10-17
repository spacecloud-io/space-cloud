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
			args: args{project: "myproject", dbAlias: "dbAlias", prefix: "users "},
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
			name: "prefix matches one table name but unable to delete provider",
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
			name: "prefix matches multiple providers but unable to survey provider",
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
			name: "prefix does not match any providers of len 1 but unable to delete provider",
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
			name: "prefix does not match any providers of len 1 and provider is succesfully deleted",
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
