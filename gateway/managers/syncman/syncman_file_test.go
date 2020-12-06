package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/stretchr/testify/mock"
)

func TestManager_SetFileStore(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		value   *config.FileStoreConfig
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
			args:    args{ctx: context.Background(), project: "2", value: &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
			wantErr: true,
		},
		{
			name: "Unable to set file store config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), project: "1", value: &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{mock.Anything, "1", &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
					paramsReturned: []interface{}{errors.New("cannot get secrets from runner")},
				},
			},
			wantErr: true,
		},
		{
			name: "Unable to set resource",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), project: "1", value: &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{mock.Anything, "1", &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreConfig, "filestore"), &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "File store config is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), project: "1", value: &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{mock.Anything, "1", &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreConfig, "filestore"), &config.FileStoreConfig{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Secret: "secret", StoreType: "local"}},
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
			tt.s.integrationMan = &mockIntegrationManager{skip: true}

			if _, err := tt.s.SetFileStore(tt.args.ctx, tt.args.project, tt.args.value, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetFileStore() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetFileStoreConfig(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "Unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args:    args{ctx: context.Background(), project: "2"},
			wantErr: true,
		},
		{
			name: "Got filestore config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{StoreType: "local"}}}}},
			args: args{ctx: context.Background(), project: "1"},
			want: []interface{}{&config.FileStoreConfig{StoreType: "local"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.integrationMan = &mockIntegrationManager{skip: true}
			_, got, err := tt.s.GetFileStoreConfig(tt.args.ctx, tt.args.project, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetFileStoreConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetFileStoreConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetFileRule(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		id      string
		value   *config.FileRule
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
			name:    "Project config not found",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args:    args{ctx: context.Background(), id: "ruleID", project: "2", value: &config.FileRule{ID: "id"}},
			wantErr: true,
		},
		{
			name: "rules does not exist and unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreSecurityRuleConfig",
					args:           []interface{}{mock.Anything, "1", config.FileStoreRules{config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "id"): &config.FileRule{ID: "id"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "id"), &config.FileRule{ID: "id"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "rules does not exist and file rule is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreSecurityRuleConfig",
					args:           []interface{}{mock.Anything, "1", config.FileStoreRules{config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "id"): &config.FileRule{ID: "id"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "id"), &config.FileRule{ID: "id"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "file rule is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreSecurityRuleConfig",
					args:           []interface{}{mock.Anything, "1", config.FileStoreRules{config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "id"): &config.FileRule{ID: "id"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "id"), &config.FileRule{ID: "id"}},
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
			tt.s.integrationMan = &mockIntegrationManager{skip: true}

			if _, err := tt.s.SetFileRule(tt.args.ctx, tt.args.project, tt.args.id, tt.args.value, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetFileRule() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetDeleteFileRule(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx      context.Context
		project  string
		filename string
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
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args:    args{ctx: context.Background(), filename: "ruleID", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to delete resource",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreRules: config.FileStoreRules{}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), filename: "ruleID", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreSecurityRuleConfig",
					args:           []interface{}{mock.Anything, "1", config.FileStoreRules{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "ruleID")},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "file rule deleted",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreRules: config.FileStoreRules{}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args: args{ctx: context.Background(), filename: "ruleID", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreSecurityRuleConfig",
					args:           []interface{}{mock.Anything, "1", config.FileStoreRules{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "ruleID")},
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
			tt.s.integrationMan = &mockIntegrationManager{skip: true}

			if _, err := tt.s.SetDeleteFileRule(tt.args.ctx, tt.args.project, tt.args.filename, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetDeleteFileRule() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetFileStoreRules(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		ruleID  string
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
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args:    args{ctx: context.Background(), project: "2", ruleID: "ruleID"},
			wantErr: true,
		},
		{
			name:    "file rule not present in config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}}}}},
			args:    args{ctx: context.Background(), project: "1", ruleID: "notRuleID"},
			wantErr: true,
		},
		{
			name: "got file store rule",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}, FileStoreRules: config.FileStoreRules{config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "ruleID"): &config.FileRule{ID: "ruleID"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleID: "ruleID"},
			want: []interface{}{&config.FileRule{ID: "ruleID"}},
		},
		{
			name: "got all file store rule",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, FileStoreConfig: &config.FileStoreConfig{}, FileStoreRules: config.FileStoreRules{config.GenerateResourceID("chicago", "1", config.ResourceFileStoreRule, "ruleID"): &config.FileRule{ID: "ruleID"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleID: "*"},
			want: []interface{}{&config.FileRule{ID: "ruleID"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.integrationMan = &mockIntegrationManager{skip: true}
			_, got, err := tt.s.GetFileStoreRules(tt.args.ctx, tt.args.project, tt.args.ruleID, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetFileStoreRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetFileStoreRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
