package filestore

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

func Test_deleteFileStoreConfig(t *testing.T) {
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
			name: "Unable to delete filstore config",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodPost,
						"/v1/config/projects/myproject/file-storage/config/filestore-config",
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
			name: "Filstore config deleted successfully",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodPost,
						"/v1/config/projects/myproject/file-storage/config/filestore-config",
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

			if err := deleteFileStoreConfig(tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("deleteFileStoreConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
		})
	}
}

func Test_deleteFileStoreRule(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "l"
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
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to get filestore rules",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/file-storage/rules",
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
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/file-storage/rules/local-admin",
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
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/file-storage/rules/local-admin",
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
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
								map[string]interface{}{
									"id":     "local",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
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
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
								map[string]interface{}{
									"id":     "local",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/file-storage/rules/local-admin",
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
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
								map[string]interface{}{
									"id":     "local",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/file-storage/rules/local-admin",
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
			name: "Prefix does not match any rules and unable to survey rule ID",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
								map[string]interface{}{
									"id":     "local",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
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
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix does not match any rules but unable to delete rule",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
								map[string]interface{}{
									"id":     "local",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/file-storage/rules/local-admin",
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
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix does not match any rules and rule successfully deleted",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/file-storage/rules",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":     "local-admin",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
								map[string]interface{}{
									"id":     "local",
									"prefix": "hello",
									"rule": map[string]interface{}{
										"rule": "allow",
										"eval": "==",
									},
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/config/projects/myproject/file-storage/rules/local-admin",
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
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local-admin"},
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

			if err := deleteFileStoreRule(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteFileStoreRule() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
