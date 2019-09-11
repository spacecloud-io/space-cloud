package schema

import (
	"context"
	"fmt"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

func TestSchema_schemaInspection(t *testing.T) {
	type fields struct {
		schemaDoc schemaType
		crud      *crud.Module
		project   string
	}
	type args struct {
		ctx    context.Context
		dbType string
		col    string
		dbName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "SQL schema Persons",
			args: args{
				dbType: "sql-mysql",
				col:    "persons",
				dbName: "testDB",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testDB",
			},
		},
		{
			name: "SQL schema Orders",
			args: args{
				dbType: "sql-mysql",
				col:    "orders",
				dbName: "testDB",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testDB",
			},
		},
		{
			name: "Postgress schema so_headers",
			args: args{
				dbType: "sql-postgres",
				col:    "so_headers",
				dbName: "testDB",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testDB",
			},
		},
		{
			name: "Postgress schema so_items",
			args: args{
				dbType: "sql-postgres",
				col:    "so_items",
				dbName: "testDB",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testDB",
			},
		},
		{
			name: "Postgress schema account",
			args: args{
				dbType: "sql-postgres",
				col:    "account",
				dbName: "testDB",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testDB",
			},
		},
		{
			name: "Postgress schema sandeep",
			args: args{
				dbType: "sql-postgres",
				col:    "sandeep",
				dbName: "testDB",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testDB",
			},
		},
	}

	db := config.Crud{
		"sql-mysql": &config.CrudStub{
			Conn: "root:1234@tcp(172.17.0.2:3306)/testDB",
			Collections: map[string]*config.TableRule{
				"Persons": &config.TableRule{},
				"Orders":  &config.TableRule{},
			},
			Enabled: true,
		},
		"sql-postgres": &config.CrudStub{
			Conn: "postgres://postgres:1234@172.17.0.3:5432/testdb?sslmode=disable",
			Collections: map[string]*config.TableRule{
				"Persons": &config.TableRule{},
				"Orders":  &config.TableRule{},
			},
			Enabled: true,
		},
	}

	crud := crud.Init()
	if err := crud.SetConfig(db); err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				schemaDoc: tt.fields.schemaDoc,
				crud:      crud,
				project:   tt.fields.project,
			}

			result, err := s.SchemaInspection(tt.args.ctx, tt.args.dbType, s.project, tt.args.col)
			if err != nil {
				t.Errorf("Schema.schemaInspection() error = %v", err)
			}
			fmt.Println(result)
		})
	}
}
