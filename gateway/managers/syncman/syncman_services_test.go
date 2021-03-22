package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_SetService(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		service string
		value   *config.Service
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name: "Project config not found",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "test", service: "greeter", value: &config.Service{ID: "greeter"}},
			wantErr: true,
		},
		{
			name: "Services are nil and add a new service",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter", value: &config.Service{ID: "greeter"}},
			modulesMockArgs: []mockArgs{
				{
					method: "SetRemoteServiceConfig",
					args: []interface{}{mock.Anything, "myproject", config.Services{
						config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter"): &config.Service{
							ID: "greeter",
						},
					}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method: "SetResource",
					args: []interface{}{
						mock.Anything,
						config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter"),
						&config.Service{ID: "greeter"}},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
		},
		{
			name: "Services are nil and update existing service",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter", value: &config.Service{ID: "greeter", URL: "https://httpbin.org/"}},
			modulesMockArgs: []mockArgs{
				{
					method: "SetRemoteServiceConfig",
					args: []interface{}{mock.Anything, "myproject", config.Services{
						config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter"): &config.Service{
							ID:  "greeter",
							URL: "https://httpbin.org/",
						},
					}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method: "SetResource",
					args: []interface{}{
						mock.Anything,
						config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter"),
						&config.Service{ID: "greeter", URL: "https://httpbin.org/"}},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
		},
		{
			name: "unable to set remote service config",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter", value: &config.Service{ID: "greeter"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetRemoteServiceConfig",
					args:           []interface{}{mock.Anything, "myproject", mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid templating engine provided")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set resource",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter", value: &config.Service{ID: "greeter"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetRemoteServiceConfig",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter"), &config.Service{ID: "greeter"}},
					paramsReturned: []interface{}{errors.New("unable to set resource")},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			_, err := tt.s.SetService(tt.args.ctx, tt.args.project, tt.args.service, tt.args.value, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetService() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_DeleteService(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		service string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name: "Project config not found",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "test", service: "greeter"},
			wantErr: true,
		},
		{
			name: "Services exists, gets deleted",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							RemoteService: config.Services{
								config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter"): &config.Service{ID: "greeter"},
							},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetRemoteServiceConfig",
					args:           []interface{}{mock.Anything, "myproject", config.Services{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter")},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
		},
		{
			name: "unable to set remote service config",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetRemoteServiceConfig",
					args:           []interface{}{mock.Anything, "myproject", mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid templating engine provided")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to delete resource",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", service: "greeter"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetRemoteServiceConfig",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter")},
					paramsReturned: []interface{}{errors.New("unable to delete resource")},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.DeleteService(tt.args.ctx, tt.args.project, tt.args.service, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.DeleteService() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetServices(t *testing.T) {
	type args struct {
		ctx       context.Context
		project   string
		serviceID string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "Project config not found",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "test", serviceID: "greeter"},
			wantErr: true,
		},
		{
			name: "Get all services",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							RemoteService: config.Services{
								config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter1"): &config.Service{ID: "greeter1"},
							},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "myproject", serviceID: "*"},
			want:    []interface{}{&config.Service{ID: "greeter1"}},
			wantErr: false,
		},
		{
			name: "Get specific service",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							RemoteService: config.Services{
								config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter1"): &config.Service{ID: "greeter1"},
								config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter2"): &config.Service{ID: "greeter2"},
								config.GenerateResourceID("chicago", "myproject", config.ResourceRemoteService, "greeter3"): &config.Service{ID: "greeter3"},
							},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "myproject", serviceID: "greeter1"},
			want:    []interface{}{&config.Service{ID: "greeter1"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetServices(tt.args.ctx, tt.args.project, tt.args.serviceID, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetServices() = %v, want %v", got, tt.want)
			}
		})
	}
}
