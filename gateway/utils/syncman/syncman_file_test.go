package syncman

import (
	"context"
	"errors"
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
