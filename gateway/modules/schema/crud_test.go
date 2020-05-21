package schema

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSchema_CrudPostProcess(t *testing.T) {
	b, err := json.Marshal(model.ReadRequest{Operation: "hello"})
	if err != nil {
		logrus.Errorf("err=%v", err)
	}
	type fields struct {
		lock      sync.RWMutex
		SchemaDoc model.Type
		crud      model.CrudSchemaInterface
		project   string
		config    config.Crud
	}
	type args struct {
		ctx     context.Context
		dbAlias string
		col     string
		result  interface{}
	}
	crudPostgres := crud.Init()
	_ = crudPostgres.SetConfig("test", config.Crud{"postgres": {Type: "sql-postgres", Enabled: false}})

	crudMySQL := crud.Init()
	_ = crudMySQL.SetConfig("test", config.Crud{"mysql": {Type: "sql-mysql", Enabled: false}})

	crudSQLServer := crud.Init()
	_ = crudSQLServer.SetConfig("test", config.Crud{"sqlserver": {Type: "sql-sqlserver", Enabled: false}})
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Empty SchemaDoc provided",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  []interface{}{},
			},
			fields:  fields{crud: crudMySQL, project: "test"},
			wantErr: true,
		},
		{
			name: "Empty col Provided",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  []interface{}{},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
		{
			name: "Empty result provided",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  []interface{}{},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON}}}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
		{
			name: "Checking with v as interface",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  []interface{}{map[string]interface{}{}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON}}}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
		{
			name: "Checking with v as mapstring",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  map[string]interface{}{"col2": b},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON}}}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
		{
			name: "Unable to assert interface to []byte",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  map[string]interface{}{"col2": "mock"},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON}}}}, crud: crudMySQL, project: "test"},
			wantErr: true,
		},
		{
			name: "Unable to unmarshal",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  map[string]interface{}{"col2": []byte("mock")},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON}}}}, crud: crudMySQL, project: "test"},
			wantErr: true,
		},
		{
			name: "Set kind as datetime",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  []interface{}{map[string]interface{}{"col2": time.Now()}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
		{
			name: "Set kind as primitive.datetime",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				result:  []interface{}{map[string]interface{}{"col2": primitive.DateTime(1)}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				lock:      tt.fields.lock,
				SchemaDoc: tt.fields.SchemaDoc,
				crud:      tt.fields.crud,
				project:   tt.fields.project,
				config:    tt.fields.config,
			}
			if err := s.CrudPostProcess(tt.args.ctx, tt.args.dbAlias, tt.args.col, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("Schema.CrudPostProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchema_AdjustWhereClause(t *testing.T) {
	type fields struct {
		lock      sync.RWMutex
		SchemaDoc model.Type
		crud      model.CrudSchemaInterface
		project   string
		config    config.Crud
	}
	type args struct {
		dbAlias string
		dbType  utils.DBType
		col     string
		find    map[string]interface{}
	}
	crudPostgres := crud.Init()
	_ = crudPostgres.SetConfig("test", config.Crud{"postgres": {Type: "sql-postgres", Enabled: false}})

	crudMySQL := crud.Init()
	_ = crudMySQL.SetConfig("test", config.Crud{"mysql": {Type: "sql-mysql", Enabled: false}})

	crudSQLServer := crud.Init()
	_ = crudSQLServer.SetConfig("test", config.Crud{"sqlserver": {Type: "sql-sqlserver", Enabled: false}})
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{}}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
			wantErr: false,
		},
		{
			name: "param with map[string]interface{} having value time.time",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": map[string]interface{}{"time": time.Now()}},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
			wantErr: true,
		},
		{
			name: "Param as time.time",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": time.Now()},
			},
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
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
			fields:  fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}}, crud: crudMySQL, project: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				lock:      tt.fields.lock,
				SchemaDoc: tt.fields.SchemaDoc,
				crud:      tt.fields.crud,
				project:   tt.fields.project,
				config:    tt.fields.config,
			}
			if err := s.AdjustWhereClause(tt.args.dbAlias, tt.args.dbType, tt.args.col, tt.args.find); (err != nil) != tt.wantErr {
				t.Errorf("Schema.AdjustWhereClause() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
