package schema

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema/helpers"
)

func Test_generateInspection(t *testing.T) {
	type args struct {
		dbType    string
		col       string
		fields    []model.InspectorFieldType
		indexKeys []model.IndexType
	}

	type testGenerateInspection struct {
		name    string
		args    args
		want    model.Collection
		wantErr bool
	}

	var checkColumnType = []testGenerateInspection{
		// Mysql
		{
			name: "MySQL field col1 with type ID",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer (mediumint)",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "mediumint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer (smallint)",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeSmallInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer (bigint)",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeBigInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "YES", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeJSON}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Unsupported type",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "wrongType", FieldNull: "YES"}},
			},
			wantErr: true,
		},
		// Postgres
		{
			name: "Postgres field col1 with type ID",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "character varying", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type String",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Boolean",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Integer",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "integer", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Integer (smallint)",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeSmallInteger}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Integer (bigint)",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeBigInteger}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Float",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "YES", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type JSON",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeJSON}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type DateTime",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Unsupported type",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "wrongType", FieldNull: "YES"}},
			},
			wantErr: true,
		},
		// 	Sql server
		{
			name: "SQL-Server field col1 with type ID",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type String",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type Boolean",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type Integer",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type Integer (smallint)",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeSmallInteger}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type Integer (bigint)",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeBigInteger}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type Float",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "YES", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type JSON",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeJSON}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type DateTime",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 with type Unsupported type",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "wrongType", FieldNull: "YES"}},
			},
			wantErr: true,
		},
	}

	var checkColumnTypeWithNotNull = []testGenerateInspection{
		// Mysql
		{
			name: "MySQL field col1 which is not null with type ID ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type String ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Boolean ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type int ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type mediumint ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "mediumint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type smallint ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeSmallInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type bigint ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBigInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Float ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type DateTime ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type JSON ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON}}},
			wantErr: false,
		},
		// Postgres
		{
			name: "Postgres field col1 which is not null with type ID ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "character varying", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type String ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type Boolean ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type integer ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "integer", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type smallint ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeSmallInteger}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type bigint ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBigInteger}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type Float ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type DateTime ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type JSON ",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON}}},
			wantErr: false,
		},
		// SQL server
		{
			name: "SQL-server field col1 which is not null with type ID ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type String ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type Boolean ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type int ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type smallint ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeSmallInteger}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type bigint ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBigInteger}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type Float ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type DateTime ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "SQL-server field col1 which is not null with type JSON ",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON}}},
			wantErr: false,
		},
	}

	var defaultTestCases = []testGenerateInspection{
		// Mysql
		{
			name: "MySQL field col1 which is not null with type ID having default value INDIA",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeID, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type String having default value INDIA",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeString, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Boolean having default value true",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO", FieldDefault: "true", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBoolean, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type int having default value 100",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type mediumint having default value 100",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "mediumint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type smallint having default value 100",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeSmallInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type bigint having default value 100",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBigInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Float having default value 9.8",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", FieldDefault: "9.8", AutoIncrement: "false", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Default: "9.8"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeDateTime, Default: "\"2020-05-30T00:42:05+00:00\""}}},
			wantErr: false,
		},
		{
			name: `MySQL field col1 which is not null with type JSON having default value {"id":"zerfvnex","name":"john"}`,
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO", FieldDefault: `{"id":"zerfvnex","name":"john"}`, AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeJSON, Default: `{"id":"zerfvnex","name":"john"}`}}},
			wantErr: false,
		},
		// postgres
		{
			name: "Postgres field col1 with type ID which is not null having default value INDIA",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "character varying", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IsDefault: true, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type String which is not null having default value INDIA",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IsDefault: true, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type integer having default value 100",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "integer", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type smallint having default value 100",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeSmallInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 which is not null with type bigint having default value 100",
			args: args{
				dbType: string(model.Postgres),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBigInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Float which is not null having default value 9.8",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", FieldDefault: "9.8", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, IsDefault: true, Default: "9.8"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Boolean which is not null having default value true",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO", FieldDefault: "true", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IsDefault: true, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type DateTime which is not null having default value 2020-05-30T00:42:05+00:00",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "timestamp", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IsDefault: true, Default: "\"2020-05-30T00:42:05+00:00\""}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Unsupported type",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "some-type", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00"}},
			},
			wantErr: true,
		},
		{
			name: `Postgres field col1 which is not null with type JSON having default value {"id":"zerfvnex","name":"john"}`,
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "jsonb", FieldNull: "NO", FieldDefault: `{"id":"zerfvnex","name":"john"}`, AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeJSON, Default: `{"id":"zerfvnex","name":"john"}`}}},
			wantErr: false,
		},
		// sql server
		{
			name: "SQL-Server field col1 which is not null with type ID having default value INDIA",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeID, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type String having default value INDIA",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(-1)", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeString, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Boolean having default value true",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO", FieldDefault: "1", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBoolean, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Boolean having default value false",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO", FieldDefault: "0", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBoolean, Default: "false"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type int having default value 100",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type smallint having default value 100",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "smallint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeSmallInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type bigint having default value 100",
			args: args{
				dbType: string(model.SQLServer),
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBigInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Float having default value 9.8",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", FieldDefault: "9.8", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeFloat, Default: "9.8"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeDateTime, Default: "\"2020-05-30T00:42:05+00:00\""}}},
			wantErr: false,
		},
	}

	var foreignKeyTestCases = []testGenerateInspection{
		{
			name: "MySQL field col1 with type ID which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IsForeign: true, JointTable: &model.TableProperties{To: "col2", Table: "table2", ConstraintName: helpers.GetConstraintName("table1", "column1"), OnDelete: "NO_ACTION"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5, TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: helpers.GetConstraintName("table1", "column1"), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
	}

	var uniqueKeyTestCases = []testGenerateInspection{
		{
			name: "MySQL field col1 with type ID which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: "column2", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: "column2", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column2", IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
	}

	var indexKeyTestCases = []testGenerateInspection{
		{
			name: "MySQL field col1 with type ID which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: "column2", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "text", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "int", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "int", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: "column2", FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "boolean", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "datetime", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "json", FieldNull: "NO"}, {ColumnName: "column2", FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: "column1", IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: "column2", IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column1", IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				"column2": &model.FieldType{FieldName: "column2", IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{Field: "column2", IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
	}

	var primaryKeyTestCases = []testGenerateInspection{
		{
			name: "MySQL field col1 with type ID which is not null having primary key constraint",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(50)", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: model.TypeID}}},
			wantErr: false,
		},
	}

	var miscellaneousTestCases = []testGenerateInspection{
		{
			name: "identify varchar with any size",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: "column1", FieldType: "varchar(5550)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{IsPrimary: true, ColumnName: "column1", Order: 1}},
			},
			want:    model.Collection{"table1": model.Fields{"column1": &model.FieldType{FieldName: "column1", IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true, PrimaryKeyInfo: &model.TableProperties{Order: 1}}}},
			wantErr: false,
		},
	}

	// Test cases for each database follows a pattern of
	// 1) Checking each individual column type
	// 2) Checking each individual column type with not null
	// 3) Checking each individual column type with specific directives e.g -> @createdAt...
	// 4) Checking each individual column type with default value
	// 5) Type ID having primary key which is not null
	// 6) Individual columns having External & Internal foreign key which is not null
	// 7) Individual & Multiple columns having External & normal index key which is not null
	// 8) Individual & Multiple columns having External & normal unique index key which is not null
	// 9) Miscellaneous

	testCases := make([]testGenerateInspection, 0)
	testCases = append(testCases, checkColumnType...)
	testCases = append(testCases, checkColumnTypeWithNotNull...)
	testCases = append(testCases, defaultTestCases...)
	testCases = append(testCases, primaryKeyTestCases...)
	testCases = append(testCases, foreignKeyTestCases...)
	testCases = append(testCases, uniqueKeyTestCases...)
	testCases = append(testCases, indexKeyTestCases...)
	testCases = append(testCases, miscellaneousTestCases...)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateInspection(tt.args.dbType, tt.args.col, tt.args.fields, tt.args.indexKeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateInspection() error = %v, wantErr %v", err, tt.wantErr)
			}
			if arr := deep.Equal(got, tt.want); len(arr) > 0 {
				t.Errorf("generateInspection() differences = %v", arr)
			}
		})
	}
}
