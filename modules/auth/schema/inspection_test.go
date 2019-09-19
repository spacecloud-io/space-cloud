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
		SchemaDoc schemaType
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
			name: "MySQL schema Persons",
			args: args{
				dbType: "sql-mysql",
				col:    "persons",
				dbName: "testdb",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testdb",
			},
		},
		{
			name: "MySQL schema Orders",
			args: args{
				dbType: "sql-mysql",
				col:    "orders",
				dbName: "testdb",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testdb",
			},
		},
		{
			name: "Postgress schema persons",
			args: args{
				dbType: "sql-postgres",
				col:    "persons",
				dbName: "testdb",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testdb",
			},
		},
		{
			name: "Postgress schema orders",
			args: args{
				dbType: "sql-postgres",
				col:    "orders",
				dbName: "testdb",
			},
			fields: fields{
				crud:    &crud.Module{},
				project: "testdb",
			},
		},
	}

	db := config.Crud{
		"sql-mysql": &config.CrudStub{
			Conn: "root:1234@tcp(172.17.0.2:3306)/testdb",
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
				SchemaDoc: tt.fields.SchemaDoc,
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
