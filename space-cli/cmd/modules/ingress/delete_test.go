package ingress

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

func Test_deleteIngressGlobalConfig(t *testing.T) {
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
			name: "Unable to delete ingress global config",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/routing/ingress/global", map[string]string{}, new(model.Response)},
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
			name: "Ingress global config deleted successfully",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{http.MethodDelete, "/v1/config/projects/myproject/routing/ingress/global", map[string]string{}, new(model.Response)},
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

			if err := deleteIngressGlobalConfig(tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("deleteIngressGlobalConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
		})
	}
}

func Test_deleteIngressRoute(t *testing.T) {
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
			name: "Unable to get ingress routes",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
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
			name: "Prefix matches one route but unable to delete route",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
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
						"/v1/config/projects/myproject/routing/ingress/local-admin",
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
			wantErr: true,
		},
		{
			name: "Prefix matches one route and route deleted successfully",
			args: args{project: "myproject", prefix: "local-admin"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
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
						"/v1/config/projects/myproject/routing/ingress/local-admin",
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
		},
		{
			name: "Prefix matches multiple routes but unable to survey route IDs",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
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
			name: "Prefix matches multiple routes but unable to delete route",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
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
						"/v1/config/projects/myproject/routing/ingress/local-admin",
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
			name: "Prefix matches multiple routes and route successfully deleted",
			args: args{project: "myproject", prefix: "l"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
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
						"/v1/config/projects/myproject/routing/ingress/local-admin",
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
			name: "Prefix does not match any routes and unable to survey route ID",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
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
			name: "Prefix does not match any routes but unable to delete route",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
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
						"/v1/config/projects/myproject/routing/ingress/local-admin",
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
			name: "Prefix does not match any routes and route successfully deleted",
			args: args{project: "myproject", prefix: "b"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/myproject/routing/ingress",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"id": "local-admin",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
									},
								},
								map[string]interface{}{
									"id": "local",
									"source": map[string]interface{}{
										"url":   "/v1/config/projects/myproject/routing/ingress",
										"hosts": []string{"www.google.com", "www.facebook.com"},
									},
									"targets": []interface{}{
										map[string]interface{}{
											"version": "v0.18.0",
											"host":    "greeting.myproject.svc.cluster.local",
										},
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
						"/v1/config/projects/myproject/routing/ingress/local-admin",
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

			if err := deleteIngressRoute(tt.args.project, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("deleteIngressRoute() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
