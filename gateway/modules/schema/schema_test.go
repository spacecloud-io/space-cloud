package schema

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

func TestParseSchema(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		schema        model.Type
		Data          config.DatabaseSchemas
	}{
		{
			name:          "compulsory field with different datatypes/primary key on list",
			IsErrExpected: true,
			schema:        nil,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
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
		{
			name:          "invalid collection name",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tes"): &config.DatabaseSchema{
					Table:   "tes",
					DbAlias: "mongo",
					Schema: `type test {
						 id : ID @id
						 person : sharad @link(table:sharad, from:Name, to:isMale)
						}`,
				},
			},
		},
		{
			name:          "invalid linked field and valid directives",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "test"): &config.DatabaseSchema{
					Table:   "test",
					DbAlias: "mongo",
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
		{
			name:          "collection could not be found in schema",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "test"): &config.DatabaseSchema{
					Table:   "test",
					DbAlias: "mongo",
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
		{
			name:          "field not provided in schema",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
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
		{
			name:          "value not provided for default",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
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
		{
			name:          "wrong directive provided",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
					Schema: `type test {
						 id : ID @primary
						 person : sharad @de
						}`,
				},
			},
		},
		{
			name:          "wrong args provided for group in directive-index",
			schema:        nil,
			IsErrExpected: true,
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
					Schema: `type test {
						 id : ID @primary
						 first_name: ID! @index(group: 10, order: 1, sort: "asc")
						}`,
				},
			},
		},
		{
			name: "OnDelete with NO ACTION",
			schema: model.Type{
				"mongo": model.Collection{
					"tweet": model.Fields{
						"ID": &model.FieldType{
							FieldName:         "ID",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							Kind:              model.TypeID,
							TypeIDSize:        model.SQLTypeIDSize,
							IsPrimary:         true,
						},
						"age": &model.FieldType{
							FieldName:         "age",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							TypeIDSize:        model.SQLTypeIDSize,
							Kind:              model.TypeFloat,
						},
						"spec": &model.FieldType{
							FieldName:         "spec",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							TypeIDSize:        model.SQLTypeIDSize,
							Kind:              model.TypeJSON,
						},
						"customer_id": &model.FieldType{
							FieldName:           "customer_id",
							AutoIncrementInfo:   new(model.AutoIncrementInfo),
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
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
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
					Schema: `type tweet {
						 ID : ID @primary
						 age: Float
						 spec: JSON
						 customer_id: ID! @foreign(table: "customer", field: "id", onDelete: "ca")
						}`,
				},
			},
		},
		{
			name: "valid schema",
			schema: model.Type{
				"mongo": model.Collection{
					"tweet": model.Fields{
						"ID": &model.FieldType{
							FieldName:         "ID",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							Kind:              model.TypeID,
							TypeIDSize:        model.SQLTypeIDSize,
							IsPrimary:         true,
						},
						"age": &model.FieldType{
							FieldName:         "age",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							Kind:              model.TypeFloat,
							TypeIDSize:        model.SQLTypeIDSize,
						},
						"role": &model.FieldType{
							FieldName:           "role",
							AutoIncrementInfo:   new(model.AutoIncrementInfo),
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
							IsDefault:           true,
							Default:             "user",
						},
						"spec": &model.FieldType{
							FieldName:         "spec",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							Kind:              model.TypeJSON,
							TypeIDSize:        model.SQLTypeIDSize,
						},
						"createdAt": &model.FieldType{
							FieldName:         "createdAt",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							Kind:              model.TypeDateTime,
							TypeIDSize:        model.SQLTypeIDSize,
							IsCreatedAt:       true,
						},
						"updatedAt": &model.FieldType{
							FieldName:         "updatedAt",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							Kind:              model.TypeDateTime,
							TypeIDSize:        model.SQLTypeIDSize,
							IsUpdatedAt:       true,
						},
						"first_name": &model.FieldType{
							FieldName:           "first_name",
							AutoIncrementInfo:   new(model.AutoIncrementInfo),
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
							IsIndex:             true,
							IndexInfo: &model.TableProperties{
								Group: "user_name",
								Order: 1,
								Sort:  "asc",
							},
						},
						"name": &model.FieldType{
							FieldName:           "name",
							AutoIncrementInfo:   new(model.AutoIncrementInfo),
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
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
							AutoIncrementInfo:   new(model.AutoIncrementInfo),
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
							IsForeign:           true,
							JointTable: &model.TableProperties{
								To:             "id",
								Table:          "customer",
								OnDelete:       "CASCADE",
								ConstraintName: "c_tweet_customer_id",
							},
						},
						"order_dates": &model.FieldType{
							FieldName:         "order_dates",
							AutoIncrementInfo: new(model.AutoIncrementInfo),
							IsList:            true,
							Kind:              model.TypeDateTime,
							TypeIDSize:        model.SQLTypeIDSize,
							IsLinked:          true,
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
			Data: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
					Table:   "tweet",
					DbAlias: "mongo",
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
	}

	s := Init("chicago", &crud.Module{})
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r, err := s.Parser(testCase.Data)
			if (err != nil) != testCase.IsErrExpected {
				t.Errorf("\n Schema.parseSchema() error = expected error-%v,got error-%v", testCase.IsErrExpected, err)
				return
			}
			if arr := deep.Equal(r, testCase.schema); len(arr) > 0 {
				t.Errorf("generateInspection() differences = %v", arr)
			}
		})
	}
}
