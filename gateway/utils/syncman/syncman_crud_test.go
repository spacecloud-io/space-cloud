package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/stretchr/testify/mock"
)

func TestManager_SetDeleteCollection(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
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
			name:    "unable to get project",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "notTableName", dbAlias: "notAlias", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "notTableName", dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "collection deleted successfully",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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

			if err := tt.s.SetDeleteCollection(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetDeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetDatabaseConnection(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		v       config.CrudStub
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
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "alias exists already and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "alias exists already and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
		},
		{
			name: "alias exists already and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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

			if err := tt.s.SetDatabaseConnection(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetDatabaseConnection() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_RemoveDatabaseConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
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
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("couldn't set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "database config is removed",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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

			if err := tt.s.RemoveDatabaseConfig(tt.args.ctx, tt.args.project, tt.args.dbAlias); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RemoveDatabaseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetLogicalDatabaseName(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "got db name",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{DBName: "DBName", Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			want: "DBName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetLogicalDatabaseName(tt.args.ctx, tt.args.project, tt.args.dbAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetLogicalDatabaseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.GetLogicalDatabaseName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetPreparedQuery(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		id      string
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "responseID", project: "2"},
			wantErr: true,
		},
		{
			name: "dbAlias is empty",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "*", id: "responseID", project: "1"},
			want: []interface{}{&preparedQueryResponse{ID: "key", SQL: "field"}},
		},
		{
			name:    "dbAlias is not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "responseID", project: "1"},
			wantErr: true,
		},
		{
			name:    "id is not empty but not present in prepared queries",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "notKey", project: "1"},
			wantErr: true,
		},
		{
			name: "id is not empty and present in prepared queries",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "key", project: "1"},
			want: []interface{}{&preparedQueryResponse{ID: "key", SQL: "field"}},
		},
		{
			name: "id is empty",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "*", project: "1"},
			want: []interface{}{&preparedQueryResponse{ID: "key", SQL: "field"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetPreparedQuery(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetPreparedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetPreparedQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetPreparedQueries(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		id      string
		v       *config.PreparedQuery
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
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "2", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "id", project: "1", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
			wantErr: true,
		},
		{
			name: "prepared queries are set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetPreparedQueries(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.id, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetPreparedQueries() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_RemovePreparedQueries(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		id      string
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
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "id", project: "1"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
			wantErr: true,
		},
		{
			name: "prepared queries is removed",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.RemovePreparedQueries(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RemovePreparedQueries() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetModifySchema(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
		schema  string
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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name: "collections in config is nil and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collections in config is nil and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
			wantErr: true,
		},
		{
			name: "collections in config is nil and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
		},
		{
			name: "table name doesn't exist and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "table name doesn't exist and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
			wantErr: true,
		},
		{
			name: "table name doesn't exist and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
		},
		{
			name: "table name exists and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "table name exists and project is not set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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
			wantErr: true,
		},
		{
			name: "table name exists and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetModifySchema(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetModifySchema() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetCollectionRules(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
		v       *config.TableRule
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
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			wantErr: true,
		},
		{
			name: "collection already present and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection already present and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "collection already present and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
		},
		{
			name: "collection not present and collectons is nil in config and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection not present and collectons is nil in config and project is not set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "collection not present and collectons is nil in config and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
		},
		{
			name: "collection not present and project is not set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			wantErr: true,
		},
		{
			name: "collection not present and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.TableRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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

			if err := tt.s.SetCollectionRules(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetCollectionRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetSecrets(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project    string
		secretName string
		key        string
	}
	tests := []struct {
		name          string
		s             *Manager
		args          args
		adminMockArgs []mockArgs
		want          string
		wantErr       bool
	}{
		{
			name: "unable to get internal access token",
			s:    &Manager{runnerAddr: "runnerAddr"},
			args: args{key: "key", project: "project", secretName: "secretName"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					args:           []interface{}{},
					paramsReturned: []interface{}{"", errors.New("unable to get signed string to get token")},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAdmin := mockAdminSyncmanInterface{}

			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.adminMan = &mockAdmin

			got, err := tt.s.GetSecrets(tt.args.project, tt.args.secretName, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.GetSecrets() = %v, want %v", got, tt.want)
			}

			mockAdmin.AssertExpectations(t)
		})
	}
}

func TestManager_GetSchemas(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "dbAlias and col are not empty but collection not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1"},
			wantErr: true,
		},
		{
			name: "dbAlias and col are not empty and got schemas",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbSchemaResponse{"alias-tableName": {Schema: "type event {id: ID! title: String}"}}},
		},
		{
			name: "dbAlias is not empty and got schemas",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbSchemaResponse{"alias-tableName": {Schema: "type event {id: ID! title: String}"}}},
		},
		{
			name: "dbAlias and col are empty and got schemas",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "*", project: "1"},
			want: []interface{}{map[string]*dbSchemaResponse{"alias-tableName": {Schema: "type event {id: ID! title: String}"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetSchemas(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSchemas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetSchemas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetReloadSchema(t *testing.T) {
	project := "1"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"create": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init())
	if err := s.SetConfig(rule, project); err != nil {
		t.Errorf("error setting config of schema - %s", err.Error())
	}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx       context.Context
		dbAlias   string
		project   string
		schemaArg *schema.Schema
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		want            map[string]interface{}
		wantErr         bool
	}{
		{
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "mongo", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notMongo", project: "1"},
			wantErr: true,
		},
		{
			name: "colName is default and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "mongo", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "colName is default and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "mongo", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			want:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "colName is default and project is set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "mongo", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
			want: map[string]interface{}{},
		},
		{
			name:    "unable to inspect schema",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "mongo", project: "1", schemaArg: s},
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

			got, err := tt.s.SetReloadSchema(tt.args.ctx, tt.args.dbAlias, tt.args.project, tt.args.schemaArg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetReloadSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.SetReloadSchema() = %v, want %v", got, tt.want)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetSchemaInspection(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
		schema  string
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
			name:    "unable to get project",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name: "collections nil and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collections nil and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collections nil and project is set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
		},
		{
			name: "collection not present and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection not present and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection not present and project is set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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
		},
		{
			name: "collection present and unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection present and unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection present and project is set",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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

			if err := tt.s.SetSchemaInspection(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetSchemaInspection() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_RemoveSchemaInspection(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "collections are nil in config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "schema inspection is removed",
			s:    &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
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

			if err := tt.s.RemoveSchemaInspection(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RemoveSchemaInspection() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_applySchemas(t *testing.T) {
	mockSchema := mockSchemaEventingInterface{}
	mockErrorSchema := mockSchemaEventingInterface{}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx           context.Context
		project       string
		dbAlias       string
		projectConfig *config.Project
		v             config.CrudStub
	}
	tests := []struct {
		name                string
		s                   *Manager
		args                args
		modulesMockArgs     []mockArgs
		schemaErrorMockArgs []mockArgs
		schemaMockArgs      []mockArgs
		wantErr             bool
	}{
		{
			name:    "database not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1", projectConfig: &config.Project{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"anotherTableName": {}}}},
			wantErr: true,
		},
		{
			name: "unable modify all schema",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"anotherTableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockErrorSchema},
				},
			},
			schemaErrorMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collections are nil and unable set crud config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "collections are nil and schemas are applied",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "unable set crud config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "schemas are applied",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.schemaErrorMockArgs {
				mockErrorSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules

			if err := tt.s.applySchemas(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.projectConfig, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Manager.applySchemas() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockErrorSchema.AssertExpectations(t)
			mockSchema.AssertExpectations(t)
		})
	}
}

func TestManager_SetModifyAllSchema(t *testing.T) {
	mockSchema := mockSchemaEventingInterface{}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		dbAlias string
		project string
		v       config.CrudStub
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		schemaMockArgs  []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			wantErr: true,
		},
		{
			name:    "unable to apply schemas",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
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
			wantErr: true,
		},
		{
			name: "modified all schema successfully",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", mock.Anything, mock.Anything},
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
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if err := tt.s.SetModifyAllSchema(tt.args.ctx, tt.args.dbAlias, tt.args.project, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetModifyAllSchema() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
			mockSchema.AssertExpectations(t)
		})
	}
}

func TestManager_GetDatabaseConfig(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Conn: "mongo:conn", Enabled: true, Type: "mongo"}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "db alias not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Conn: "mongo:conn", Enabled: true, Type: "mongo"}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "got db alias config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Conn: "mongo:conn", Enabled: true, Type: "mongo"}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			want: []interface{}{config.Crud{"alias": {Enabled: true, Conn: "mongo:conn", Type: "mongo"}}},
		},
		{
			name: "got services config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Conn: "mongo:conn", Enabled: true, Type: "mongo"}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "*", project: "1"},
			want: []interface{}{config.Crud{"alias": {Enabled: true, Conn: "mongo:conn", Type: "mongo"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetDatabaseConfig(tt.args.ctx, tt.args.project, tt.args.dbAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetDatabaseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetDatabaseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetCollectionRules(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		dbAlias string
		col     string
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}, Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "specified collection not present in config for dbAlias",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}, Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args:    args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1"},
			wantErr: true,
		},
		{
			name: "got collection rules for specific db alias and collection",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}, Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbRulesResponse{"alias-tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}}}},
		},
		{
			name: "col is empty and got collection rules",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}, Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbRulesResponse{"alias-tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}}}},
		},
		{
			name: "col and dbalias is empty and got collection rules",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}, Schema: "type event {id: ID! title: String}"}}}}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "*", project: "1"},
			want: []interface{}{map[string]*dbRulesResponse{"alias-tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"rule": {}}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetCollectionRules(tt.args.ctx, tt.args.project, tt.args.dbAlias, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetCollectionRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetCollectionRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
