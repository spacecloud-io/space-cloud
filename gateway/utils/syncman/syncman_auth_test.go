package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/stretchr/testify/mock"
)

func TestManager_SetUserManagement(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx      context.Context
		project  string
		provider string
		value    *config.AuthStub
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{}}}}}},
			args:    args{ctx: context.Background(), project: "2", provider: "provider", value: &config.AuthStub{ID: "1"}},
			wantErr: true,
		},
		{
			name: "userman config is not set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{}}}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider", value: &config.AuthStub{ID: "1"}},
			modulesMockArgs: []mockArgs{
				{
					method: "SetUsermanConfig",
					args:   []interface{}{"1", config.Auth{"provider": &config.AuthStub{ID: "provider", Enabled: false}}},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "userman config is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{}}}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider", value: &config.AuthStub{ID: "1"}},
			modulesMockArgs: []mockArgs{
				{
					method: "SetUsermanConfig",
					args:   []interface{}{"1", config.Auth{"provider": &config.AuthStub{ID: "provider", Enabled: false}}},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
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

			if err := tt.s.SetUserManagement(tt.args.ctx, tt.args.project, tt.args.provider, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetUserManagement() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetUserManagement(t *testing.T) {
	type args struct {
		ctx        context.Context
		project    string
		providerID string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "unable to get project",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args:    args{ctx: context.Background(), project: "2", providerID: "provider"},
			wantErr: true,
		},
		{
			name: "providerID is empty",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{"provider": &config.AuthStub{ID: "id"}}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "*"},
			want: []interface{}{&config.AuthStub{ID: "id"}},
		},
		{
			name:    "providerID is not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{"provider": &config.AuthStub{ID: "id"}}}}}}},
			args:    args{ctx: context.Background(), project: "1", providerID: "notProvider"},
			wantErr: true,
		},
		{
			name: "providerID is present in config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{"provider": &config.AuthStub{ID: "id"}}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "provider"},
			want: []interface{}{&config.AuthStub{ID: "id"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetUserManagement(tt.args.ctx, tt.args.project, tt.args.providerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetUserManagement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetUserManagement() = %v, want %v", got, tt.want)
			}
		})
	}
}
