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
		v       *config.DatabaseConfig
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
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2", v: &config.DatabaseConfig{DbAlias: "alias"}},
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and unable to set crud config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: &config.DatabaseConfig{DbAlias: "alias"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "alias exists already and unable to set crud config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: &config.DatabaseConfig{DbAlias: "alias"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: &config.DatabaseConfig{DbAlias: "alias"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "notAlias"): &config.DatabaseConfig{DbAlias: "notAlias"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "notAlias"), &config.DatabaseConfig{DbAlias: "notAlias"}},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "alias exists already and unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: &config.DatabaseConfig{DbAlias: "alias"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"), mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "alias doesn't exist already and project is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: &config.DatabaseConfig{DbAlias: "alias"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"), &config.DatabaseConfig{DbAlias: "alias"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "alias exists already and project is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: &config.DatabaseConfig{DbAlias: "alias"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"), &config.DatabaseConfig{DbAlias: "alias"}},
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

			if _, err := tt.s.SetDatabaseConnection(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.v, model.RequestParams{}); (err != nil) != tt.wantErr {
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
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseConfigs{}},
					paramsReturned: []interface{}{errors.New("couldn't set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseConfigs{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias")},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "database config is removed",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseSchemas: config.DatabaseSchemas{}, DatabaseRules: config.DatabaseRules{}, DatabasePreparedQueries: config.DatabasePreparedQueries{}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseConfigs{}},
					paramsReturned: []interface{}{nil},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{}},
					paramsReturned: []interface{}{nil},
				},
				{
					method:         "SetDatabaseRulesConfig",
					args:           []interface{}{mock.Anything, config.DatabaseRules{}},
					paramsReturned: []interface{}{nil},
				},
				{
					method:         "SetDatabasePreparedQueryConfig",
					args:           []interface{}{mock.Anything, config.DatabasePreparedQueries{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias")},
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
			_, err := tt.s.RemoveDatabaseConfig(context.Background(), tt.args.project, tt.args.dbAlias, model.RequestParams{})
			if (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "got db name",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "DBName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			want: "DBName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetLogicalDatabaseName(context.Background(), tt.args.project, tt.args.dbAlias)
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
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "responseID", project: "2"},
			wantErr: true,
		},
		{
			name: "dbAlias is empty",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabasePreparedQueries: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "id"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "id", SQL: "field"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "*", id: "id", project: "1"},
			want: []interface{}{&preparedQueryResponse{ID: "id", DBAlias: "alias", SQL: "field"}},
		},
		{
			name:    "dbAlias is not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabasePreparedQueries: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "id"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "id", SQL: "field"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "id", project: "1"},
			wantErr: true,
		},
		{
			name:    "id is not empty but not present in prepared queries",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "DBName", DbAlias: "alias"}}, DatabasePreparedQueries: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "alias", "id"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "id", SQL: "field"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "notKey", project: "1"},
			wantErr: true,
		},
		{
			name: "id is not empty and present in prepared queries",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "DBName", DbAlias: "alias"}}, DatabasePreparedQueries: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "alias", "key"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "key", SQL: "field"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "key", project: "1"},
			want: []interface{}{&preparedQueryResponse{ID: "key", DBAlias: "alias", SQL: "field"}},
		},
		{
			name: "id is empty",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "DBName", DbAlias: "alias"}}, DatabasePreparedQueries: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "alias", "key"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "key", SQL: "field"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "*", project: "1"},
			want: []interface{}{&preparedQueryResponse{ID: "key", DBAlias: "alias", SQL: "field"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetPreparedQuery(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.id, model.RequestParams{})
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
		v       *config.DatbasePreparedQuery
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
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabasePreparedQueries: config.DatabasePreparedQueries{"resourceId": &config.DatbasePreparedQuery{DbAlias: "alias", ID: "id", SQL: "field"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", id: "id", project: "2", v: &config.DatbasePreparedQuery{ID: "queryID", SQL: "field"}},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabasePreparedQueries: config.DatabasePreparedQueries{"resourceId": &config.DatbasePreparedQuery{DbAlias: "alias", ID: "id", SQL: "field"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", id: "id", project: "1", v: &config.DatbasePreparedQuery{ID: "queryID", SQL: "field"}},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: map[string]*config.DatabaseConfig{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "queryID", project: "1", v: &config.DatbasePreparedQuery{ID: "queryID", SQL: "field"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabasePreparedQueryConfig",
					args:           []interface{}{mock.Anything, config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "alias", "queryID"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "queryID", SQL: "field"}}},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "prepared queries are set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "DBName", DbAlias: "alias"}}, ProjectConfig: &config.ProjectConfig{ID: "1"}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", id: "queryID", project: "1", v: &config.DatbasePreparedQuery{ID: "queryID", SQL: "field"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabasePreparedQueryConfig",
					args:           []interface{}{mock.Anything, config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "alias", "queryID"): &config.DatbasePreparedQuery{DbAlias: "alias", ID: "queryID", SQL: "field"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabasePreparedQuery, "alias", "queryID"), &config.DatbasePreparedQuery{DbAlias: "alias", ID: "queryID", SQL: "field"}},
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

			if _, err := tt.s.SetPreparedQueries(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.id, tt.args.v, model.RequestParams{}); (err != nil) != tt.wantErr {
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
			name: "Project config not found",
			s: &Manager{
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), dbAlias: "db", id: "fetchInstruments", project: "test"},
			wantErr: true,
		},
		{
			name: "DBAlias not found in config while removing prepared queries ",
			s: &Manager{
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							DatabaseConfigs: config.DatabaseConfigs{
								config.GenerateResourceID("", "myproject", config.ResourceDatabaseConfig, "db"): &config.DatabaseConfig{
									DbAlias: "db",
								},
							},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), dbAlias: "postgres", id: "fetchInstruments", project: "myproject"},
			wantErr: true,
		},
		{
			name: "Unable to set database prepared query config",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							DatabaseConfigs: config.DatabaseConfigs{
								config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "db"): &config.DatabaseConfig{
									DbAlias: "db",
								},
							},
							DatabasePreparedQueries: config.DatabasePreparedQueries{
								config.GenerateResourceID("chicago", "myproject", config.ResourceDatabasePreparedQuery, "db", "id"): &config.DatbasePreparedQuery{
									DbAlias: "db",
									ID:      "id",
								},
							},
						},
					},
				},
			},
			args: args{ctx: context.Background(), dbAlias: "db", id: "id", project: "myproject"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabasePreparedQueryConfig",
					args:           []interface{}{mock.Anything, config.DatabasePreparedQueries{}},
					paramsReturned: []interface{}{errors.New("unable to set database prepared query config")},
				},
			},
			wantErr: true,
		},
		{
			name: "Prepared query is removed from the config",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							DatabaseConfigs: config.DatabaseConfigs{
								config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "db"): &config.DatabaseConfig{
									DbAlias: "db",
								},
							},
							DatabasePreparedQueries: config.DatabasePreparedQueries{
								config.GenerateResourceID("chicago", "myproject", config.ResourceDatabasePreparedQuery, "db", "fetchInstruments"): &config.DatbasePreparedQuery{
									DbAlias: "db",
									SQL:     "select * from instruments;",
									ID:      "fetchInstruments",
								},
							},
						},
					},
				},
			},
			args: args{ctx: context.Background(), dbAlias: "db", id: "fetchInstruments", project: "myproject"},
			modulesMockArgs: []mockArgs{
				{
					method: "SetDatabasePreparedQueryConfig",
					args: []interface{}{
						mock.Anything, config.DatabasePreparedQueries{},
					},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "myproject", config.ResourceDatabasePreparedQuery, "db", "fetchInstruments")},
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

			if _, err := tt.s.RemovePreparedQueries(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.id, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RemovePreparedQueries() error = %v, wantErr %v", err, tt.wantErr)
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
		v       *config.DatabaseRule
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
			args:    args{ctx: context.Background(), project: "test", dbAlias: "alias", col: "table", v: &config.DatabaseRule{}},
			wantErr: true,
		},
		{
			name:    "Database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1", v: &config.DatabaseRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			wantErr: true,
		},
		{
			name: "collection already present and unable to set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", v: &config.DatabaseRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseRulesConfig",
					args:           []interface{}{mock.Anything, config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}}},
					paramsReturned: []interface{}{errors.New("error setting db module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection already present and unable to set resource",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", v: &config.DatabaseRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseRulesConfig",
					args:           []interface{}{mock.Anything, config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"), &config.DatabaseRule{Table: "tableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection already present and project is set",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", v: &config.DatabaseRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseRulesConfig",
					args:           []interface{}{mock.Anything, config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"), &config.DatabaseRule{Table: "tableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "collection not present and collectons is nil in config and unable to set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.DatabaseRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseRulesConfig",
					args:           []interface{}{mock.Anything, config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "notTableName", "rule"): &config.DatabaseRule{Table: "notTableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}}},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection not present and project is set",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", v: &config.DatabaseRule{Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseRulesConfig",
					args:           []interface{}{mock.Anything, config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "notTableName", "rule"): &config.DatabaseRule{Table: "notTableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "notTableName", "rule"), &config.DatabaseRule{Table: "notTableName", DbAlias: "alias", Rules: map[string]*config.Rule{"DB_INSERT": {ID: "rule1"}}}},
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

			if _, err := tt.s.SetCollectionRules(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.v, model.RequestParams{}); (err != nil) != tt.wantErr {
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
			s:    &Manager{clusterID: "chicago", runnerAddr: "runnerAddr"},
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

	mockSchema := mockSchemaEventingInterface{}
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
		format  string
	}
	tests := []struct {
		name                string
		s                   *Manager
		args                args
		modulesMockArgs     []mockArgs
		schemaErrorMockArgs []mockArgs
		schemaMockArgs      []mockArgs
		want                []interface{}
		wantErr             bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "dbAlias and col are not empty but collection not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1"},
			wantErr: true,
		},
		{
			name: "dbAlias and col are not empty and got schemas",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbSchemaResponse{"alias-tableName": {Schema: "type event {id: ID! title: String}"}}},
		},
		{
			name: "dbAlias is not empty and got schemas",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbSchemaResponse{"alias-tableName": {Schema: "type event {id: ID! title: String}"}}},
		},
		{
			name: "dbAlias and col are empty and got schemas",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "*", project: "1"},
			want: []interface{}{map[string]*dbSchemaResponse{"alias-tableName": {Schema: "type event {id: ID! title: String}"}}},
		},
		{
			name: "dbAlias and col are not empty and format JSON and got schemas",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", format: "json"},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method: "GetSchema",
					args:   []interface{}{"alias", "tableName"},
					paramsReturned: []interface{}{model.Fields{
						"alias": &model.FieldType{
							FieldName: "abcd",
						},
					}, true},
				},
			},
			want: []interface{}{map[string]*dbJSONSchemaResponse{"alias-tableName": {Fields: []*model.FieldType{
				{
					FieldName: "abcd",
				},
			}}}},
		},
		{
			name: "dbAlias is not empty and format is JSON and got schemas",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{"": &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "alias", project: "1", format: "json"},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
			},
			want: []interface{}{map[string]*dbJSONSchemaResponse{"alias-tableName": {Fields: []*model.FieldType{
				{
					FieldName: "abcd",
				},
			}}}},
		},
		{
			name: "dbAlias and col are empty and format is JSON and got schemas",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{"": &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "*", project: "1", format: "json"},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
			},
			want: []interface{}{map[string]*dbJSONSchemaResponse{"alias-tableName": {Fields: []*model.FieldType{
				{
					FieldName: "abcd",
				},
			}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModules := mockModulesInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			_, got, err := tt.s.GetSchemas(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.format, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSchemas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetSchemas() = %v, want %v", got, tt.want)
			}
			mockModules.AssertExpectations(t)
			mockSchema.AssertExpectations(t)
		})
	}
}

