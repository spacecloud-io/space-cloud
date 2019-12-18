package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/utils"
)

func Test_generateInspection(t *testing.T) {
	type args struct {
		dbType      string
		col         string
		fields      []utils.FieldType
		foreignkeys []utils.ForeignKeysType
	}
	tests := []struct {
		name    string
		args    args
		want    schemaCollection
		wantErr bool
	}{
		//TODO: Add test cases.
		{
			name: "primary-!null-ID",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Integer",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "bigint", FieldNull: "NO", FieldKey: "UNI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Integer", IsUnique: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-String",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "text", FieldNull: "NO", FieldKey: "UNI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "String", IsUnique: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Boolean",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "boolean", FieldNull: "NO", FieldKey: "UNI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Boolean", IsUnique: true}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-Float",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "float", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Float", IsForeign: true, JointTable: &TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-DateTime",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "datetime", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "DateTime", IsForeign: true, JointTable: &TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-wrongDataType",
			args: args{
				dbType:      "sql-mysql",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "wrongType", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			wantErr: true,
		},
		//postgres
		{
			name: "primary-!null-ID",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "character varying(50)", FieldNull: "NO", FieldKey: "PRI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Integer",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "bigint", FieldNull: "NO", FieldKey: "UNI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Integer", IsUnique: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-String",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "text", FieldNull: "NO", FieldKey: "UNI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "String", IsUnique: true}}},
			wantErr: false,
		},
		{
			name: "unique-!null-Boolean",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "boolean", FieldNull: "NO", FieldKey: "UNI"}},
				foreignkeys: []utils.ForeignKeysType{},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Boolean", IsUnique: true}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-Float",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "float", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "Float", IsForeign: true, JointTable: &TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-DateTime",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "datetime", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			want:    schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", IsFieldTypeRequired: true, Kind: "DateTime", IsForeign: true, JointTable: &TableProperties{To: "col2", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "foreign-!null-wrongDataType",
			args: args{
				dbType:      "sql-postgres",
				col:         "table1",
				fields:      []utils.FieldType{utils.FieldType{FieldName: "col1", FieldType: "wrongType", FieldNull: "NO", FieldKey: "MUL"}},
				foreignkeys: []utils.ForeignKeysType{utils.ForeignKeysType{TableName: "table1", ColumnName: "col1", RefTableName: "table2", RefColumnName: "col2"}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateInspection(tt.args.dbType, tt.args.col, tt.args.fields, tt.args.foreignkeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateInspection() error = %v, wantErr %v", err, tt.wantErr)
				b, err1 := json.MarshalIndent(got, "", "  ")
				if err1 != nil {
					fmt.Println("error:", err1)
				}
				fmt.Print(string(b))
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateInspection() = %v, want %v", got, tt.want)
				b, err1 := json.MarshalIndent(got, "", "  ")
				if err1 != nil {
					fmt.Println("error:", err1)
				}
				fmt.Print(tt.name, string(b))
				return
			}
		})
	}
}
