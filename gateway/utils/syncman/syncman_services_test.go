package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/stretchr/testify/mock"
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
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args:    args{ctx: context.Background(), project: "2", service: "serviceID", value: &config.Service{ID: "id"}},
			wantErr: true,
		},
		{
			name: "services are nil and unable to set services config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID", value: &config.Service{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid templating engine provided")},
				},
			},
			wantErr: true,
		},
		{
			name: "services are nil and unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID", value: &config.Service{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "services are nil and service is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID", value: &config.Service{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "unable to set services config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID", value: &config.Service{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid templating engine provided")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID", value: &config.Service{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "service is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID", value: &config.Service{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
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

			if err := tt.s.SetService(tt.args.ctx, tt.args.project, tt.args.service, tt.args.value); (err != nil) != tt.wantErr {
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
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args:    args{ctx: context.Background(), project: "2", service: "serviceID"},
			wantErr: true,
		},
		{
			name: "unable to set services config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid templating engine provided")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "service is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "serviceID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", service: "serviceID"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetServicesConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
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

			if err := tt.s.DeleteService(tt.args.ctx, tt.args.project, tt.args.service); (err != nil) != tt.wantErr {
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
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "id"}}}}}}}},
			args:    args{ctx: context.Background(), project: "2", serviceID: "id"},
			wantErr: true,
		},
		{
			name:    "service not present in config config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "id"}}}}}}}},
			args:    args{ctx: context.Background(), project: "1", serviceID: "notService"},
			wantErr: true,
		},
		{
			name: "got service",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "id"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", serviceID: "service"},
			want: []interface{}{&config.Service{ID: "id"}},
		},
		{
			name: "got all services",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Services: &config.ServicesModule{Services: config.Services{"service": &config.Service{ID: "id"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", serviceID: "*"},
			want: []interface{}{&config.Service{ID: "id"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetServices(tt.args.ctx, tt.args.project, tt.args.serviceID)
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