// func TestManager_SetReloadSchema(t *testing.T) {
//
// 	project := "1"
// 	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"create": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
// 	s := schema.Init(crud.Init())
// 	if err := s.SetConfig(rule, project); err != nil {
// 		t.Errorf("error setting config of schema - %s", err.Error())
// 	}
// 	type mockArgs struct {
// 		method         string
// 		args           []interface{}
// 		paramsReturned []interface{}
// 	}
// 	type args struct {
// 		ctx       context.Context
// 		dbAlias   string
// 		project   string
// 		schemaArg *schema.Schema
// 	}
// 	tests := []struct {
// 		name            string
// 		s               *Manager
// 		args            args
// 		modulesMockArgs []mockArgs
// 		storeMockArgs   []mockArgs
// 		want            map[string]interface{}
// 		wantErr         bool
// 	}{
// 		{
// 			name:    "unable to get project",
// 			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
// 			args:    args{ctx: context.Background(), dbAlias: "mongo", project: "2"},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "database not present in config",
// 			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
// 			args:    args{ctx: context.Background(), dbAlias: "notMongo", project: "1"},
// 			wantErr: true,
// 		},
// 		{
// 			name: "colName is default and unable to set crud config",
// 			s:    &Manager{clusterID:"chicago",storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {}}}}}}}}},
// 			args: args{ctx: context.Background(), dbAlias: "mongo", project: "1"},
// 			modulesMockArgs: []mockArgs{
// 				{
// 					method:         "SetDatabaseConfig",
// 					args:           []interface{}{"1", mock.Anything},
// 					paramsReturned: []interface{}{errors.New("error setting db module config")},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "colName is default and unable to set project",
// 			s:    &Manager{clusterID:"chicago",storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {}}}}}}}}},
// 			args: args{ctx: context.Background(), dbAlias: "mongo", project: "1"},
// 			modulesMockArgs: []mockArgs{
// 				{
// 					method:         "SetDatabaseConfig",
// 					args:           []interface{}{"1", mock.Anything},
// 					paramsReturned: []interface{}{nil},
// 				},
// 			},
// 			storeMockArgs: []mockArgs{
// 				{
// 					method:         "SetResource",
// 					args:           []interface{}{context.Background(), mock.Anything},
// 					paramsReturned: []interface{}{errors.New("Invalid config file type")},
// 				},
// 			},
// 			want:    map[string]interface{}{},
// 			wantErr: true,
// 		},
// 		{
// 			name: "colName is default and project is set",
// 			s:    &Manager{clusterID:"chicago",storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {}}}}}}}}},
// 			args: args{ctx: context.Background(), dbAlias: "mongo", project: "1"},
// 			modulesMockArgs: []mockArgs{
// 				{
// 					method:         "SetDatabaseConfig",
// 					args:           []interface{}{"1", mock.Anything},
// 					paramsReturned: []interface{}{nil},
// 				},
// 			},
// 			storeMockArgs: []mockArgs{
// 				{
// 					method:         "SetResource",
// 					args:           []interface{}{context.Background(), mock.Anything},
// 					paramsReturned: []interface{}{nil},
// 				},
// 			},
// 			want: map[string]interface{}{},
// 		},
// 		{
// 			name:    "unable to inspect schema",
// 			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Crud: config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}}}}}}},
// 			args:    args{ctx: context.Background(), dbAlias: "mongo", project: "1", schemaArg: s},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
//
// 			mockModules := mockModulesInterface{}
// 			mockStore := mockStoreInterface{}
//
// 			for _, m := range tt.modulesMockArgs {
// 				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
// 			}
// 			for _, m := range tt.storeMockArgs {
// 				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
// 			}
//
// 			tt.s.modules = &mockModules
// 			tt.s.store = &mockStore
//
// 			got, err := tt.s.SetReloadSchema(context.Background(), tt.args.dbAlias, tt.args.project)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Manager.SetReloadSchema() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Manager.SetReloadSchema() = %v, want %v", got, tt.want)
// 			}
//
// 			mockModules.AssertExpectations(t)
// 			mockStore.AssertExpectations(t)
// 		})
// 	}
// }

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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{"": &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{"": &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1", schema: "type event {id: ID! title: String}"},
			wantErr: true,
		},
		{
			name: "collections nil and unable to set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"): &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collections nil and unable to set project",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"): &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"), &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collections nil and project is set",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"): &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"), &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "collection not present and unable to set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"): &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection not present and unable to set project",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"): &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"), &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection not present and project is set",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"): &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}, config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "notTableName"), &config.DatabaseSchema{Table: "notTableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "collection present and unable to set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection present and unable to set project",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "collection present and project is set",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1", schema: "type event {id: ID! title: String}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
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

			if _, err := tt.s.SetSchemaInspection(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.schema, model.RequestParams{}); (err != nil) != tt.wantErr {
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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "unable to set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{}},
					paramsReturned: []interface{}{errors.New("unable to set db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "schema inspection is removed",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
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

			if _, err := tt.s.RemoveCollection(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.col, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RemoveCollection() error = %v, wantErr %v", err, tt.wantErr)
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
		storeMockArgs       []mockArgs
		wantErr             bool
	}{
		{
			name:    "database not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1", projectConfig: &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"anotherTableName": {}}}},
			wantErr: true,
		},
		{
			name: "unable modify all schema",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"anotherTableName": {}}}},
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
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
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
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "collections are nil and schemas are applied",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "unable set crud config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "schemas are applied",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", projectConfig: &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}, v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), mock.Anything},
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
			for _, m := range tt.schemaErrorMockArgs {
				mockErrorSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if err := tt.s.applySchemas(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.projectConfig, tt.args.v); (err != nil) != tt.wantErr {
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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			wantErr: true,
		},
		{
			name:    "unable to apply schemas",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{"": &config.DatabaseRule{Table: "tableName", DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1", v: config.CrudStub{Collections: map[string]*config.TableRule{"tableName": {}}}},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "1", DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{DBName: "1", Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "SchemaModifyAll",
					args:           []interface{}{context.Background(), "alias", "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "modified all schema successfully",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DBName: "1", DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1", v: config.CrudStub{DBName: "1", Collections: map[string]*config.TableRule{"tableName": {Schema: "type event {id: ID! title: String}"}}}},
			modulesMockArgs: []mockArgs{
				{
					method:         "GetSchemaModuleForSyncMan",
					paramsReturned: []interface{}{&mockSchema},
				},
				{
					method:         "SetDatabaseSchemaConfig",
					args:           []interface{}{mock.Anything, "1", config.DatabaseSchemas{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"): &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{context.Background(), config.GenerateResourceID("chicago", "1", config.ResourceDatabaseSchema, "alias", "tableName"), &config.DatabaseSchema{Schema: "type event {id: ID! title: String}", Table: "tableName", DbAlias: "alias"}},
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

			if _, err := tt.s.SetModifyAllSchema(context.Background(), tt.args.dbAlias, tt.args.project, tt.args.v, model.RequestParams{}); (err != nil) != tt.wantErr {
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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{Conn: "mongo:conn", Enabled: true, DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "db alias not present in config",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{Conn: "mongo:conn", Type: "mongo", Enabled: true, DbAlias: "alias"}}}}}},
			args:    args{ctx: context.Background(), dbAlias: "notAlias", project: "1"},
			wantErr: true,
		},
		{
			name: "got db alias config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{Conn: "mongo:conn", Type: "mongo", Enabled: true, DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "alias", project: "1"},
			want: []interface{}{config.Crud{"alias": {Enabled: true, Conn: "mongo:conn", Type: "mongo"}}},
		},
		{
			name: "got services config",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{Conn: "mongo:conn", Type: "mongo", Enabled: true, DbAlias: "alias"}}}}}},
			args: args{ctx: context.Background(), dbAlias: "*", project: "1"},
			want: []interface{}{config.Crud{"alias": {Enabled: true, Conn: "mongo:conn", Type: "mongo"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetDatabaseConfig(context.Background(), tt.args.project, tt.args.dbAlias, model.RequestParams{})
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
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{"": &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args:    args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "2"},
			wantErr: true,
		},
		{
			name:    "specified collection not present in config for dbAlias",
			s:       &Manager{storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{"resourceId": &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseSchemas: config.DatabaseSchemas{"": &config.DatabaseSchema{Table: "tableName", DbAlias: "alias", Schema: "type event {id: ID! title: String}"}}}}}},
			args:    args{ctx: context.Background(), col: "notTableName", dbAlias: "alias", project: "1"},
			wantErr: true,
		},
		{
			name: "got collection rules for specific db alias and collection",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", IsRealTimeEnabled: true, DbAlias: "alias", Rules: map[string]*config.Rule{"create": &config.Rule{ID: "id"}}}}}}}},
			args: args{ctx: context.Background(), col: "tableName", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbRulesResponse{"alias-tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"create": {ID: "id"}}}}},
		},
		{
			name: "col is empty and got collection rules",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", IsRealTimeEnabled: true, DbAlias: "alias", Rules: map[string]*config.Rule{"create": &config.Rule{ID: "id"}}}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "alias", project: "1"},
			want: []interface{}{map[string]*dbRulesResponse{"alias-tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"create": {ID: "id"}}}}},
		},
		{
			name: "col and dbalias is empty and got collection rules",
			s:    &Manager{clusterID: "chicago", storeType: "local", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, DatabaseConfigs: config.DatabaseConfigs{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseConfig, "alias"): &config.DatabaseConfig{DbAlias: "alias"}}, DatabaseRules: config.DatabaseRules{config.GenerateResourceID("chicago", "1", config.ResourceDatabaseRule, "alias", "tableName", "rule"): &config.DatabaseRule{Table: "tableName", IsRealTimeEnabled: true, DbAlias: "alias", Rules: map[string]*config.Rule{"create": &config.Rule{ID: "id"}}}}}}}},
			args: args{ctx: context.Background(), col: "*", dbAlias: "*", project: "1"},
			want: []interface{}{map[string]*dbRulesResponse{"alias-tableName": {IsRealTimeEnabled: true, Rules: map[string]*config.Rule{"create": {ID: "id"}}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetCollectionRules(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.col, model.RequestParams{})
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
