package ingress

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func TestGetIngressRoutes(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project     string
		commandName string
		params      map[string]string
		filters     []string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		want              []*model.SpecObject
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			name: "Successful test",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": map[string]interface{}{
								"version": "v0.18.0",
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{"source": map[string]interface{}{"url": "/v1/config/projects/myproject/routing/ingress"}, "targets": map[string]interface{}{"version": "v0.18.0"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Url Filter Passes",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"url=/v1/config/projects/myproject/routing/ingress"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": map[string]interface{}{
								"version": "v0.18.0",
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{"source": map[string]interface{}{"url": "/v1/config/projects/myproject/routing/ingress"}, "targets": map[string]interface{}{"version": "v0.18.0"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Url Filter Fails",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"url=/v1/config/projects/myproject/routing"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": map[string]interface{}{
								"version": "v0.18.0",
							},
						},
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: false,
		},
		{
			name: "Service Filter Passes",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"service=greeting"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": []interface{}{
								map[string]interface{}{
									"version": "v0.18.0",
									"host":    "greeting.myproject.svc.cluster.local",
								},
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{
						"source": map[string]interface{}{"url": "/v1/config/projects/myproject/routing/ingress"},
						"targets": []interface{}{map[string]interface{}{
							"version": "v0.18.0",
							"host":    "greeting.myproject.svc.cluster.local",
						},
						}},
				},
			},
			wantErr: false,
		},
		{
			name: "Service Filter Fails",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"url=abc"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": []interface{}{
								map[string]interface{}{
									"version": "v0.18.0",
									"host":    "greeting.myproject.svc.cluster.local",
								},
							},
						},
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: false,
		},
		{
			name: "Request Host Filter Passes",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"request-host=www.google.com"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
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
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{
						"source": map[string]interface{}{
							"url":   "/v1/config/projects/myproject/routing/ingress",
							"hosts": []interface{}{"www.google.com", "www.facebook.com"},
						},
						"targets": []interface{}{map[string]interface{}{
							"version": "v0.18.0",
							"host":    "greeting.myproject.svc.cluster.local",
						},
						}},
				},
			},
			wantErr: false,
		},
		{
			name: "Request Host Filter Fails",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"request-host=abc"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
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
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: false,
		},

		{
			name: "Target Host Filter Passes",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
				filters:     []string{"target-host=greeting.myproject.svc.cluster.local", "target-host=basic.myproject.svc.cluster.local"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{
							map[string]interface{}{
								"id": "local-admin-1",
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
								"id": "local-admin-2",
								"source": map[string]interface{}{
									"url":   "/v1/config/projects/myproject/routing/ingress",
									"hosts": []string{"www.google.com", "www.facebook.com"},
								},
								"targets": []interface{}{
									map[string]interface{}{
										"version": "v0.18.0",
										"host":    "basic.myproject.svc.cluster.local",
									},
								},
							},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin-1"},
					Spec: map[string]interface{}{
						"source": map[string]interface{}{
							"url":   "/v1/config/projects/myproject/routing/ingress",
							"hosts": []interface{}{"www.google.com", "www.facebook.com"},
						},
						"targets": []interface{}{
							map[string]interface{}{
								"version": "v0.18.0",
								"host":    "greeting.myproject.svc.cluster.local",
							},
						},
					},
				},
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin-2"},
					Spec: map[string]interface{}{
						"source": map[string]interface{}{
							"url":   "/v1/config/projects/myproject/routing/ingress",
							"hosts": []interface{}{"www.google.com", "www.facebook.com"},
						},
						"targets": []interface{}{
							map[string]interface{}{
								"version": "v0.18.0",
								"host":    "basic.myproject.svc.cluster.local",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": map[string]interface{}{
								"version": "v0.18.0",
							},
						},
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSchema := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockSchema
			got, err := GetIngressRoutes(tt.args.project, tt.args.commandName, tt.args.params, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIngressRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetIngressRoutes() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetIngressRoutes() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetIngressGlobal(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project     string
		commandName string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		want              []*model.SpecObject
		wantErr           bool
	}{
		{
			name: "unable to get response",
			args: args{commandName: "ingress-global", project: "project"},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/project/routing/ingress/global", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{errors.New("unable to unmarshall"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"headers":    []interface{}{},
							"resHeaders": []interface{}{},
						}},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "got ingress global",
			args: args{commandName: "ingress-global", project: "project"},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/project/routing/ingress/global", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"headers":    []map[string]string{{"key": "key", "value": "value", "op": "option"}},
							"resHeaders": []map[string]string{{"key": "key", "value": "value", "op": "option"}},
						}},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/global",
					Type: "ingress-global",
					Meta: map[string]string{"project": "project"},
					Spec: map[string]interface{}{
						"headers":    []interface{}{map[string]interface{}{"key": "key", "value": "value", "op": "option"}},
						"resHeaders": []interface{}{map[string]interface{}{"key": "key", "value": "value", "op": "option"}},
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

			got, err := GetIngressGlobal(tt.args.project, tt.args.commandName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIngressGlobal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) > 0 {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetIngressGlobal() = %v, want %v", got, tt.want)
				}
			}

			mockTransport.AssertExpectations(t)
		})
	}
}
