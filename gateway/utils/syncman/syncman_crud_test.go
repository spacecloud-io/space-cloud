package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
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
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and unable to set crud config",
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "alias exists already and unable to set project",
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
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
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
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
			args: args{ctx: context.Background(), dbAlias: "", id: "responseID", project: "1"},
			want: []interface{}{&response{ID: "key", SQL: "field"}},
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
			want: []interface{}{&response{ID: "key", SQL: "field"}},
		},
		{
			name: "id is empty",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "", project: "1"},
			want: []interface{}{&response{ID: "key", SQL: "field"}},
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
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "2", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "id", project: "1", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1", v: &config.PreparedQuery{ID: "queryID", SQL: "field"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
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
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "id", project: "1"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{PreparedQueries: map[string]*config.PreparedQuery{"key": {ID: "id", SQL: "field"}}}}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
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
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name: "collections in config is nil and unable to set crud config",
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
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
			s:    &Manager{storeType: "none", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"alias": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetCrudConfig",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
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
