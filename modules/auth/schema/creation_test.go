package schema

import (
	"context"
	"fmt"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
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
			name: "SQL: no changes copy of actual db",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : String      college  : String     personId : ID! @id     age : Integer }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: adding not null to address, college, age field",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : String!      college  : String!     personId : ID! @id     age : Integer! }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: removing not null from address, college, age field",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : String      college  : String     personId : ID! @id     age : Integer }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: changing data type of address string - int & age int to string",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Integer!      college  : String!     personId : ID! @id     age : String }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: changing data type of address int - Boolean & college string - DateTime & age int to Float",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personId : ID! @id      age : Float }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: adding unique key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean! @unique     college  : DateTime!     personId : ID! @id   age : Float }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: removing unique key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personId : ID! @id    age : Float }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: removing primary key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personId : ID! @id     age : Float }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: adding new column",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personId : ID! @id     age : Float    contact : Integer! }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: removing column",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Boolean!      college  : DateTime!     personId : ID! @id     age : Float   }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: adding foreign key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Persons! @relation(field:iD)     college  : DateTime!     personId : ID! @id     age : Float   }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
		{
			name: "SQL: removing foreign key",
			args: args{dbType: "sql-mysql", col: "orders", project: "testDB", schema: `type Persons { lastName  : String!     firstName : String     age : Integer     iD : Integer! @id	}	type Orders { address  : Integer    college  : DateTime!     personId : ID! @id     age : Float   }`},
			fields: fields{crud: crud.Init(), project: "testDB"},
		},
	}
	dbSql := config.Crud{
		"sql-mysql": &config.CrudStub{
			Conn: "root:1234@tcp(172.17.0.2:3306)/testDB",
			Collections: map[string]*config.TableRule{
				"Persons": &config.TableRule{},
				"Orders":  &config.TableRule{},
			},
			Enabled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc: tt.fields.SchemaDoc,
				crud:      tt.fields.crud,
				project:   tt.fields.project,
			}
			if err := s.crud.SetConfig(dbSql); err != nil {
				t.Fatal(err)
			}
			fmt.Print("")
			if err := s.SchemaCreation(tt.args.ctx, tt.args.dbType, tt.args.col, tt.args.project, tt.args.schema); err != nil {
				t.Errorf("Schema.SchemaCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// dbPostgres := config.Crud{
	// 	"sql-postgres": &config.CrudStub{
	// 		Conn: "postgres://postgres:1234@172.17.0.3:5432/testdb?sslmode=disable",
	// 		Collections: map[string]*config.TableRule{
	// 			"Persons": &config.TableRule{},
	// 			"Orders":  &config.TableRule{},
	// 		},
	// 		Enabled: true,
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {

	// 		s := &Schema{
	// 			SchemaDoc: tt.fields.SchemaDoc,
	// 			crud:      tt.fields.crud,
	// 			project:   tt.fields.project,
	// 		}
	// 		if err := s.crud.SetConfig(dbPostgres); err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		s.project = "testdb"
	// 		tt.args.dbType = "sql-postgres"
	// 		tt.args.project = "testdb"
	// 		if err := s.SchemaCreation(tt.args.ctx, tt.args.dbType, tt.args.col, tt.args.project, tt.args.schema); err != nil {
	// 			t.Errorf("Schema.SchemaCreation() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }
}
