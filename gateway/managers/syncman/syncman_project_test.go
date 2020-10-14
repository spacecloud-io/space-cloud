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
		project *config.ProjectConfig
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
			name: "Unable to validate project sync operation",
			s:    &Manager{clusterID: "chicago", projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.ProjectConfig{}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{config.GenerateEmptyConfig(), &config.ProjectConfig{}},
					paramsReturned: []interface{}{false},
				},
			},
			want:    http.StatusUpgradeRequired,
			wantErr: true,
		},
		{
			name: "Could not get internal access token",
			s:    &Manager{clusterID: "chicago", projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.ProjectConfig{}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateProjectSyncOperation",
					args:           []interface{}{config.GenerateEmptyConfig(), &config.ProjectConfig{}},
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
			name: "project doesn't exist but cannot set resource",
			s:    &Manager{clusterID: "chicago", storeType: "kube", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}}}}},
			args: args{ctx: context.Background(), project: &config.ProjectConfig{ID: "2"}},
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
					args:           []interface{}{mock.Anything, &config.ProjectConfig{ID: "2", ContextTimeGraphQL: 10}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "2", config.ResourceProject, "2"), &config.ProjectConfig{ID: "2", ContextTimeGraphQL: 10}},
					paramsReturned: []interface{}{errors.New("error marshalling project config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project exists already so update existing project",
			s:    &Manager{clusterID: "chicago", storeType: "kube", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}}}}},
			args: args{ctx: context.Background(), project: &config.ProjectConfig{ID: "1"}},
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
					args:           []interface{}{mock.Anything, &config.ProjectConfig{ID: "1", ContextTimeGraphQL: 10}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceProject, "1"), &config.ProjectConfig{ID: "1", ContextTimeGraphQL: 10}},
					paramsReturned: []interface{}{nil},
				},
			},
			want: http.StatusOK,
		},
		{
			name: "project doesn't exist so add a new project",
			s:    &Manager{clusterID: "chicago", storeType: "kube", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}}}}},
			args: args{ctx: context.Background(), project: &config.ProjectConfig{ID: "2"}},
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
					args:           []interface{}{mock.Anything, &config.ProjectConfig{ID: "2", ContextTimeGraphQL: 10}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "2", config.ResourceProject, "2"), &config.ProjectConfig{ID: "2", ContextTimeGraphQL: 10}},
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
			name: "Could not get internal access token",
			s:    &Manager{clusterID: "chicago", storeType: "local"},
			args: args{ctx: context.Background(), projectID: "myproject"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("could not generate signed string for token")},
				},
			},
			wantErr: true,
		},
		{
			name: "Unable to delete existing project store throws error",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject1": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject1"},
						},
						"myproject2": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject2"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), projectID: "myproject1"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "Delete",
					args:   []interface{}{"myproject1"},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "myproject1", config.ResourceProject, "myproject1")},
					paramsReturned: []interface{}{errors.New("unable to get config map")},
				},
			},
			wantErr: true,
		},
		{
			name: "Delete an existing project",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject1": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject1"},
						},
						"myproject2": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject2"},
						},
					},
				},
			},
			args: args{ctx: context.Background(), projectID: "myproject1"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "Delete",
					args:   []interface{}{"myproject1"},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "myproject1", config.ResourceProject, "myproject1")},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
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
			name: "Get all project configs",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject1": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject1"},
						},
					},
				},
			},
			args: args{projectID: "*"},
			want: []interface{}{&config.ProjectConfig{ID: "myproject1"}},
		},
		{
			name: "Get specific project config",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject1": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject1"},
						},
						"myproject2": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject2"},
						},
					},
				},
			},
			args: args{projectID: "myproject1"},
			want: []interface{}{&config.ProjectConfig{ID: "myproject1"}},
		},
		{
			name: "Throw error when you are trying to fetch specific project which doesn't exists",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject1": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject1"},
						},
						"myproject2": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject2"},
						},
					},
				},
			},
			args:    args{projectID: "myproject3"},
			want:    []interface{}{},
			wantErr: true,
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
