package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/modules/crud"

	"github.com/spaceuptech/space-cloud/config"
)

var query = `

   type user {
 	id: ID! @id
 	mentor:sharad @link(table:sharad, from:Name)
   }

   type sharad {
 	  Name : String!
 	  Surname : String!
 	  age : Integer!
 	  isMale : Boolean!
 	  dob : DateTime@createdAt
   }
   type event_logs {
		id:ID@unique
	  	owner: [String]@foreign
   }
 `
var parsedata = []struct {
	name   string
	want   error
	schema schemaType
	Data   config.Crud
}{
	{
		name: "compulsory field with different datatypes",
		want: errors.New("invalid type for field owner - primary and foreign keys cannot be made on lists"),
		schema: schemaType{
			"mongo": schemaCollection{
				"tweet": SchemaFields{
					"id": &SchemaFieldType{
						FieldName:           "id",
						Kind:                TypeID,
						IsFieldTypeRequired: true,
					},
					"createdAt": &SchemaFieldType{
						FieldName: "createdAt",
						Kind:      typeDateTime,
					},
					"exp": &SchemaFieldType{
						FieldName: "exp",
						Kind:      typeInteger,
					},
					"age": &SchemaFieldType{
						FieldName:           "age",
						Kind:                typeFloat,
						IsFieldTypeRequired: true,
					},
					"isMale": &SchemaFieldType{
						FieldName: "age",
						Kind:      typeBoolean,
					},
					"text": &SchemaFieldType{
						FieldName: "text",
						Kind:      typeString,
					},
					"owner": &SchemaFieldType{
						FieldName: "owner",
						Kind:      typeString,
					},
				},
			},
		},
		Data: config.Crud{
			"mongo": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					"tweet": &config.TableRule{
						Schema: `
						type tweet {
							id: ID!
							createdAt:DateTime
							text: String
							isMale: Boolean
							age: Float!
							exp: Integer
							owner:[String]@primary
						  }`,
					},
				},
			},
		},
	},
	{
		name: "invalid collection name",
		schema: schemaType{
			"mongo": schemaCollection{
				"tweet": SchemaFields{
					"id": &SchemaFieldType{
						FieldName: "id",
						Kind:      TypeID,
					},
					"person": &SchemaFieldType{
						FieldName: "createdAt",
						Kind:      typeDateTime,
						IsLinked:  true,
					},
				},
			},
		},
		want: errors.New("collection tes could not be found in schema"),
		Data: config.Crud{
			"mongo": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					"tes": &config.TableRule{
						Schema: `type test {
						 id : ID @id
						 person : sharad @link(table:sharad, from:Name, to:isMale)
						}`,
					},
				},
			},
		},
	},
	{
		name: "invalid linked field and valid directives",
		schema: schemaType{
			"mongo": schemaCollection{
				"tweet": SchemaFields{
					"id": &SchemaFieldType{
						FieldName: "id",
						Kind:      TypeID,
						IsPrimary: true,
					},
					"text": &SchemaFieldType{
						FieldName: "text",
						Kind:      typeString,
						IsUnique:  true,
					},
					"person": &SchemaFieldType{
						FieldName: "person",
						Kind:      typeObject,
					},
					"createdAt": &SchemaFieldType{
						FieldName:   "createdAt",
						Kind:        typeDateTime,
						IsCreatedAt: true,
					},
					"updatedAt": &SchemaFieldType{
						FieldName:   "id",
						Kind:        typeDateTime,
						IsUpdatedAt: true,
					},
					"loc": &SchemaFieldType{
						FieldName: "loc",
						Kind:      typeObject,
						IsForeign: true,
					},
				},
				"location": SchemaFields{
					"latitude": &SchemaFieldType{
						FieldName: "latitude",
						Kind:      typeFloat,
					},
					"text": &SchemaFieldType{
						FieldName: "text",
						Kind:      typeFloat,
					},
				},
			},
		},
		want: errors.New("link directive must be accompanied with to and from fields"),
		Data: config.Crud{
			"mongo": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					"test": &config.TableRule{
						Schema: `type test {
						 id : ID @primary
						 text: String@unique
						 createdAt:DateTime@createdAt
						 updatedAt:DateTime@updatedAt
						 loc:location@foreign(table:location,field:latitude)
						 person : sharad @link(table:sharad, from:Name)
						}
						type location{
							latitude:Float
							longitude:Float
						}`,
					},
				},
			},
		},
	},
	{
		name: "invalid linked field and valid directives",
		schema: schemaType{
			"mongo": schemaCollection{
				"tweet": SchemaFields{
					"id": &SchemaFieldType{
						FieldName: "id",
						Kind:      TypeID,
						IsPrimary: true,
					},
					"person": &SchemaFieldType{
						FieldName: "person",
					},
				},
			},
		},
		want: errors.New("collection Integera could not be found in schema"),
		Data: config.Crud{
			"mongo": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					"test": &config.TableRule{
						Schema: `type test {
						 id : ID @primary
						 text: String@unique
						 createdAt:DateTime@createdAt
						 updatedAt:DateTime@updatedAt
						 exp:Integera
						 person : sharad @link(table:sharad, from:Name)
						 
						}`,
					},
				},
			},
		},
	},
	{
		name: "valid schema",
		schema: schemaType{
			"mongo": schemaCollection{
				"tweet": SchemaFields{
					"ID": &SchemaFieldType{
						FieldName: "ID",
						Kind:      TypeID,
						IsPrimary: true,
					},
				},
			},
		},
		want: nil,
		Data: config.Crud{
			"mongo": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					"test": &config.TableRule{
						Schema: `type test {
						 ID : ID @primary
						 person:sharad@link(table:sharad,from:ID,to:isMale,field:surname)
						}`,
					},
				},
			},
		},
	},
}

func TestParseSchema(t *testing.T) {
	temp := crud.Module{}
	s := Init(&temp, false)

	for _, value := range parsedata {
		t.Run("Schema Parser", func(t *testing.T) {
			if _, err := s.parser(value.Data); err != nil {
				if !reflect.DeepEqual(err, value.want) {
					t.Errorf("\n Schema.parseSchema() error = (%v,%v)", err, value.want)
				}
				/*if !reflect.DeepEqual(r, value.schema) {
					t.Errorf("parser()=%v,want%v", r, value.schema)
				}*/
			}
			// uncomment the below statements to see the reuslt
			b, err := json.MarshalIndent(s.SchemaDoc, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Print(string(b))
			t.Log("Logging Test Output :: ", s.SchemaDoc)
		})
	}
}
