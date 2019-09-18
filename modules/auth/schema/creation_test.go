package schema

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

func TestSchema_SchemaCreation(t *testing.T) {
	type fields struct {
		SchemaDoc schemaType
		crud      *crud.Module
		project   string
	}
	type args struct {
		ctx     context.Context
		dbType  string
		col     string
		project string
		schema  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: ": no changes copy of actual db",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : String      college  : String     personid : ID! @id     age : Integer }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": adding not null to address, college, age field",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : String!      college  : String!     personid : ID! @id     age : Integer! }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": removing not null from address, college, age field",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : String      college  : String     personid : ID! @id     age : Integer }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": changing data type of address string - int & age int to string",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Integer!      college  : String!     personid : ID! @id     age : String }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": changing data type of address int - Boolean & college string - DateTime & age int to Float",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personid : ID! @id      age : Float }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": adding unique key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean! @unique     college  : DateTime!     personid : ID! @id   age : Float }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": removing unique key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personid : ID! @id    age : Float }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		// {
		// 	name: ": removing primary key",
		// 	args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personid : ID! @id   age : Float }`},
		// 	fields: fields{crud: crud.Init(), project: "testdb"},
		// },
		{
			name: ": adding new column",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personid : ID! @id     age : Float    contact : Integer! }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": removing column",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personid : ID! @id     age : Float   }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": adding foreign key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Persons! @relation(field:id)     college  : DateTime!     personid : ID! @id     age : Float   }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name: ": removing foreign key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testdb", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Integer    college  : DateTime!     personid : ID! @id     age : Float   }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
		{
			name:   ": adding new table",
			args:   args{dbType: "sql-mysql", col: "college", project: "testdb", schema: `  type College { name : String!     sirname : String!     age : ID! @id   ismale : Boolean!     dob : DateTime! }`},
			fields: fields{crud: crud.Init(), project: "testdb"},
		},
	}
	dbSQL := config.Crud{
		"sql-mysql": &config.CrudStub{
			Conn: "root:1234@tcp(172.17.0.2:3306)/testdb",
			Collections: map[string]*config.TableRule{
				"Persons": &config.TableRule{},
				"Orders":  &config.TableRule{},
			},
			Enabled: true,
		},
	}
	for _, tt := range tests {
		tt.name = string(utils.MySQL) + tt.name
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc: tt.fields.SchemaDoc,
				crud:      tt.fields.crud,
				project:   tt.fields.project,
			}
			if err := s.crud.SetConfig(dbSQL); err != nil {
				t.Fatal(err)
			}
			tt.args.dbType = string(utils.MySQL)

			if err := s.SchemaCreation(tt.args.ctx, tt.args.dbType, tt.args.col, tt.args.project, tt.args.schema); err != nil {
				t.Errorf("Schema.SchemaCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	dbPostgres := config.Crud{
		"sql-postgres": &config.CrudStub{
			Conn: "postgres://postgres:1234@172.17.0.3:5432/testdb?sslmode=disable",
			Collections: map[string]*config.TableRule{
				"Persons": &config.TableRule{},
				"Orders":  &config.TableRule{},
			},
			Enabled: true,
		},
	}
	for _, tt := range tests {
		tt.name = string(utils.Postgres) + tt.name
		t.Run(tt.name, func(t *testing.T) {

			s := &Schema{
				SchemaDoc: tt.fields.SchemaDoc,
				crud:      tt.fields.crud,
				project:   tt.fields.project,
			}
			if err := s.crud.SetConfig(dbPostgres); err != nil {
				t.Fatal(err)
			}
			tt.args.dbType = string(utils.Postgres)
			if err := s.SchemaCreation(tt.args.ctx, tt.args.dbType, tt.args.col, tt.args.project, tt.args.schema); err != nil {
				t.Errorf("Schema.SchemaCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
