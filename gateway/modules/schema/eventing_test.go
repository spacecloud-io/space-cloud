package schema

import (
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

func TestSchema_CheckIfEventingIsPossible(t *testing.T) {
	type fields struct {
		lock      sync.RWMutex
		SchemaDoc model.Type
		crud      model.CrudSchemaInterface
		project   string
		config    config.Crud
	}
	type args struct {
		dbAlias string
		col     string
		obj     map[string]interface{}
		isFind  bool
	}
	crudPostgres := crud.Init()
	_ = crudPostgres.SetConfig("test", config.Crud{"postgres": {Type: "sql-postgres", Enabled: false}})

	crudMySQL := crud.Init()
	_ = crudMySQL.SetConfig("test", config.Crud{"mysql": {Type: "sql-mysql", Enabled: false}})

	crudSQLServer := crud.Init()
	_ = crudSQLServer.SetConfig("test", config.Crud{"sqlserver": {Type: "sql-sqlserver", Enabled: false}})
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantFindForUpdate map[string]interface{}
		wantPresent       bool
	}{
		// TODO: Add test cases.
		{
			name: "dbAlias not provided",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{"col2": "xyz"},
				isFind:  false,
			},
			fields:            fields{SchemaDoc: model.Type{}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{},
			wantPresent:       false,
		},
		{
			name: "Col not provided",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{"col2": "xyz"},
				isFind:  false,
			},
			fields:            fields{SchemaDoc: model.Type{"mysql": model.Collection{}}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{},
			wantPresent:       false,
		},
		{
			name: "fieldSchema with IsIndex and IsUnique",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{"col2": "xyz"},
				isFind:  false,
			},
			fields:            fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON, IsIndex: true, IsUnique: true, IndexInfo: &model.TableProperties{Group: "abcd"}}}}}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{"col2": "xyz"},
			wantPresent:       true,
		},
		{
			name: "fieldSchema with IsIndex and IsUnique, obj map[string]interface{} ,isFind true",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{"col2": map[string]interface{}{"$eq": "xyz"}},
				isFind:  true,
			},
			fields:            fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON, IsIndex: true, IsUnique: true, IndexInfo: &model.TableProperties{Group: "abcd"}}}}}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{"col2": "xyz"},
			wantPresent:       true,
		},
		{
			name: "fieldSchema with IsIndex and IsUnique, obj map[string]interface{} and isFind false",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{"col2": map[string]interface{}{"$eq": 10}},
				isFind:  false,
			},
			fields:            fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON, IsIndex: true, IsUnique: true, IndexInfo: &model.TableProperties{Group: "abcd"}}}}}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{"col2": map[string]interface{}{"$eq": 10}},
			wantPresent:       true,
		},
		{
			name: "fieldSchema with IsIndex and IsUnique obj not provided",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{},
				isFind:  false,
			},
			fields:            fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON, IsIndex: true, IsUnique: true, IndexInfo: &model.TableProperties{Group: "abcd"}}}}}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{},
			wantPresent:       false,
		},
		{
			name: "fieldSchema with IsPrimary",
			args: args{
				dbAlias: "mysql",
				col:     "table1",
				obj:     map[string]interface{}{"col2": "xyz"},
				isFind:  false,
			},
			fields:            fields{SchemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeJSON, IsPrimary: true}}}}, crud: crudMySQL, project: "test"},
			wantFindForUpdate: map[string]interface{}{"col2": "xyz"},
			wantPresent:       true,
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
			gotFindForUpdate, gotPresent := s.CheckIfEventingIsPossible(tt.args.dbAlias, tt.args.col, tt.args.obj, tt.args.isFind)
			if !reflect.DeepEqual(len(gotFindForUpdate), len(tt.wantFindForUpdate)) {
				t.Errorf("Schema.CheckIfEventingIsPossible() gotFindForUpdate = %v, want %v", gotFindForUpdate, tt.wantFindForUpdate)
			} else if len(gotFindForUpdate) != 0 {
				if !reflect.DeepEqual(gotFindForUpdate, tt.wantFindForUpdate) {
					t.Errorf("Schema.CheckIfEventingIsPossible() gotFindForUpdate = %v, want %v", gotFindForUpdate, tt.wantFindForUpdate)
				}
			}
			if gotPresent != tt.wantPresent {
				t.Errorf("Schema.CheckIfEventingIsPossible() gotPresent = %v, want %v", gotPresent, tt.wantPresent)
			}
		})
	}
}
