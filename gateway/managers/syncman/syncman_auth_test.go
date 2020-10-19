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
		integrationArgs []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: make(config.Auths)}}}},
			args:    args{ctx: context.Background(), project: "2", provider: "provider", value: &config.AuthStub{ID: "1"}},
			wantErr: true,
		},
		{
			name: "userman config is not set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: make(config.Auths)}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider", value: &config.AuthStub{ID: "1"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetUsermanConfig",
					args:           []interface{}{mock.Anything, "1", config.Auths{config.GenerateResourceID("chicago", "1", config.ResourceAuthProvider, "provider"): &config.AuthStub{ID: "provider", Enabled: false}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceAuthProvider, "provider"), &config.AuthStub{ID: "provider"}},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockHookResponse{}},
				},
			},
			wantErr: true,
		},
		{
			name: "userman config is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: make(config.Auths)}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider", value: &config.AuthStub{ID: "1"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetUsermanConfig",
					args:           []interface{}{mock.Anything, "1", config.Auths{config.GenerateResourceID("chicago", "1", config.ResourceAuthProvider, "provider"): &config.AuthStub{ID: "provider", Enabled: false}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceAuthProvider, "provider"), &config.AuthStub{ID: "provider"}},
					paramsReturned: []interface{}{nil},
				},
			},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockHookResponse{}},
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
			tt.s.integrationMan = &mockIntegrationManager{skip: true}

			if _, err := tt.s.SetUserManagement(context.Background(), tt.args.project, tt.args.provider, tt.args.value, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetUserManagement() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetUserManagement(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx        context.Context
		project    string
		providerID string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		integrationArgs []mockArgs
		want            []interface{}
		wantErr         bool
	}{
		{
			name:    "unable to get project",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: make(config.Auths)}}}},
			args:    args{ctx: context.Background(), project: "2", providerID: "provider"},
			wantErr: true,
		},
		{
			name: "providerID is empty",
			s:    &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: map[string]*config.AuthStub{"provider": {ID: "id"}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "*"},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockHookResponse{}},
				},
			},
			want: []interface{}{&config.AuthStub{ID: "id"}},
		},
		{
			name: "providerID is not present in config",
			s:    &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: map[string]*config.AuthStub{"provider": {ID: "id"}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "notProvider"},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockHookResponse{}},
				},
			},
			wantErr: true,
		},
		{
			name: "providerID is present in config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: map[string]*config.AuthStub{config.GenerateResourceID("chicago", "1", config.ResourceAuthProvider, "provider"): {ID: "provider"}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "provider"},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockHookResponse{}},
				},
			},
			want: []interface{}{&config.AuthStub{ID: "provider"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.s.integrationMan = &mockIntegrationManager{skip: true}

			_, got, err := tt.s.GetUserManagement(context.Background(), tt.args.project, tt.args.providerID, model.RequestParams{})
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
