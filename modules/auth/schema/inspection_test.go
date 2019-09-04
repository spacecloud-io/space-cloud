package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

func TestSchema_schemaInspection(t *testing.T) {
	type fields struct {
		SchemaDoc SchemaType
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
				col:    "Persons",
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
				col:    "Orders",
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

			result, err := s.schemaInspection(tt.args.ctx, tt.args.dbType, s.project, tt.args.col)
			if err != nil {
				t.Errorf("Schema.schemaInspection() error = %v", err)
			}
			b, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Print(string(b))

		})
	}
}
