package schema

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func Test_generateInspection(t *testing.T) {
	type args struct {
		dbType      string
		col         string
		fields      []utils.FieldType
		foreignkeys []utils.ForeignKeysType
		indexkeys   []utils.IndexType
	}
	tests := []struct {
		name    string
		args    args
		want    model.Collection
		wantErr bool
	}{
		// TODO TEST CASES REMAINING FOR
		// Detecting external index & normal index
		// Detecting external foreign keys & normal index
		{
			name: "identify varchar with any size",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "varchar(5550)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "foreign keys with constraint name not matching gateways convention name",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []utils.FieldType{{FieldName: "col1", FieldType: "varchar(5550)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{{
					TableName:      "table1",
					ColumnName:     "col1",
					ConstraintName: "some-random-name",
					DeleteRule:     "NO_ACTION",
					RefTableName:   "table2",
					RefColumnName:  "id",
				}},
			},
			want: model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsForeign: true, IsFieldTypeRequired: true, Kind: model.TypeID, IsPrimary: true, JointTable: &model.TableProperties{
				To:             "id",
				Table:          "table2",
				OnDelete:       "NO_ACTION",
				ConstraintName: "some-random-name",
			}}}},
			wantErr: false,
		},
		{
			name: "primary-!null-ID",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Integer",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "bigint", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Integer"}}},
			wantErr: false,
		},
		{
			name: "unique-!null-String",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "text", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "String"}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Boolean",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "boolean", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Boolean"}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-Float",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "float", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Float", IsForeign: true, JointTable: &model.TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-DateTime",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "datetime", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "DateTime", IsForeign: true, JointTable: &model.TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-wrongDataType",
			args: args{
				dbType:      "mysql",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "wrongType", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			wantErr: true,
		},
		// postgres
		{
			name: "JSON with not null",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "jsonb", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "JSON"}}},
			wantErr: false,
		},
		{
			name: "default key -!null-ID",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "float", FieldNull: "NO", FieldDefault: "9.8"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Float", IsDefault: true, Default: "9.8"}}},
			wantErr: false,
		},
		{
			name: "default key string -!null-ID",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "text", FieldNull: "NO", FieldDefault: "'string'::text"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "String", IsDefault: true, Default: "\"string\""}}},
			wantErr: false,
		},
		{
			name: "default key boolean -!null-ID",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "boolean", FieldNull: "NO", FieldDefault: "true"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Boolean", IsDefault: true, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "primary-!null-ID",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "character varying", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Integer",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "bigint", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Integer"}}},
			wantErr: false,
		},
		{
			name: "unique-!null-String",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "text", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "String"}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Boolean",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "boolean", FieldNull: "NO"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Boolean"}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-Float",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "float", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Float", IsForeign: true, JointTable: &model.TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-DateTime",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "datetime", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "DateTime", IsForeign: true, JointTable: &model.TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-wrongDataType",
			args: args{
				dbType:      "postgres",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "wrongType", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			wantErr: true,
		},
		// sql server
		{
			name: "identify type id with any size in varchar",
			args: args{
				dbType:      "sqlserver",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "varchar(5520)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "sqlserver type string",
			args: args{
				dbType:      "sqlserver",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "varchar(-1)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: model.TypeString, IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "primary-!null-ID",
			args: args{
				dbType:      "sqlserver",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "default key string -!null-ID",
			args: args{
				dbType:      "sqlserver",
				col:         "table1",
				fields:      []utils.FieldType{{FieldName: "col1", FieldType: "text", FieldNull: "NO", FieldDefault: "((string))"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    model.Collection{"table1": model.Fields{"col1": &model.FieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "String", IsDefault: true, Default: "\"string\""}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateInspection(tt.args.dbType, tt.args.col, tt.args.fields, tt.args.foreignkeys, tt.args.indexkeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateInspection() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateInspection() = %s, want %s", print(got), print(tt.want))
			}
		})
	}
}

func print(val interface{}) string {
	b, _ := json.MarshalIndent(val, "", "  ")
	return string(b)
}
