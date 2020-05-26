package schema

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

func TestParseSchema(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		schema        model.Type
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
			name:          "field not provided in schema",
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
						 person : sharad @link()

						}`,
						},
					},
				},
			},
		},
		{
			name:          "value not provided for default",
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
						 person : sharad @default

						}`,
						},
					},
				},
			},
		},
		{
			name:          "wrong directive provided",
			schema:        nil,
			IsErrExpected: true,
			Data: config.Crud{
				"mongo": &config.CrudStub{
					Collections: map[string]*config.TableRule{
						"test": &config.TableRule{
							Schema: `type test {
						 id : ID @primary
						 person : sharad @de

						}`,
						},
					},
				},
			},
		},
		{
			name:          "wrong args provided for group in directive-index",
			schema:        nil,
			IsErrExpected: true,
			Data: config.Crud{
				"mongo": &config.CrudStub{
					Collections: map[string]*config.TableRule{
						"test": &config.TableRule{
							Schema: `type test {
						 id : ID @primary
						 first_name: ID! @index(group: 10, order: 1, sort: "asc")
						}`,
						},
					},
				},
			},
		},
		{
			name: "OnDelete with NO ACTION",
			schema: model.Type{
				"mongo": model.Collection{
					"tweet": model.Fields{
						"ID": &model.FieldType{
							FieldName: "ID",
							Kind:      model.TypeID,
							IsPrimary: true,
						},
						"age": &model.FieldType{
							FieldName: "age",
							Kind:      model.TypeFloat,
						},
						"spec": &model.FieldType{
							FieldName: "spec",
							Kind:      model.TypeJSON,
						},
						"customer_id": &model.FieldType{
							FieldName:           "customer_id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsForeign:           true,
							JointTable: &model.TableProperties{
								To:             "id",
								Table:          "customer",
								OnDelete:       "NO ACTION",
								ConstraintName: "c_tweet_customer_id",
							},
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
						 spec: JSON
						 customer_id: ID! @foreign(table: "customer", field: "id", onDelete: "ca")
						}`,
						},
					},
				},
			},
		},
		{
			name: "valid schema",
			schema: model.Type{
				"mongo": model.Collection{
					"tweet": model.Fields{
						"ID": &model.FieldType{
							FieldName: "ID",
							Kind:      model.TypeID,
							IsPrimary: true,
						},
						"age": &model.FieldType{
							FieldName: "age",
							Kind:      model.TypeFloat,
						},
						"role": &model.FieldType{
							FieldName:           "role",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsDefault:           true,
							Default:             "user",
						},
						"spec": &model.FieldType{
							FieldName: "spec",
							Kind:      model.TypeJSON,
						},
						"createdAt": &model.FieldType{
							FieldName:   "createdAt",
							Kind:        model.TypeDateTime,
							IsCreatedAt: true,
						},
						"updatedAt": &model.FieldType{
							FieldName:   "updatedAt",
							Kind:        model.TypeDateTime,
							IsUpdatedAt: true,
						},
						"first_name": &model.FieldType{
							FieldName:           "first_name",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsIndex:             true,
							IndexInfo: &model.TableProperties{
								Group: "user_name",
								Order: 1,
								Sort:  "asc",
							},
						},
						"name": &model.FieldType{
							FieldName:           "name",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsIndex:             true,
							IsUnique:            true,
							IndexInfo: &model.TableProperties{
								Group: "user_name",
								Order: 1,
								Sort:  "asc",
							},
						},
						"customer_id": &model.FieldType{
							FieldName:           "customer_id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsForeign:           true,
							JointTable: &model.TableProperties{
								To:             "id",
								Table:          "customer",
								OnDelete:       "CASCADE",
								ConstraintName: "c_tweet_customer_id",
							},
						},
						"order_dates": &model.FieldType{
							FieldName: "order_dates",
							IsList:    true,
							Kind:      model.TypeDateTime,
							IsLinked:  true,
							LinkedTable: &model.TableProperties{
								Table:  "order",
								From:   "id",
								To:     "customer_id",
								Field:  "order_date",
								DBType: "mongo",
							},
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
						 spec: JSON
						 createdAt:DateTime@createdAt
						 updatedAt:DateTime@updatedAt
						 role: ID! @default(value: "user")
						 first_name: ID! @index(group: "user_name", order: 1, sort: "asc")
						 name: ID! @unique(group: "user_name", order: 1)
						 customer_id: ID! @foreign(table: "customer", field: "id", onDelete: "cascade")
						 order_dates: [DateTime] @link(table: "order", field: "order_date", from: "id", to: "customer_id")
						}`,
						},
					},
				},
			},
		},
	}

	s := Init(&crud.Module{})
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r, err := s.Parser(testCase.Data)
			if (err != nil) != testCase.IsErrExpected {
				t.Errorf("\n Schema.parseSchema() error = expected error-%v,got error-%v", testCase.IsErrExpected, err)
				return
			}
			if !reflect.DeepEqual(r, testCase.schema) {
				t.Errorf("parser()=got return value-%v,expected schema-%v", r, testCase.schema)
			}
		})
	}
}
