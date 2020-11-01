package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

func TestSchema_CrudPostProcess(t *testing.T) {

	b, err := json.Marshal(model.ReadRequest{Operation: "hello"})
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.Background()), "err", err, nil)
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.Background()), "err", err, nil)
	}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type fields struct {
		SchemaDoc model.Type
	}
	type args struct {
		dbAlias string
		col     string
		result  interface{}
	}
	crudPostgres := crud.Init()
	_ = crudPostgres.SetConfig("test", config.DatabaseConfigs{config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "postgres"): &config.DatabaseConfig{DbAlias: "postgres", Type: "sql-postgres", Enabled: false}})

	crudMySQL := crud.Init()
	_ = crudMySQL.SetConfig("test", config.DatabaseConfigs{config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "mysql"): &config.DatabaseConfig{DbAlias: "mysql", Type: "sql-mysql", Enabled: false}})

	crudSQLServer := crud.Init()
	_ = crudSQLServer.SetConfig("test", config.DatabaseConfigs{config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "sqlserver"): &config.DatabaseConfig{DbAlias: "sqlserver", Type: "sql-sqlserver", Enabled: false}})
	tests := []struct {
		name         string
		fields       fields
		args         args
		crudMockArgs []mockArgs
		want         interface{}
		wantErr      bool
	}{
		// TODO: Add test cases for mongo
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc: tt.fields.SchemaDoc,
			}

			mockCrud := mockCrudSchemaInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			s.crud = &mockCrud

			err := s.CrudPostProcess(context.Background(), tt.args.dbAlias, tt.args.col, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.CrudPostProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.result, tt.want) {
				t.Errorf("Schema.CrudPostProcess() tt.args.result = %v, tt.want %v", tt.args.result, tt.want)
			}
		})
	}
}

func returntime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("invalid string format of datetime (%s)", s), map[string]interface{}{"error": err})
		return time.Now()
	}
	return t
}
func TestSchema_AdjustWhereClause(t *testing.T) {

	type fields struct {
		SchemaDoc model.Type
	}
	type args struct {
		dbAlias string
		dbType  model.DBType
		col     string
		find    map[string]interface{}
	}
	crudPostgres := crud.Init()
	_ = crudPostgres.SetConfig("test", config.DatabaseConfigs{config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "postgres"): &config.DatabaseConfig{DbAlias: "postgres", Type: "sql-postgres", Enabled: false}})

	crudMySQL := crud.Init()
	_ = crudMySQL.SetConfig("test", config.DatabaseConfigs{config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "mysql"): &config.DatabaseConfig{DbAlias: "mysql", Type: "sql-mysql", Enabled: false}})

	crudSQLServer := crud.Init()
	_ = crudSQLServer.SetConfig("test", config.DatabaseConfigs{config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseConfig, "sqlserver"): &config.DatabaseConfig{DbAlias: "sqlserver", Type: "sql-sqlserver", Enabled: false}})
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "db is not mongo",
			args: args{
				dbAlias: "mysql",
				dbType:  "sql",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "SchemaDoc not provided",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			},
			fields:  fields{},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "Col not provided",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{}}},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "Tableinfo not provided",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{}}}},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "Using param as string",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": returntime("2014-11-12T11:45:26.371Z")},
			wantErr: false,
		},
		{
			name: "Error string format provided",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12"},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": "2014-11-12"},
			wantErr: true,
		},
		{
			name: "param as map[string]interface{}",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": map[string]interface{}{"time": "2014-11-12T11:45:26.371Z"}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": returntime("2014-11-12T11:45:26.371Z")}},
			wantErr: false,
		},
		{
			name: "param with map[string]interface{} having value time.time",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": map[string]interface{}{"time": time.Now().Round(time.Second)}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": time.Now().Round(time.Second)}},
			wantErr: false,
		},
		{
			name: "Error foramt provided as value to map[string]interface{} ",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": map[string]interface{}{"time": "string"}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": "string"}},
			wantErr: true,
		},
		{
			name: "Param as time.time",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": time.Now().Round(time.Second)},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": time.Now().Round(time.Second)},
			wantErr: false,
		},
		{
			name: "Param as default",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": 10},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}},
			want:    map[string]interface{}{"col2": 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc: tt.fields.SchemaDoc,
			}
			err := s.AdjustWhereClause(context.Background(), tt.args.dbAlias, tt.args.dbType, tt.args.col, tt.args.find)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.AdjustWhereClause() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, tt.args.find) {
				t.Errorf("Schema.AdjustWhereClause() find = %v, want %v", tt.args.find, tt.want)
			}
		})
	}
}
