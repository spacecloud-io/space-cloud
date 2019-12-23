package schema

import (
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"reflect"
	"testing"
)

func TestParseSchema(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		schema        schemaType
		Data          config.Crud
	}{
		{
			name:          "compulsory field with different datatypes/primary key on list",
			IsErrExpected: true,
			schema:        nil,
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
			name:          "invalid collection name",
			schema:        nil,
			IsErrExpected: true,
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
			name:          "invalid linked field and valid directives",
			schema:        nil,
			IsErrExpected: true,
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
			name:          "collection could not be found in schema",
			schema:        nil,
			IsErrExpected: true,
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
						"age": &SchemaFieldType{
							FieldName: "age",
							Kind:      typeFloat,
						},
					},
				},
			},
			IsErrExpected: false,
			Data: config.Crud{
				"mongo": &config.CrudStub{
					Collections: map[string]*config.TableRule{
						"tweet": &config.TableRule{
							Schema: `type tweet {
						 ID : ID @primary
						 age: Float
						}`,
						},
					},
				},
			},
		},
	}

	s := Init(&crud.Module{}, false)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r, err := s.parser(testCase.Data)
			if (err != nil) != testCase.IsErrExpected {
				t.Errorf("\n Schema.parseSchema() error = expected error-%v,got error-%v", testCase.IsErrExpected, err)
			} else if !reflect.DeepEqual(r, testCase.schema) {
				t.Errorf("parser()=got return value-%v,expected schema-%v", r, testCase.schema)
			}
		})
	}
}
