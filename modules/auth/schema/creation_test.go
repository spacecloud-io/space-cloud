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
			name: "Sql",
			args: args{
				dbType:  "sql-mysql",
				col:     "orders",
				project: "testDB",
				schema: `type Persons {
					lastName  : String!
					firstName : String       
					age       : Integer
					iD 		  : Integer! @id					
				}	
				type Orders {
					address  : String!   
					college  : String!  
					personId : ID! @id      
					age      : Integer       
				}
				`,
			},
			fields: fields{
				crud:    crud.Init(),
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				SchemaDoc: tt.fields.SchemaDoc,
				crud:      tt.fields.crud,
				project:   tt.fields.project,
			}
			if err := s.crud.SetConfig(db); err != nil {
				t.Fatal(err)
			}
			fmt.Println("db", tt.args.dbType)
			if err := s.SchemaCreation(tt.args.ctx, tt.args.dbType, tt.args.col, tt.args.project, tt.args.schema); err != nil {
				t.Errorf("Schema.SchemaCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
