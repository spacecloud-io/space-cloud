package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
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
		value   *config.FileStore
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{}}}}}},
			args:    args{ctx: context.Background(), project: "2", value: &config.FileStore{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Rules: []*config.FileRule{}, Secret: "secret", StoreType: "local"}},
			wantErr: true,
		},
		{
			name: "unable to set filestore config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{}}}}}},
			args: args{ctx: context.Background(), project: "1", value: &config.FileStore{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Rules: []*config.FileRule{}, Secret: "secret", StoreType: "local"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("cannot get secrets from runner")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{}}}}}},
			args: args{ctx: context.Background(), project: "1", value: &config.FileStore{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Rules: []*config.FileRule{}, Secret: "secret", StoreType: "local"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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
			name: "filestore is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{}}}}}},
			args: args{ctx: context.Background(), project: "1", value: &config.FileStore{Enabled: true, Bucket: "bucket", Conn: "conn", Endpoint: "endpoint", Rules: []*config.FileRule{}, Secret: "secret", StoreType: "local"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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

			if err := tt.s.SetFileStore(tt.args.ctx, tt.args.project, tt.args.value); (err != nil) != tt.wantErr {
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
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{}}}}}},
			args:    args{ctx: context.Background(), project: "2"},
			wantErr: true,
		},
		{
			name: "got filestore config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{}}}}}},
			args: args{ctx: context.Background(), project: "1"},
			want: []interface{}{config.FileStore{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetFileStoreConfig(tt.args.ctx, tt.args.project)
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
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args:    args{ctx: context.Background(), id: "ruleID", project: "2", value: &config.FileRule{ID: "id"}},
			wantErr: true,
		},
		{
			name: "rules does not exist and unable to set file store config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("cannot get secrets from runner")},
				},
			},
			wantErr: true,
		},
		{
			name: "rules does not exist and unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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
			name: "rules does not exist and file rule is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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
			name: "unable to set file store config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), id: "ruleID", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("cannot get secrets from runner")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), id: "ruleID", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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
			name: "file rule is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), id: "ruleID", project: "1", value: &config.FileRule{ID: "id"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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

			if err := tt.s.SetFileRule(tt.args.ctx, tt.args.project, tt.args.id, tt.args.value); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args:    args{ctx: context.Background(), filename: "ruleID", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set file store config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), filename: "ruleID", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("cannot get secrets from runner")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), filename: "ruleID", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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
			name: "file rule deleted",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), filename: "ruleID", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetFileStoreConfig",
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

			if err := tt.s.SetDeleteFileRule(tt.args.ctx, tt.args.project, tt.args.filename); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args:    args{ctx: context.Background(), project: "2", ruleID: "ruleID"},
			wantErr: true,
		},
		{
			name:    "file rule not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args:    args{ctx: context.Background(), project: "1", ruleID: "notRuleID"},
			wantErr: true,
		},
		{
			name: "got file store rule",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleID: "ruleID"},
			want: []interface{}{&config.FileRule{ID: "ruleID"}},
		},
		{
			name: "got all file store rule",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{FileStore: &config.FileStore{Rules: []*config.FileRule{{ID: "ruleID"}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleID: "*"},
			want: []interface{}{&config.FileRule{ID: "ruleID"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetFileStoreRules(tt.args.ctx, tt.args.project, tt.args.ruleID)
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
