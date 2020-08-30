package syncman

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_ApplyProjectConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project *config.Project
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		adminMockArgs   []mockArgs
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		want            int
		wantErr         bool
	}{
		{
			name: "sync operation is not valid",
			s:    &Manager{projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.Project{}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{config.GenerateEmptyConfig(), &config.Project{}},
					paramsReturned: []interface{}{false},
				},
			},
			want:    http.StatusUpgradeRequired,
			wantErr: true,
		},
		{
			name: "could not get internal access token",
			s:    &Manager{projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.Project{}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{config.GenerateEmptyConfig(), &config.Project{}},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("could not generate signed string for token")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project exists already and store type kube and can not set project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetProjectConfig",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("error marshalling project config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project doesn't exist and store type kube and can not set project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "2"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetProjectConfig",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("error marshalling project config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project exists already and store type kube and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetProjectConfig",
					args:           []interface{}{mock.Anything},
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
		{
			name: "project doesn't exist and store type kube and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "2"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetProjectConfig",
					args:           []interface{}{mock.Anything},
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

			mockAdmin := mockAdminSyncmanInterface{}
			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.adminMan = &mockAdmin
			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			got, err := tt.s.ApplyProjectConfig(tt.args.ctx, tt.args.project, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.ApplyProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.ApplyProjectConfig() = %v, want %v", got, tt.want)
			}

			mockAdmin.AssertExpectations(t)
			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_DeleteProjectConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		adminMockArgs   []mockArgs
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name: "could not get internal access token",
			s:    &Manager{storeType: "local"},
			args: args{ctx: context.Background(), projectID: "project"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("could not generate signed string for token")},
				},
			},
			wantErr: true,
		},
		{
			name: "store type kube and couldn't delete project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "notProject"}}}},
			args: args{ctx: context.Background(), projectID: "project"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "Delete",
					args:   []interface{}{"project"},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteProject",
					args:           []interface{}{mock.Anything, "project"},
					paramsReturned: []interface{}{errors.New("unable to get config map")},
				},
			},
			wantErr: true,
		},
		{
			name: "store type kube and project is deleted",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "notProject"}}}},
			args: args{ctx: context.Background(), projectID: "project"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "Delete",
					args:   []interface{}{"project"},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteProject",
					args:           []interface{}{mock.Anything, "project"},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAdmin := mockAdminSyncmanInterface{}
			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.adminMan = &mockAdmin
			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.DeleteProjectConfig(tt.args.ctx, tt.args.projectID, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.DeleteProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockAdmin.AssertExpectations(t)
			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetProjectConfig(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "project not present in state",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}, {ID: "2"}}}},
			args:    args{projectID: "3"},
			want:    []interface{}{},
			wantErr: true,
		},
		{
			name: "projectID is *",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}, {ID: "2"}}}},
			args: args{projectID: "*"},
			want: []interface{}{config.Project{ID: "1"}, config.Project{ID: "2"}},
		},
		{
			name: "projectID matches an existing project's id",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}, {ID: "2"}}}},
			args: args{projectID: "1"},
			want: []interface{}{config.Project{ID: "1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetProjectConfig(context.Background(), tt.args.projectID, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetProjectConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
