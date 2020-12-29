package services

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

func Test_deleteSecret(t *testing.T) {
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
			name: "Unable to get remote services",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
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
			name: "Prefix matches one service but unable to delete service",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/secrets/local-admin",
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
			name: "Prefix matches one service and service deleted successfully",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/secrets/local-admin",
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
			name: "Prefix matches multiple services but unable to survey service IDs",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
								map[string]interface{}{
									"id": "local",
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
			name: "Prefix matches multiple services but unable to delete service",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
								map[string]interface{}{
									"id": "local",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/secrets/local-admin",
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
			name: "Prefix matches multiple services and service successfully deleted",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
								map[string]interface{}{
									"id": "local",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/secrets/local-admin",
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
			name: "Prefix does not match any services and unable to survey service ID",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
								map[string]interface{}{
									"id": "local",
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
			name: "Prefix does not match any services but unable to delete service",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
								map[string]interface{}{
									"id": "local",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/secrets/local-admin",
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
			name: "Prefix does not match any services and service successfully deleted",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/secrets",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
								},
								map[string]interface{}{
									"id": "local",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/secrets/local-admin",
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

			if err := deleteSecret(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_deleteService(t *testing.T) {
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
		prefix  map[string]string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			name: "Unable to get remote services",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "local-admin", "version": "v1"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
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
			name: "Prefix matches one service but unable to delete service",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/services/local/v1",
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
							Options: []string{"local::v1"},
							Default: []string{"local::v1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local::v1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one service and service deleted successfully",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/services/local/v1",
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
							Options: []string{"local::v1"},
							Default: []string{"local::v1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local::v1"},
				},
			},
		},
		{
			name: "Prefix matches multiple services but unable to survey service IDs",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "localAdmin",
									"version": "v1",
								},
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
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
							Options: []string{"localAdmin::v1", "local::v1"},
							Default: []string{"localAdmin::v1", "local::v1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "localAdmin::v1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple services but unable to delete service",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "localAdmin",
									"version": "v1",
								},
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/services/localAdmin/v1",
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
							Options: []string{"localAdmin::v1", "local::v1"},
							Default: []string{"localAdmin::v1", "local::v1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::v1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple services and service successfully deleted",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "localAdmin",
									"version": "v1",
								},
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/services/localAdmin/v1",
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
							Options: []string{"localAdmin::v1", "local::v1"},
							Default: []string{"localAdmin::v1", "local::v1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::v1"},
				},
			},
		},
		{
			name: "Prefix does not match any services and unable to survey service ID",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "b"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "localAdmin",
									"version": "v1",
								},
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
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
							Options: []string{"localAdmin::v1", "local::v1"},
							Default: []string{"localAdmin::v1", "local::v1"}[0],
						},
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "localAdmin::v1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix does not match any services but unable to delete service",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "b"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "localAdmin",
									"version": "v1",
								},
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/services/localAdmin/v1",
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
							Options: []string{"localAdmin::v1", "local::v1"},
							Default: []string{"localAdmin::v1", "local::v1"}[0],
						},
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::v1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix does not match any services and service successfully deleted",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "b"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/services",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "localAdmin",
									"version": "v1",
								},
								map[string]interface{}{
									"id":      "local",
									"version": "v1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/services/localAdmin/v1",
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
							Options: []string{"localAdmin::v1", "local::v1"},
							Default: []string{"localAdmin::v1", "local::v1"}[0],
						},
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::v1"},
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

			if err := deleteService(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteService() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_deleteServiceRole(t *testing.T) {
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
		prefix  map[string]string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		surveyMockArgs    []mockArgs
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			name: "Unable to get remote service-roles",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "local-admin", "version": "v1"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
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
			name: "Prefix matches one service-roles but unable to delete service-roles",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "local",
									"id":      "role1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/service-roles/local/role1",
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
							Options: []string{"local::role1"},
							Default: []string{"local::role1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local::role1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches one service-role and service-roles deleted successfully",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "local",
									"id":      "role1",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/service-roles/local/role1",
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
							Options: []string{"local::role1"},
							Default: []string{"local::role1"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "local::role1"},
				},
			},
		},
		{
			name: "Prefix matches multiple service-roles but unable to survey service IDs",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "localAdmin",
									"id":      "role1",
								},
								map[string]interface{}{
									"service": "local",
									"id":      "role2",
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
							Options: []string{"localAdmin::role1", "local::role2"},
							Default: []string{"localAdmin::role1", "local::role2"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "localAdmin::role1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple service-roles but unable to delete service-roles",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "localAdmin",
									"id":      "role1",
								},
								map[string]interface{}{
									"service": "local",
									"id":      "role2",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/service-roles/localAdmin/role1",
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
							Options: []string{"localAdmin::role1", "local::role2"},
							Default: []string{"localAdmin::role1", "local::role2"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::role1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix matches multiple service-roles and service-role successfully deleted",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "l"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "localAdmin",
									"id":      "role1",
								},
								map[string]interface{}{
									"service": "local",
									"id":      "role2",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/service-roles/localAdmin/role1",
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
							Options: []string{"localAdmin::role1", "local::role2"},
							Default: []string{"localAdmin::role1", "local::role2"}[0],
						},
						&surveyMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::role1"},
				},
			},
		},
		{
			name: "Prefix does not match any service-roles and unable to survey service ID",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "b"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "localAdmin",
									"id":      "role1",
								},
								map[string]interface{}{
									"service": "local",
									"id":      "role2",
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
							Options: []string{"localAdmin::role1", "local::role2"},
							Default: []string{"localAdmin::role1", "local::role2"}[0],
						},
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "localAdmin::role1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix does not match any service-roles but unable to delete service-roles",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "b"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "localAdmin",
									"id":      "role1",
								},
								map[string]interface{}{
									"service": "local",
									"id":      "role2",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/service-roles/localAdmin/role1",
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
							Options: []string{"localAdmin::role1", "local::role2"},
							Default: []string{"localAdmin::role1", "local::role2"}[0],
						},
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::role1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Prefix does not match any service-roles and service-role successfully deleted",
			args: args{project: "myproject", prefix: map[string]string{"serviceId": "b"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/runner/myproject/service-roles",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"service": "localAdmin",
									"id":      "role1",
								},
								map[string]interface{}{
									"service": "local",
									"id":      "role2",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodDelete,
						"/v1/runner/myproject/service-roles/localAdmin/role1",
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
							Options: []string{"localAdmin::role1", "local::role2"},
							Default: []string{"localAdmin::role1", "local::role2"}[0],
						},
						&surveyNoMatchReturnValue,
					},
					paramsReturned: []interface{}{nil, "localAdmin::role1"},
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

			if err := deleteServiceRole(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteServiceRole() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
