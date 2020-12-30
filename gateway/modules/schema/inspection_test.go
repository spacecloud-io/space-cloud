package schema

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func Test_generateInspection(t *testing.T) {
	var firstColumn = "column1"
	var secondColumn = "column2"
	type args struct {
		dbType    string
		col       string
		fields    []model.InspectorFieldType
		indexKeys []model.IndexType
	}
	tests := []struct {
		name    string
		args    args
		want    model.Collection
		wantErr bool
	}{
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
		{
			name: "MySQL field col1 with type ID",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "YES", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeJSON}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "YES"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Unsupported type",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "wrongType", FieldNull: "YES"}},
			},
			wantErr: true,
		},
		{
			name: "MySQL field col1 which is not null with type ID ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type String ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Boolean ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Integer ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Float ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type DateTime ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type JSON ",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON}}},
			wantErr: false,
		},
		// There is a bug in code, inspection cannot detect @createdAt,@updatedAt directives
		// TODO: What other special directives do we have ?
		// {
		// 	name: "MySQL field col1 which is not null with type DateTime having directive @createdAt",
		// 	args: args{
		// 		dbAlias:      "mysql",
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},
		// 	want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IsCreatedAt: true}}},
		// 	wantErr: false,
		// },
		// {
		// 	name: "MySQL field col1 which is not null with type DateTime having directive @updatedAt",
		// 	args: args{
		// 		dbAlias:      "mysql",
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},
		// 	want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IsUpdatedAt: true}}},
		// 	wantErr: false,
		// },
		{
			name: "MySQL field col1 which is not null with type ID having default value INDIA",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeID, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type String having default value INDIA",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeString, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Boolean having default value true",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO", FieldDefault: "true", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBoolean, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Integer having default value 100",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type Float having default value 9.8",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", FieldDefault: "9.8", AutoIncrement: "false", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Default: "9.8"}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeDateTime, Default: "\"2020-05-30T00:42:05+00:00\""}}},
			wantErr: false,
		},
		{
			name: `MySQL field col1 which is not null with type JSON having default value {"id":"zerfvnex","name":"john"}`,
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO", FieldDefault: `{"id":"zerfvnex","name":"john"}`, AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeJSON, Default: `{"id":"zerfvnex","name":"john"}`}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having primary key constraint",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IsForeign: true, JointTable: &model.TableProperties{To: "col2", Table: "table2", ConstraintName: getConstraintName("table1", firstColumn), OnDelete: "NO_ACTION"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: getConstraintName("table1", firstColumn), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: getConstraintName("table1", firstColumn), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "decimal", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5, TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: getConstraintName("table1", firstColumn), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: getConstraintName("table1", firstColumn), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having foreign key constraint created through or not from space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO", TableName: "table1", RefTableName: "table2", RefColumnName: "col2", ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO_ACTION"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IsForeign: true, JointTable: &model.TableProperties{To: "col2", ConstraintName: getConstraintName("table1", firstColumn), OnDelete: "NO_ACTION", Table: "table2"}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single unique index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: secondColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single unique index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: secondColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: true},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: true},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsUnique: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single index constraint created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: secondColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Args: &model.FieldArgs{Precision: 10, Scale: 5}, Kind: model.TypeFloat, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "index1", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: getIndexName("table1", "index1")}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON which is not null having single index constraint not created through space cloud",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}}}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "text", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "bigint", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}, {ColumnName: secondColumn, FieldType: "float", FieldNull: "NO", NumericPrecision: 10, NumericScale: 5}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, Args: &model.FieldArgs{Precision: 10, Scale: 5}, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "boolean", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type DateTime, col2 with type DateTime which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "datetime", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint not created through space cloud",
			args: args{
				dbType: "mysql",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "json", FieldNull: "NO"}, {ColumnName: secondColumn, FieldType: "json", FieldNull: "NO"}},
				indexKeys: []model.IndexType{
					{TableName: "table1", ColumnName: firstColumn, IndexName: "custom-index", Order: 1, Sort: model.DefaultIndexSort, IsUnique: false},
					{TableName: "table1", ColumnName: secondColumn, IndexName: "custom-index", Order: 2, Sort: model.DefaultIndexSort, IsUnique: false},
				},
			},
			want: model.Collection{"table1": model.Fields{
				firstColumn:  &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 1, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
				secondColumn: &model.FieldType{FieldName: secondColumn, IsFieldTypeRequired: true, Kind: model.TypeJSON, IndexInfo: []*model.TableProperties{{IsIndex: true, Group: "custom-index", Order: 2, Sort: model.DefaultIndexSort, ConstraintName: "custom-index"}}},
			}},
			wantErr: false,
		},
		{
			name: "identify varchar with any size",
			args: args{
				dbType:    "mysql",
				col:       "table1",
				fields:    []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(5550)", FieldNull: "NO"}},
				indexKeys: []model.IndexType{{IsPrimary: true, ColumnName: firstColumn}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: "ID", IsPrimary: true}}},
			wantErr: false,
		},
		// postgres
		{
			name: "Postgres field col1 with type ID which is not null having default value INDIA",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "character varying", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeID, IsDefault: true, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type String which is not null having default value INDIA",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "text", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeString, IsDefault: true, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Integer which is not null having default value 100",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeInteger, IsDefault: true, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Float which is not null having default value 9.8",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", FieldDefault: "9.8", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeFloat, IsDefault: true, Default: "9.8"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Boolean which is not null having default value true",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO", FieldDefault: "true", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeBoolean, IsDefault: true, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type DateTime which is not null having default value 2020-05-30T00:42:05+00:00",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "timestamp", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, Kind: model.TypeDateTime, IsDefault: true, Default: "\"2020-05-30T00:42:05+00:00\""}}},
			wantErr: false,
		},
		{
			name: "Postgres field col1 with type Unsupported type",
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "some-type", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00"}},
			},
			wantErr: true,
		},
		{
			name: `Postgres field col1 which is not null with type JSON having default value {"id":"zerfvnex","name":"john"}`,
			args: args{
				dbType: "postgres",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "jsonb", FieldNull: "NO", FieldDefault: `{"id":"zerfvnex","name":"john"}`, AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeJSON, Default: `{"id":"zerfvnex","name":"john"}`}}},
			wantErr: false,
		},
		// sql server
		{
			name: "SQL-Server field col1 which is not null with type ID having default value INDIA",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeID, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type String having default value INDIA",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "varchar(-1)", FieldNull: "NO", FieldDefault: "INDIA", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeString, Default: "\"INDIA\""}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Boolean having default value true",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO", FieldDefault: "1", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBoolean, Default: "true"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Boolean having default value false",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "boolean", FieldNull: "NO", FieldDefault: "0", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeBoolean, Default: "false"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Integer having default value 100",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "bigint", FieldNull: "NO", FieldDefault: "100", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeInteger, Default: "100"}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type Float having default value 9.8",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "float", FieldNull: "NO", FieldDefault: "9.8", AutoIncrement: "false", NumericPrecision: 10, NumericScale: 5}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeFloat, Default: "9.8", Args: &model.FieldArgs{Precision: 10, Scale: 5}}}},
			wantErr: false,
		},
		{
			name: "SQL-Server field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			args: args{
				dbType: "sqlserver",
				col:    "table1",
				fields: []model.InspectorFieldType{{ColumnName: firstColumn, FieldType: "datetime", FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", AutoIncrement: "false"}},
			},
			want:    model.Collection{"table1": model.Fields{firstColumn: &model.FieldType{FieldName: firstColumn, IsFieldTypeRequired: true, IsDefault: true, Kind: model.TypeDateTime, Default: "\"2020-05-30T00:42:05+00:00\""}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
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
