package auth

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

func Test_deleteAuthProvider(t *testing.T) {
	// surveyMatchReturnValue stores the values returned from the survey when prefix is matched
	surveyMatchReturnValue := "loc"
	// surveyNoMatchReturnValue stores the values returned from the survey when prefix is not matched
	surveyNoMatchReturnValue := "a"
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
			name: "Unable to get auth providers",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one provider but unable to delete provider",
			args: args{project: "myproject", prefix: "loc"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local-admin", map[string]string{"id": "local-admin"}, new(model.Response)},
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
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin"}, Default: []string{"local-admin"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one provider and provider is succesfully deleted",
			args: args{project: "myproject", prefix: "loc"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local-admin", map[string]string{"id": "local-admin"}, new(model.Response)},
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
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin"}, Default: []string{"local-admin"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "prefix matches multiple providers but unable to survey provider",
			args: args{project: "myproject", prefix: "loc"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin", "local"}, Default: []string{"local-admin", "local"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple providers but unable to delete provider",
			args: args{project: "myproject", prefix: "loc"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local", map[string]string{"id": "local"}, new(model.Response)},
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
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin", "local"}, Default: []string{"local-admin", "local"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "local"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple providers and provider is succesfully deleted",
			args: args{project: "myproject", prefix: "loc"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local", map[string]string{"id": "local"}, new(model.Response)},
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
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin", "local"}, Default: []string{"local-admin", "local"}[0]}, &surveyMatchReturnValue},
					paramsReturned: []interface{}{nil, "local"},
				},
			},
		},
		{
			name: "prefix does not match any providers and unable to survey provider",
			args: args{project: "myproject", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin", "local"}, Default: []string{"local-admin", "local"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any providers but unable to delete provider",
			args: args{project: "myproject", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local", map[string]string{"id": "local"}, new(model.Response)},
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
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin", "local"}, Default: []string{"local-admin", "local"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "local"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix does not match any providers and provider is succesfully deleted",
			args: args{project: "myproject", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local-admin",
									"enabled": true,
									"secret":  "hello",
								},
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local", map[string]string{"id": "local"}, new(model.Response)},
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
					args:           []interface{}{&survey.Select{Message: "Choose the resource ID: ", Options: []string{"local-admin", "local"}, Default: []string{"local-admin", "local"}[0]}, &surveyNoMatchReturnValue},
					paramsReturned: []interface{}{nil, "local"},
				},
			},
		},
		{
			name: "prefix does not match any providers of len 1 but unable to delete provider",
			args: args{project: "myproject", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local", map[string]string{"id": "local"}, new(model.Response)},
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
			args: args{project: "myproject", prefix: "a"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodGet, "/v1/config/projects/myproject/user-management/provider", map[string]string{"id": "*"}, new(model.Response)},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id":      "local",
									"enabled": true,
									"secret":  "hello",
								},
							},
						},
					},
				},
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/user-management/provider/local", map[string]string{"id": "local"}, new(model.Response)},
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

			if err := deleteAuthProvider(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteAuthProvider() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
