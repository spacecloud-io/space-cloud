// +build integration

package sql

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestSQL_DescribeTable(t *testing.T) {
	var firstColumn = "column1"
	var secondColumn = "column2"
	type args struct {
		ctx context.Context
		col string
	}
	type test struct {
		name           string
		createQuery    []string
		args           args
		fields         []model.InspectorFieldType
		foreignKeys    []model.ForeignKeysType
		indexKeys      []model.IndexType
		isMssqlSkip    bool
		isPostgresSkip bool
		isMysqlSkip    bool
		wantErr        bool
	}

	tests := []test{}

	mysqlTestCases := []test{
		{
			name:        "MySQL field col1 with type ID",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID having different size",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(500))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: "varchar(500)", FieldNull: "YES", VarcharSize: 500}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type String",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeString) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "YES", VarcharSize: 65535}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeBoolean) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeFloat) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type JSON",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeJSON) + " )"},

			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "MySQL field col1 with type DateTime",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		// Check Required fields(!)
		{
			name:        "MySQL field col1 which is not null with type ID ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type String ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeString) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: 65535}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Boolean ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeBoolean) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Integer ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Float ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeFloat) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type DateTime ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type JSON ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeJSON) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		// Currently @createdAt,@updateAt work at application level, so we cannot fetch data regarding this from our database
		// But this is possible if we do the above operation at database level
		// TODO: Do this when @createdAt & @updatedAt is pushed to database level
		// {
		// 	name: "MySQL field col1 which is not null with type DateTime having directive @createdAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// {
		// 	name: "MySQL field col1 which is not null with type DateTime having directive @updatedAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// Check @default directive
		// NOTE: JSON & text type cannot have default value
		{
			name:        "MySQL field col1 which is not null with type ID having default value INDIA",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) NOT NULL DEFAULT 'INDIA')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldDefault: "INDIA", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 boolean NOT NULL DEFAULT true)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldDefault: "true", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bit NOT NULL DEFAULT 1)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys:    []model.ForeignKeysType{},
			indexKeys:      []model.IndexType{},
			fields:         []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldDefault: "1", VarcharSize: model.SQLTypeIDSize}},
			wantErr:        false,
			isPostgresSkip: true,
			isMysqlSkip:    true,
		},
		{
			name:        "MySQL field col1 which is not null with type Integer having default value 100",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint NOT NULL DEFAULT 100)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", FieldDefault: "100", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type Float having default value 9.8",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 float NOT NULL DEFAULT 9.8)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", FieldDefault: "9.8", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldDefault: "2020-05-30 00:42:05", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys:    []model.ForeignKeysType{},
			indexKeys:      []model.IndexType{},
			fields:         []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", VarcharSize: model.SQLTypeIDSize}},
			wantErr:        false,
			isMysqlSkip:    true,
			isPostgresSkip: true,
		},
		// Check primary key constraint
		// NOTE: We are only checking for type ID as currently space cloud supports that only
		{
			name:        "MySQL field col1 with type ID which is not null having primary key constraint",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) PRIMARY KEY NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},

		// Check foreign key constraint
		// NOTE: We cannot create foreign keys on string
		{
			name:        "MySQL field col2 with type ID which is not null having foreign key constraint created through from space cloud",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeID) + " PRIMARY KEY NOT NULL)", "CREATE TABLE myproject.table2 (column2 " + getSQLType(*dbType, model.TypeID) + " NOT NULL, CONSTRAINT c_table1_column1 FOREIGN KEY (column2) REFERENCES table1(column1))"},
			args: args{
				ctx: context.Background(),
				col: "table2",
			},
			foreignKeys: []model.ForeignKeysType{{TableName: "table2", ColumnName: secondColumn, RefTableName: "table1", RefColumnName: firstColumn, ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO ACTION"}},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "MySQL field col2 with type ID which is not null having foreign key constraint with delete rule cascade created through from space cloud",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeID) + " PRIMARY KEY NOT NULL)", "CREATE TABLE myproject.table2 (column2 " + getSQLType(*dbType, model.TypeID) + " NOT NULL, CONSTRAINT c_table1_column1 FOREIGN KEY (column2) REFERENCES table1(column1) ON DELETE CASCADE )"},
			args: args{
				ctx: context.Background(),
				col: "table2",
			},
			foreignKeys: []model.ForeignKeysType{{TableName: "table2", ColumnName: secondColumn, RefTableName: "table1", RefColumnName: firstColumn, ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "CASCADE"}},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: secondColumn, FieldType: "varchar(50)", FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		// Check Unique Constraint
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type ID which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type Boolean which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type String which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: 65535}},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type Integer which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type Datetime which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type Float which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: I am getting type as primary
			name:        "MySQL field col1 with type JSON which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Id
		{
			name:        "MySQL field col1 with type ID, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type String
		{
			name:        "MySQL field col1 with type String, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Integer
		{
			name:        "MySQL field col1 with type Integer, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Float
		{
			name:        "MySQL field col1 with type Float, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Boolean
		{
			name:        "MySQL field col1 with type Boolean, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Datetime
		{
			name:        "MySQL field col1 with type Datetime, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "PRI", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type JSON
		{
			name:        "MySQL field col1 with type JSON, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},

		// Check Index Constraint
		{
			name:        "MySQL field col1 with type ID which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type String which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Integer which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type JSON which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s ON table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Id
		{
			name:        "MySQL field col1 with type ID, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type ID, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type String
		{
			name:        "MySQL field col1 with type String, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type String, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Integer
		{
			name:        "MySQL field col1 with type Integer, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Integer, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Float
		{
			name:        "MySQL field col1 with type Float, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Float, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Boolean
		{
			name:        "MySQL field col1 with type Boolean, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Boolean, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Datetime
		{
			name:        "MySQL field col1 with type Datetime, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldKey: "MUL", FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "MySQL field col1 with type Datetime, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type JSON
		{
			name:        "MySQL field col1 with type JSON, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "MySQL field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
	}
	postgresTestCases := []test{
		{
			name:        "Postgres field col1 with type ID",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID having different size",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(500))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "YES", VarcharSize: 500}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeString) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeBoolean) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeFloat) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeJSON) + " )"},

			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Postgres field col1 with type DateTime",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		// Check Required fields(!)
		{
			name:        "Postgres field col1 which is not null with type ID ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type String ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeString) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type Boolean ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeBoolean) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type Integer ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type Float ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeFloat) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type DateTime ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type JSON ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeJSON) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		// Currently @createdAt,@updateAt work at application level, so we cannot fetch data regarding this from our database
		// But this is possible if we do the above operation at database level
		// TODO: Do this when @createdAt & @updatedAt is pushed to database level
		// {
		// 	name: "Postgres field col1 which is not null with type DateTime having directive @createdAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// {
		// 	name: "Postgres field col1 which is not null with type DateTime having directive @updatedAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// Check @default directive
		// NOTE: JSON & text type cannot have default value
		{
			name:        "Postgres field col1 which is not null with type ID having default value INDIA",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) NOT NULL DEFAULT 'INDIA')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldDefault: "INDIA", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 boolean NOT NULL DEFAULT true)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldDefault: "true", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bit NOT NULL DEFAULT 1)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys:    []model.ForeignKeysType{},
			indexKeys:      []model.IndexType{},
			fields:         []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldDefault: "1", VarcharSize: model.SQLTypeIDSize}},
			wantErr:        false,
			isPostgresSkip: true,
			isMysqlSkip:    true,
		},
		{
			name:        "Postgres field col1 which is not null with type Integer having default value 100",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint NOT NULL DEFAULT 100)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", FieldDefault: "100", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type Float having default value 9.8",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 float NOT NULL DEFAULT 9.8)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", FieldDefault: "9.8", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldDefault: "2020-05-30 00:42:05", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys:    []model.ForeignKeysType{},
			indexKeys:      []model.IndexType{},
			fields:         []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", VarcharSize: model.SQLTypeIDSize}},
			wantErr:        false,
			isMysqlSkip:    true,
			isPostgresSkip: true,
		},
		// Check primary key constraint
		// NOTE: We are only checking for type ID as currently space cloud supports that only
		{
			name:        "Postgres field col1 with type ID which is not null having primary key constraint",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) PRIMARY KEY NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},

		// Check foreign key constraint
		// NOTE: We cannot create foreign keys on string,
		// TODO: FieldKey is empty
		{
			name:        "Postgres field col2 with type ID which is not null having foreign key constraint created through from space cloud",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeID) + " PRIMARY KEY NOT NULL)", "CREATE TABLE myproject.table2 (column2 " + getSQLType(*dbType, model.TypeID) + " NOT NULL, CONSTRAINT c_table1_column1 FOREIGN KEY (column2) REFERENCES myproject.table1(column1))"},
			args: args{
				ctx: context.Background(),
				col: "table2",
			},
			foreignKeys: []model.ForeignKeysType{{TableName: "table2", ColumnName: secondColumn, RefTableName: "table1", RefColumnName: firstColumn, ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO ACTION"}},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Postgres field col2 with type ID which is not null having foreign key constraint with delete rule cascade created through from space cloud",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeID) + " PRIMARY KEY NOT NULL)", "CREATE TABLE myproject.table2 (column2 " + getSQLType(*dbType, model.TypeID) + " NOT NULL, CONSTRAINT c_table1_column1 FOREIGN KEY (column2) REFERENCES myproject.table1(column1) ON DELETE CASCADE)"},
			args: args{
				ctx: context.Background(),
				col: "table2",
			},
			foreignKeys: []model.ForeignKeysType{{TableName: "table2", ColumnName: secondColumn, RefTableName: "table1", RefColumnName: firstColumn, ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "CASCADE"}},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		// Check Unique Constraint
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type ID which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type Boolean which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type String which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type Integer which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type Datetime which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type Float which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Postgres field col1 with type JSON which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Id
		{
			name:        "Postgres field col1 with type ID, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type String
		{
			name:        "Postgres field col1 with type String, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Integer
		{
			name:        "Postgres field col1 with type Integer, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Float
		{
			name:        "Postgres field col1 with type Float, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Boolean
		{
			name:        "Postgres field col1 with type Boolean, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Datetime
		{
			name:        "Postgres field col1 with type Datetime, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type JSON
		{
			name:        "Postgres field col1 with type JSON, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},

		// Check Index Constraint
		{
			// TODO: Field key is empty
			name:        "Postgres field col1 with type ID which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Id
		{
			name:        "Postgres field col1 with type ID, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type ID, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type String
		{
			name:        "Postgres field col1 with type String, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type String, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Integer
		{
			name:        "Postgres field col1 with type Integer, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Integer, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Float
		{
			name:        "Postgres field col1 with type Float, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Float, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Boolean
		{
			name:        "Postgres field col1 with type Boolean, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Boolean, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type Datetime
		{
			name:        "Postgres field col1 with type Datetime, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type Datetime, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		// Type JSON
		{
			name:        "Postgres field col1 with type JSON, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Postgres field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
	}
	mssqlTestCases := []test{
		{
			name:        "Sql Server field col1 with type ID",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID having different size",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(500))"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: "varchar(500)", FieldNull: "YES", VarcharSize: 500}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type String",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeString) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: "column1", FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "YES", VarcharSize: -1}},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeBoolean) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeFloat) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type JSON",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeJSON) + " )"},

			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 with type DateTime",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " )"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "YES", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		// Check Required fields(!)
		{
			name:        "Sql Server field col1 which is not null with type ID ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type String ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeString) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeString), FieldNull: "NO", VarcharSize: -1}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type Boolean ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeBoolean) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type Integer ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type Float ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeFloat) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type DateTime ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type JSON ",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeJSON) + " NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   []model.IndexType{},
			foreignKeys: []model.ForeignKeysType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeJSON), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		// Currently @createdAt,@updateAt work at application level, so we cannot fetch data regarding this from our database
		// But this is possible if we do the above operation at database level
		// TODO: Do this when @createdAt & @updatedAt is pushed to database level
		// {
		// 	name: "Sql Server field col1 which is not null with type DateTime having directive @createdAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// {
		// 	name: "Sql Server field col1 which is not null with type DateTime having directive @updatedAt",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col:         "table1",
		// 		fields:      []utils.FieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO"}},
		// 		foreignKeys: []utils.ForeignKeysType{},
		// 	},// 	wantErr: false,
		// },
		// Check @default directive
		// NOTE: JSON & text type cannot have default value
		{
			name:        "Sql Server field col1 which is not null with type ID having default value INDIA",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) NOT NULL DEFAULT 'INDIA')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldDefault: "INDIA", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 boolean NOT NULL DEFAULT true)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldDefault: "true", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 which is not null with type Boolean having default value true",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bit NOT NULL DEFAULT 1)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys:    []model.ForeignKeysType{},
			indexKeys:      []model.IndexType{},
			fields:         []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", FieldDefault: "1", VarcharSize: model.SQLTypeIDSize}},
			wantErr:        false,
			isPostgresSkip: true,
			isMysqlSkip:    true,
		},
		{
			name:        "Sql Server field col1 which is not null with type Integer having default value 100",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 bigint NOT NULL DEFAULT 100)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", FieldDefault: "100", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type Float having default value 9.8",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 float NOT NULL DEFAULT 9.8)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", FieldDefault: "9.8", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldDefault: "2020-05-30 00:42:05", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
			isMssqlSkip: true,
		},
		{
			name:        "Sql Server field col1 which is not null with type DateTime having default value 2020-05-30T00:42:05+00:00",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeDateTime) + " NOT NULL DEFAULT '2020-05-30T00:42:05+00:00')"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys:    []model.ForeignKeysType{},
			indexKeys:      []model.IndexType{},
			fields:         []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", FieldDefault: "2020-05-30T00:42:05+00:00", VarcharSize: model.SQLTypeIDSize}},
			wantErr:        false,
			isMysqlSkip:    true,
			isPostgresSkip: true,
		},
		// Check primary key constraint
		// NOTE: We are only checking for type ID as currently space cloud supports that only
		{
			name:        "Sql Server field col1 with type ID which is not null having primary key constraint",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 varchar(50) PRIMARY KEY NOT NULL)"},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			foreignKeys: []model.ForeignKeysType{},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "PRI", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},

		// Check foreign key constraint
		// NOTE: We cannot create foreign keys on string,
		{
			name:        "Sql Server field col2 with type ID which is not null having foreign key constraint created through from space cloud",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeID) + " PRIMARY KEY NOT NULL)", "CREATE TABLE myproject.table2 (column2 " + getSQLType(*dbType, model.TypeID) + " NOT NULL, CONSTRAINT c_table1_column1 FOREIGN KEY (column2) REFERENCES myproject.table1(column1))"},
			args: args{
				ctx: context.Background(),
				col: "table2",
			},
			foreignKeys: []model.ForeignKeysType{{TableName: "table2", ColumnName: secondColumn, RefTableName: "table1", RefColumnName: firstColumn, ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "NO ACTION"}},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col2 with type ID which is not null having foreign key constraint with delete rule cascade created through from space cloud",
			createQuery: []string{"CREATE TABLE myproject.table1 (column1 " + getSQLType(*dbType, model.TypeID) + " PRIMARY KEY NOT NULL)", "CREATE TABLE myproject.table2 (column2 " + getSQLType(*dbType, model.TypeID) + " NOT NULL, CONSTRAINT c_table1_column1 FOREIGN KEY (column2) REFERENCES myproject.table1(column1) ON DELETE CASCADE)"},
			args: args{
				ctx: context.Background(),
				col: "table2",
			},
			foreignKeys: []model.ForeignKeysType{{TableName: "table2", ColumnName: secondColumn, RefTableName: "table1", RefColumnName: firstColumn, ConstraintName: getConstraintName("table1", firstColumn), DeleteRule: "CASCADE"}},
			indexKeys:   []model.IndexType{},
			fields:      []model.InspectorFieldType{{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", FieldKey: "MUL", VarcharSize: model.SQLTypeIDSize}},
			wantErr:     false,
		},
		// Check Unique Constraint
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type ID which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type Boolean which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type String which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type Integer which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type Datetime which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type Float which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			// TODO: FieldKey is empty
			name:        "Sql Server field col1 with type JSON which is not null having single unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Id
		{
			name:        "Sql Server field col1 with type ID, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type String
		{
			name:        "Sql Server field col1 with type String, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Integer
		{
			name:        "Sql Server field col1 with type Integer, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Float
		{
			name:        "Sql Server field col1 with type Float, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Boolean
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Datetime
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "yes"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "yes"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type JSON
		{
			name:        "Sql Server field col1 with type JSON, col2 with type ID which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Integer which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Float which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type String which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Boolean which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Datetime which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type JSON which is not null having multiple unique index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE UNIQUE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},

		// Check Index Constraint
		{
			// TODO: Field key is empty
			name:        "Sql Server field col1 with type ID which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type String which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Integer which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1) ", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      []model.InspectorFieldType{{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize}},
			indexKeys:   []model.IndexType{{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"}},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type JSON which is not null having single index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s ON myproject.table1 (column1)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			fields:      nil,
			indexKeys:   nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Id
		{
			name:        "Sql Server field col1 with type ID, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type ID, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeID), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type String
		{
			name:        "Sql Server field col1 with type String, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type String, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeString), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Integer
		{
			name:        "Sql Server field col1 with type Integer, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Integer, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeInteger), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Float
		{
			name:        "Sql Server field col1 with type Float, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Float, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeFloat), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Boolean
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Boolean, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeBoolean), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type Datetime
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeID), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeInteger), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeFloat), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeBoolean), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys: []model.IndexType{
				{TableName: "table1", ColumnName: firstColumn, IndexName: getIndexName("table1", "index1"), Order: 1, Sort: model.DefaultIndexSort, IsUnique: "no"},
				{TableName: "table1", ColumnName: secondColumn, IndexName: getIndexName("table1", "index1"), Order: 2, Sort: model.DefaultIndexSort, IsUnique: "no"},
			},
			fields: []model.InspectorFieldType{
				{FieldName: firstColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
				{FieldName: secondColumn, FieldType: getSQLType(*dbType, model.TypeDateTime), FieldNull: "NO", VarcharSize: model.SQLTypeIDSize},
			},
			foreignKeys: []model.ForeignKeysType{},
			wantErr:     false,
		},
		{
			name:        "Sql Server field col1 with type Datetime, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeDateTime), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		// Type JSON
		{
			name:        "Sql Server field col1 with type JSON, col2 with type ID which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeID)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Integer which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeInteger)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Float which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeFloat)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type String which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeString)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Boolean which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeBoolean)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type Datetime which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeDateTime)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
		{
			name:        "Sql Server field col1 with type JSON, col2 with type JSON which is not null having multiple index constraint created through space cloud",
			createQuery: []string{fmt.Sprintf("CREATE TABLE myproject.table1 (column1 %s NOT NULL, column2 %s NOT NULL)", getSQLType(*dbType, model.TypeJSON), getSQLType(*dbType, model.TypeJSON)), fmt.Sprintf("CREATE INDEX %s on myproject.table1 (column1,column2)", getIndexName("table1", "index1"))},
			args: args{
				ctx: context.Background(),
				col: "table1",
			},
			indexKeys:   nil,
			fields:      nil,
			foreignKeys: nil,
			wantErr:     true,
		},
	}

	switch model.DBType(*dbType) {
	case model.MySQL:
		tests = mysqlTestCases
	case model.Postgres:
		tests = postgresTestCases
	case model.SQLServer:
		tests = mssqlTestCases
	}

	db, err := Init(model.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("DescribeTable() Couldn't establishing connection with database", dbType)
	}
	clean := func() {
		if _, err := db.client.Exec("DROP TABLE IF EXISTS myproject.table2"); err != nil {
			t.Log("DescribeTable() Couldn't truncate table", err)
		}
		if _, err := db.client.Exec("DROP TABLE IF EXISTS myproject.table1"); err != nil {
			t.Log("DescribeTable() Couldn't truncate table", err)
		}
	}
	clean()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if *dbType == string(model.MySQL) && tt.isMysqlSkip {
				return
			}
			if *dbType == string(model.Postgres) && tt.isPostgresSkip {
				return
			}
			if *dbType == string(model.SQLServer) && tt.isMssqlSkip {
				return
			}
			// create table in db
			t.Log("Creating query")
			if err := db.RawBatch(context.Background(), tt.createQuery); err != nil {
				t.Logf("DescribeTable() couldn't execute raw query error - (%v)", err)
				if tt.wantErr {
					clean()
					return
				}
			}

			got, got1, got2, err := db.DescribeTable(tt.args.ctx, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.fields) {
				t.Errorf("DescribeTable() got = %v, want %v", got, tt.fields)
			}
			if !reflect.DeepEqual(got1, tt.foreignKeys) {
				t.Errorf("DescribeTable() got1 = %v, want %v", got1, tt.foreignKeys)
			}
			if !reflect.DeepEqual(got2, tt.indexKeys) {
				t.Errorf("DescribeTable() got2 = %v, want %v", got2, tt.indexKeys)
			}
			clean()
		})
	}
}

func getIndexName(tableName, indexName string) string {
	return fmt.Sprintf("index__%s__%s", tableName, indexName)
}

func getConstraintName(tableName, columnName string) string {
	return fmt.Sprintf("c_%s_%s", tableName, columnName)
}

func getSQLType(dbType, typename string) string {

	switch typename {
	case model.TypeID:
		if dbType == string(model.Postgres) {
			return "character varying"
		}
		return fmt.Sprintf("varchar(%d)", model.SQLTypeIDSize)
	case model.TypeString:
		if dbType == string(model.SQLServer) {
			return "varchar(max)"
		}
		return "text"
	case model.TypeDateTime:
		switch dbType {
		case string(model.MySQL):
			return "datetime"
		case string(model.SQLServer):
			return "datetimeoffset"
		default:
			return "timestamp without time zone"
		}
	case model.TypeBoolean:
		if dbType == string(model.SQLServer) {
			return "bit"
		}
		if dbType == string(model.MySQL) {
			return "tinyint"
		}
		return "boolean"
	case model.TypeFloat:
		if dbType == string(model.Postgres) {
			return "double precision"
		}
		return "float"
	case model.TypeInteger:
		return "bigint"
	case model.TypeJSON:
		switch dbType {
		case string(model.Postgres):
			return "jsonb"
		case string(model.MySQL):
			return "json"
		}
	}
	return ""
}
