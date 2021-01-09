package eventing

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

func Test_deleteEventingConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to delete eventing config",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodPost,
						"/v1/config/projects/myproject/eventing/config/eventing-config",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						map[string]interface{}{
							"statusCode": 400,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Eventing config deleted succesfully",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodPost,
						"/v1/config/projects/myproject/eventing/config/eventing-config",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						map[string]interface{}{
							"statusCode": 200,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockTransport

			if err := deleteEventingConfig(tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("deleteEventingConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
		})
	}
}

func Test_deleteEventingTriggers(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "l"
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
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get eventing triggers",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one trigger but unable to delete trigger",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local-admin",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/triggers/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin"},
							Default: []string{"local-admin"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one trigger and trigger deleted successfully",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local-admin",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/triggers/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin"},
							Default: []string{"local-admin"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "Prefix matches multiple triggers but unable to survey trigger IDs",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local-admin",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple triggers but unable to delete trigger",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local-admin",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/triggers/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple triggers and trigger successfully deleted",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local-admin",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/triggers/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "Prefix does not match any triggers but unable to delete trigger",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/triggers",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local-admin",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
								map[string]interface{}{
									"type":    "mongodb",
									"retries": 2,
									"timeout": 10,
									"id":      "local",
									"url":     "/v1/config/projects/myproject/eventing/triggers",
								},
							},
						},
					},
				},
			},
			wantErr: true,
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

			if err := deleteEventingTriggers(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteEventingTriggers() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_deleteEventingSchemas(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "l"
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
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get eventing schemas",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one schema but unable to delete schema",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/schema/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin"},
							Default: []string{"local-admin"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one schema and schema deleted successfully",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/schema/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin"},
							Default: []string{"local-admin"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "Prefix matches multiple schemas but unable to survey schema IDs",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
								map[string]interface{}{
									"id":     "local",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple schemas but unable to delete schema",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
								map[string]interface{}{
									"id":     "local",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/schema/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple schemas and schema successfully deleted",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
								map[string]interface{}{
									"id":     "local",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/schema/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "Prefix does not match any schemas but unable to delete schema",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/schema",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
								map[string]interface{}{
									"id":     "local",
									"schema": "type subscribers { id: ID! @primary name: String!}",
								},
							},
						},
					},
				},
			},
			wantErr: true,
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

			if err := deleteEventingSchemas(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteEventingSchemas() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_deleteEventingRules(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "l"
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
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get eventing rules",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one rule but unable to delete rule",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":   "local-admin",
									"rule": "date",
									"eval": "==",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/rules/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin"},
							Default: []string{"local-admin"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one rule and rule deleted successfully",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":   "local-admin",
									"rule": "date",
									"eval": "==",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/rules/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin"},
							Default: []string{"local-admin"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "Prefix matches multiple rules but unable to survey rule IDs",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":   "local-admin",
									"rule": "date",
									"eval": "==",
								},
								map[string]interface{}{
									"id":   "local",
									"rule": "date",
									"eval": "==",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple rules but unable to delete rule",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":   "local-admin",
									"rule": "date",
									"eval": "==",
								},
								map[string]interface{}{
									"id":   "local",
									"rule": "date",
									"eval": "==",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/rules/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple rules and rule successfully deleted",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":   "local-admin",
									"rule": "date",
									"eval": "==",
								},
								map[string]interface{}{
									"id":   "local",
									"rule": "date",
									"eval": "==",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/eventing/rules/local-admin",
						map[string]string{},
						new(model.Response),
					},
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
					method: "AskOne",
					args: []interface{}{
						&survey.Select{
							Message: "Choose the resource ID: ",
							Options: []string{"local-admin", "local"},
							Default: []string{"local-admin", "local"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "Prefix does not match any rules but unable to delete rule",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/eventing/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":   "local-admin",
									"rule": "date",
									"eval": "==",
								},
								map[string]interface{}{
									"id":   "local",
									"rule": "date",
									"eval": "==",
								},
							},
						},
					},
				},
			},
			wantErr: true,
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

			if err := deleteEventingRules(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteEventingRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
