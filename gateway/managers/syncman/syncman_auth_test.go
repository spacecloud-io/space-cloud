package syncman

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
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

			if _, err := tt.s.SetUserManagement(context.Background(), tt.args.project, tt.args.provider, tt.args.value, model.RequestParams{}); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: make(config.Auths)}}}},
			args:    args{ctx: context.Background(), project: "2", providerID: "provider"},
			wantErr: true,
		},
		{
			name: "providerID is empty",
			s:    &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: map[string]*config.AuthStub{"provider": {ID: "id"}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "*"},
			want: []interface{}{&config.AuthStub{ID: "id"}},
		},
		{
			name:    "providerID is not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: map[string]*config.AuthStub{"provider": {ID: "id"}}}}}},
			args:    args{ctx: context.Background(), project: "1", providerID: "notProvider"},
			wantErr: true,
		},
		{
			name: "providerID is present in config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, Auths: map[string]*config.AuthStub{config.GenerateResourceID("chicago", "1", config.ResourceAuthProvider, "provider"): {ID: "provider"}}}}}},
			args: args{ctx: context.Background(), project: "1", providerID: "provider"},
			want: []interface{}{&config.AuthStub{ID: "provider"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestManager_DeleteUserManagement(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx       context.Context
		project   string
		provider  string
		reqParams model.RequestParams
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		want            int
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{}}}}}},
			args:    args{ctx: context.Background(), project: "2", provider: "provider"},
			want:    http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "unable to set userman config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{"provider": &config.AuthStub{ID: "provider"}}}}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetUsermanConfig",
					args:           []interface{}{"1", config.Auth{}},
					paramsReturned: []interface{}{errors.New("unable to set userman config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{"provider": &config.AuthStub{ID: "provider"}}}}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetUsermanConfig",
					args:           []interface{}{"1", config.Auth{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "auth provider is succesfully deleted",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Auth: config.Auth{"provider": &config.AuthStub{ID: "provider"}}}}}}},
			args: args{ctx: context.Background(), project: "1", provider: "provider"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetUsermanConfig",
					args:           []interface{}{"1", config.Auth{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			want: http.StatusOK,
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

			got, err := tt.s.DeleteUserManagement(tt.args.ctx, tt.args.project, tt.args.provider, tt.args.reqParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.DeleteUserManagement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.DeleteUserManagement() = %v, want %v", got, tt.want)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}
