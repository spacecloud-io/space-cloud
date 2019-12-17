package schema

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

// func TestSchema_SchemaModifyAll(t *testing.T) {
// 	type fields struct {
// 		lock               sync.RWMutex
// 		SchemaDoc          schemaType
// 		crud               *crud.Module
// 		project            string
// 		config             config.Crud
// 		removeProjectScope bool
// 	}
// 	type args struct {
// 		ctx     context.Context
// 		dbAlias string
// 		project string
// 		tables  map[string]*config.TableRule
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		queries []string
// 		wantErr bool
// 		want    error
//	}{
//TODO: Add test cases.
// {
// 	name: ": add a new table with column of String type",
// 	args: args{
// 		dbAlias: "sql-mysql",
// 		project: "test",
// 		tables: map[string]*config.TableRule{
// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String}`},
// 		},
// 	},
// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 	wantErr: false,
// 	want:    nil,
// 	queies: []string{
// 		// "CREATE TABLE test.table1 ()"
// 	},
// },
// 		// {
// 		// 	name: ": add a new column of floating type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a new column of integer type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float, col4: Integer}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a new column of boolean type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a new column of datetime type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean col6: DateTime}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a new column of ID type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean col6: DateTime col7: ID}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a primary key to ID type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean col6: DateTime col7: ID! @primary}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": remove primary key",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean col6: DateTime col7: ID}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a unique key to integer type column",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer @unique col5: Boolean col6: DateTime col7: ID}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": remove unique key",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean col6: DateTime col7: ID}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a foreign key to integer type column",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table2": &config.TableRule{Schema: `type table2 {col1: Float! @unique col2: String}`},
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float! @foreign(table: "table2", to: "col1") col4: Integer col5: Boolean col6: DateTime col7: ID}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": remove the foreign key constraint ",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String col3: Float col4: Integer col5: Boolean col6: DateTime col7: ID!}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": bring table1 to initial state",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"table1": &config.TableRule{Schema: `type table1 {col1: String col2: String}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": wrong name for dbAlias",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysl",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"student": &config.TableRule{Schema: `type student {id: ID! @primary name: String college_id: ID @foreign(table: "college", field: "id") age: Integer stipend: Float dob: DateTime}`},
// 		// 			"college": &config.TableRule{Schema: `type college {id: ID! @primary name: String city_code: Integer estd: DateTime}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: true,
// 		// 	want:    errors.New("crud module dbalias \"sql-mysl\" not found"),
// 		// },
// 		// {
// 		// 	name: ": only single column added in table",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"coll": &config.TableRule{Schema: `type coll {id: ID! @primary }`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": empty info.schema",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"student": &config.TableRule{Schema: ``},
// 		// 			"college": &config.TableRule{Schema: ``},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },

// 		// {
// 		// 	name: ": no changes copy of actual db",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"student": &config.TableRule{Schema: `type student {id: ID! @primary name: String college_id: ID @foreign(table: "college", field: "id") age: Integer stipend: Float dob: DateTime}`},
// 		// 			"college": &config.TableRule{Schema: `type college {id: ID! @primary name: String city_code: Integer estd: DateTime}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },

// 		// {
// 		// 	name: ": add new table",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"student": &config.TableRule{Schema: `type student {id: ID! @primary name: String college_id: ID @foreign(table: "college", field: "id") age: Integer stipend: Float dob: DateTime}`},
// 		// 			"coll":    &config.TableRule{Schema: `type coll {id: ID! @primary name: String city_code: Integer estd: DateTime}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": changes in columns name and type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"coll": &config.TableRule{Schema: `type coll {id: ID! @primary name: String city_code: Integer estd: Integer}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add new table and primary key",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"primaryTable": &config.TableRule{Schema: `type primaryTable {col1: String col2: Float col3: ID @primary}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a foreign key",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"foreignTable": &config.TableRule{Schema: `type foreignTable {col1: String col2: ID@foreign(table: "college", field: "id")}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a unique key",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"uniqueTable": &config.TableRule{Schema: `type uniqueTable {col1: String col2: ID@unique}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },

// 		// {
// 		// 	name: ": add a integer type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"integerTable": &config.TableRule{Schema: `type integerTable {id: ID! @primary  col1: Integer}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a boolean type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"booleanTable": &config.TableRule{Schema: `type booleanTable {id: ID! @primary  col1: Boolean}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 		// {
// 		// 	name: ": add a datatime type",
// 		// 	args: args{
// 		// 		dbAlias: "sql-mysql",
// 		// 		project: "test",
// 		// 		tables: map[string]*config.TableRule{
// 		// 			"datetimeTable": &config.TableRule{Schema: `type datetimeTable {id: ID! @primary  col1: DateTime}`},
// 		// 		},
// 		// 	},
// 		// 	fields:  fields{crud: crud.Init(false), project: "test"},
// 		// 	wantErr: false,
// 		// 	want:    nil,
// 		// },
// 	}

// 	dbNames := config.Crud{
// 		"sql-mysql": &config.CrudStub{
// 			Conn: "root:1234@tcp(localhost:3306)/test",
// 			Collections: map[string]*config.TableRule{
// 				"Persons": &config.TableRule{},
// 				"Orders":  &config.TableRule{},
// 			},
// 			Enabled: true,
// 		},
// 		/*
// 			"sql-postgres": &config.CrudStub{
// 				Conn: "postgres://postgres:1234@172.17.0.3:5432/testdb?sslmode=disable",
// 				Collections: map[string]*config.TableRule{
// 					"Persons": &config.TableRule{},
// 					"Orders":  &config.TableRule{},
// 				},
// 				Enabled: true,
// 			},
// 			"sql-sqlserver": &config.CrudStub{
// 				Conn: "postgres://postgres:1234@172.17.0.3:5432/testdb?sslmode=disable",
// 				Collections: map[string]*config.TableRule{
// 					"Persons": &config.TableRule{},
// 					"Orders":  &config.TableRule{},
// 				},
// 				Enabled: true,
// 			},
// 		*/
// 	}

// 	s := &Schema{
// 		crud:    crud.Init(false),
// 		project: "test",
// 	}
// 	if err := s.crud.SetConfig("test", dbNames); err != nil {
// 		t.Fatal(err)
// 	}
// 	for i := range dbNames {
// 		for _, tt := range tests {
// 			tt.name = string(i) + tt.name
// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()
// 			t.Run(tt.name, func(t *testing.T) {
// 				tt.args.ctx = ctx
// 				if err := s.SchemaModifyAll(tt.args.ctx, tt.args.dbAlias, tt.args.project, tt.args.tables); (tt.wantErr && err == tt.want) || (!tt.wantErr && err != tt.want) {
// 					t.Errorf("Schema.SchemaModifyAll() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 			})
// 		}
// 	}
// }

func TestSchema_generateCreationQueries(t *testing.T) {
	type fields struct {
		SchemaDoc          schemaType
		crud               *crud.Module
		project            string
		config             config.Crud
		removeProjectScope bool
	}
	type args struct {
		ctx           context.Context
		dbAlias       string
		tableName     string
		project       string
		parsedSchema  schemaType
		currentSchema schemaCollection
	}

	crudPostgres := crud.Init(false)
	crudPostgres.SetConfig("test", config.Crud{"postgres": {Type: "sql-postgres", Enabled: false}})

	crudMySql := crud.Init(false)
	crudMySql.SetConfig("test", config.Crud{"mysql": {Type: "sql-mysql", Enabled: false}})

	crudSqlServer := crud.Init(false)
	crudSqlServer.SetConfig("test", config.Crud{"sqlserver": {Type: "sql-sqlserver", Enabled: false}})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "adding two columns",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}, "col2": &SchemaFieldType{FieldName: "col2", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)", "ALTER TABLE test.table1 ADD col2 text"},
			wantErr: false,
		},
		{
			name: "removing one column",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "required to unrequired",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NOT NULL"},
			wantErr: false,
		},
		{
			name: "unrequired to required",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NULL"},
			wantErr: false,
		},
		{
			name: "integer to string",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 text"},
			wantErr: false,
		},
		{
			name: "string to integer",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "integer to float",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "float to integer",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "float to dateTime",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 datetime"},
			wantErr: false,
		},
		{
			name: "datetime to float",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "datetime to id",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},
		{
			name: "id to datetime",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 datetime"},
			wantErr: false,
		},
		{
			name: "adding unique key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 bigint NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 UNIQUE (col1)"},
			wantErr: false,
		},
		{
			name: "removing unique key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: true}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP INDEX c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding primary key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsPrimary: false}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NOT NULL", "ALTER TABLE test.table1 ADD PRIMARY KEY (col1)"},
			wantErr: false,
		},
		{
			name: "removing primary key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsPrimary: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 MODIFY col1 varchar(50) NULL", "ALTER TABLE test.table1 DROP PRIMARY KEY"},
			wantErr: false,
		},
		{
			name: "adding foreign key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: false}}, "table2": SchemaFields{}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 FOREIGN KEY (col1) REFERENCES test.table2 (id)"},
			wantErr: false,
		},
		{
			name: "removing foreign key",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}, "table2": SchemaFields{"id": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP FOREIGN KEY c_table1_col1", "ALTER TABLE test.table1 DROP INDEX c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding link",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "removing link",
			args: args{
				dbAlias:       "mysql",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"mysql": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudMySql, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},

		// //sql-server

		{
			name: "adding two columns",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}, "col2": &SchemaFieldType{FieldName: "col2", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)", "ALTER TABLE test.table1 ADD col2 text"},
			wantErr: false,
		},
		{
			name: "removing one column",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "required to unrequired",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NOT NULL"},
			wantErr: false,
		},
		{
			name: "unrequired to required",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NULL"},
			wantErr: false,
		},
		{
			name: "integer to string",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 text"},
			wantErr: false,
		},
		{
			name: "string to integer",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "integer to float",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "float to integer",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 bigint"},
			wantErr: false,
		},
		{
			name: "float to dateTime",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 timestamp NULL"},
			wantErr: false,
		},
		{
			name: "datetime to float",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 float"},
			wantErr: false,
		},
		{
			name: "datetime to id",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},
		{
			name: "id to datetime",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD col1 timestamp NULL"},
			wantErr: false,
		},
		{
			name: "adding unique key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 bigint NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 UNIQUE (col1)"},
			wantErr: false,
		},
		{
			name: "removing unique key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 bigint NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding primary key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsPrimary: false}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 PRIMARY KEY CLUSTERED (col1)"},
			wantErr: false,
		},
		{
			name: "removing primary key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsPrimary: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 varchar(50) NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding foreign key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: false}}, "table2": SchemaFields{}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 FOREIGN KEY (col1) REFERENCES test.table2 (id)"},
			wantErr: false,
		},
		{
			name: "removing foreign key",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}, "table2": SchemaFields{"id": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding link",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "removing link",
			args: args{
				dbAlias:       "sqlserver",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"sqlserver": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudSqlServer, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD col1 varchar(50)"},
			wantErr: false,
		},

		// //postgres
		{
			name: "adding two columns",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}, "col2": &SchemaFieldType{FieldName: "col2", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD COLUMN col1 varchar(50)", "ALTER TABLE test.table1 ADD COLUMN col2 text"},
			wantErr: false,
		},
		{
			name: "removing one column",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "required to unrequired",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 SET NOT NULL"},
			wantErr: false,
		},
		{
			name: "unrequired to required",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 DROP NOT NULL"},
			wantErr: false,
		},
		{
			name: "integer to string",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 text"},
			wantErr: false,
		},
		{
			name: "string to integer",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeString}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 bigint"},
			wantErr: false,
		},
		{
			name: "integer to float",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 float"},
			wantErr: false,
		},
		{
			name: "float to integer",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 bigint"},
			wantErr: false,
		},
		{
			name: "float to dateTime",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 timestamp"},
			wantErr: false,
		},
		{
			name: "datetime to float",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeFloat}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 float"},
			wantErr: false,
		},
		{
			name: "datetime to id",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 varchar(50)"},
			wantErr: false,
		},
		{
			name: "id to datetime",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeDateTime}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1", "ALTER TABLE test.table1 ADD COLUMN col1 timestamp"},
			wantErr: false,
		},
		{
			name: "adding unique key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 SET NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 UNIQUE (col1)"},
			wantErr: false,
		},
		{
			name: "removing unique key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsUnique: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: typeInteger, IsFieldTypeRequired: true, IsUnique: true}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 DROP NOT NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding primary key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsPrimary: false}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 SET NOT NULL", "ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 PRIMARY KEY (col1)"},
			wantErr: false,
		},
		{
			name: "removing primary key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsPrimary: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: true, IsPrimary: true}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ALTER COLUMN col1 DROP NOT NULL", "ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding foreign key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: false}}, "table2": SchemaFields{}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD CONSTRAINT c_table1_col1 FOREIGN KEY (col1) REFERENCES test.table2 (id)"},
			wantErr: false,
		},
		{
			name: "removing foreign key",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}, "table2": SchemaFields{"id": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsForeign: true, JointTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP CONSTRAINT c_table1_col1"},
			wantErr: false,
		},
		{
			name: "adding link",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 DROP COLUMN col1"},
			wantErr: false,
		},
		{
			name: "removing link",
			args: args{
				dbAlias:       "postgres",
				tableName:     "table1",
				project:       "test",
				parsedSchema:  schemaType{"postgres": schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsFieldTypeRequired: false, IsForeign: false}}}},
				currentSchema: schemaCollection{"table1": SchemaFields{"col1": &SchemaFieldType{FieldName: "col1", Kind: TypeID, IsLinked: true, LinkedTable: &TableProperties{Table: "table2", To: "id"}}}},
			},
			fields:  fields{crud: crudPostgres, project: "test"},
			want:    []string{"ALTER TABLE test.table1 ADD COLUMN col1 varchar(50)"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc:          tt.fields.SchemaDoc,
				crud:               tt.fields.crud,
				project:            tt.fields.project,
				config:             tt.fields.config,
				removeProjectScope: tt.fields.removeProjectScope,
			}
			got, err := s.generateCreationQueries(tt.args.ctx, tt.args.dbAlias, tt.args.tableName, tt.args.project, tt.args.parsedSchema, tt.args.currentSchema)
			if (err != nil) != tt.wantErr {
				t.Errorf("name = %v, Schema.generateCreationQueries() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("name = %v, Schema.generateCreationQueries() = %v, want %v", tt.name, got, tt.want)
				return
			}

			for i, v := range got {
				if tt.want[i] != v {
					t.Errorf("name = %v, Schema.generateCreationQueries() = %v, want %v", tt.name, got, tt.want)
					break
				}
			}
		})
	}
}
