package remoteservices

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

func Test_deleteRemoteService(t *testing.T) {
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
						"/v1/config/projects/myproject/remote-service/service",
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
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
						"/v1/config/projects/myproject/remote-service/service/local-admin",
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
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
						"/v1/config/projects/myproject/remote-service/service/local-admin",
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
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
									},
								},
								map[string]interface{}{
									"id":  "local",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
			name: "Prefix matches multiple services but unable to delete service",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
									},
								},
								map[string]interface{}{
									"id":  "local",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
						"/v1/config/projects/myproject/remote-service/service/local-admin",
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
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
									},
								},
								map[string]interface{}{
									"id":  "local",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
						"/v1/config/projects/myproject/remote-service/service/local-admin",
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
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
									},
								},
								map[string]interface{}{
									"id":  "local",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
			name: "Prefix does not match any services but unable to delete service",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
									},
								},
								map[string]interface{}{
									"id":  "local",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
						"/v1/config/projects/myproject/remote-service/service/local-admin",
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
						"/v1/config/projects/myproject/remote-service/service",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":  "local-admin",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
									},
								},
								map[string]interface{}{
									"id":  "local",
									"url": "/v1/config/projects/myproject/remote-service/service",
									"endpoints": map[string]interface{}{
										"path": "/v1/config/projects/{project}/remote-service/service/{id}",
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
						"/v1/config/projects/myproject/remote-service/service/local-admin",
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

			if err := deleteRemoteService(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteRemoteService() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
