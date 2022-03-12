package schema

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
)

func TestParseSchema(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		schema        model.DBSchemas
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
			schema: model.DBSchemas{
				"mongo": model.CollectionSchemas{
					"tweet": model.FieldSchemas{
						"ID": &model.FieldType{
							FieldName:      "ID",
							Kind:           model.TypeID,
							TypeIDSize:     model.DefaultCharacterSize,
							IsPrimary:      true,
							PrimaryKeyInfo: &model.TableProperties{},
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
							TypeIDSize:          model.DefaultCharacterSize,
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
			schema: model.DBSchemas{
				"mongo": model.CollectionSchemas{
					"tweet": model.FieldSchemas{
						"ID": &model.FieldType{
							FieldName:      "ID",
							Kind:           model.TypeID,
							TypeIDSize:     model.DefaultCharacterSize,
							IsPrimary:      true,
							PrimaryKeyInfo: &model.TableProperties{},
						},
						"age": &model.FieldType{
							FieldName: "age",
							Kind:      model.TypeFloat,
						},
						"amount": &model.FieldType{
							FieldName:           "amount",
							Kind:                model.TypeDecimal,
							IsFieldTypeRequired: true,
							Args: &model.FieldArgs{
								Scale:     5,
								Precision: 10,
							},
						},
						"role": &model.FieldType{
							FieldName:           "role",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
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
							Args: &model.FieldArgs{
								Precision: model.DefaultDateTimePrecision,
							},
						},
						"updatedAt": &model.FieldType{
							FieldName:   "updatedAt",
							Kind:        model.TypeDateTime,
							IsUpdatedAt: true,
							Args: &model.FieldArgs{
								Precision: model.DefaultDateTimePrecision,
							},
						},
						"first_name": &model.FieldType{
							FieldName:           "first_name",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
							IndexInfo: []*model.TableProperties{
								{
									IsIndex: true,
									Group:   "user_name",
									Order:   1,
									Sort:    "asc",
									Field:   "first_name",
								},
							},
						},
						"name": &model.FieldType{
							FieldName:           "name",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
							IndexInfo: []*model.TableProperties{
								{
									IsUnique: true,
									Group:    "user_name",
									Order:    1,
									Sort:     "asc",
									Field:    "name",
								},
							},
						},
						"customer_id": &model.FieldType{
							FieldName:           "customer_id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
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
							Args: &model.FieldArgs{
								Precision: model.DefaultDateTimePrecision,
							},
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
						 amount: Decimal! @args(precision:10,scale:5)
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r, err := Parser(testCase.Data)
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

func TestSchema_SchemaValidate(t *testing.T) {
	testCases := []struct {
		dbAlias, dbType, coll, name string
		Document                    map[string]interface{}
		IsErrExpected               bool
		IsSkipable                  bool
	}{
		{
			coll:          "test",
			dbAlias:       "mongo",
			dbType:        string(model.Mongo),
			name:          "inserting value for linked field",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"person": "12PM",
			},
		},
		{
			coll:          "tweet",
			dbAlias:       "mongo",
			dbType:        string(model.Mongo),
			name:          "required field not present",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"latitude": "12PM",
			},
		},
		{
			coll:          "tweet",
			dbAlias:       "mongo",
			dbType:        string(model.Mongo),
			name:          "field having directive createdAt",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"id":        "1234",
				"createdAt": "2019-12-23 05:52:16.5366853 +0000 UTC",
				"age":       12.5,
			},
		},
		{
			coll:          "tweet",
			dbAlias:       "mongo",
			dbType:        string(model.Mongo),
			name:          "valid field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"text": "12PM",
				"age":  12.65,
			},
		},
		{
			coll:          "location",
			dbAlias:       "mongo",
			dbType:        string(model.Mongo),
			name:          "valid field with mutated doc",
			IsErrExpected: false,
			IsSkipable:    false,
			Document: map[string]interface{}{
				"id":        "1234",
				"latitude":  21.3,
				"longitude": 64.5,
			},
		},
	}

	schemaDoc, err := Parser(Parsedata)
	if err != nil {
		t.Errorf("unable to generate test queries - (%v)", err)
		return
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := Validate(context.Background(), testCase.dbAlias, testCase.dbType, testCase.coll, schemaDoc["mongo"][testCase.coll], testCase.Document)
			if (err != nil) != testCase.IsErrExpected {
				t.Errorf("\n SchemaValidateOperation() error : expected error-%v, got-%v)", testCase.IsErrExpected, err)
			}
			if !testCase.IsSkipable {
				if !reflect.DeepEqual(result, testCase.Document) {
					t.Errorf("\n SchemaValidateOperation() error : got  %v,expected %v)", result, testCase.Document)
				}
			}
		})
	}
}

func TestSchema_CheckType(t *testing.T) {
	testCases := []struct {
		dbAlias       string
		coll          string
		name          string
		dbType        string
		Document      map[string]interface{}
		result        interface{}
		IsErrExpected bool
		IsSkipable    bool
	}{
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			name:          "integer value for float field",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        float64(12),
			Document: map[string]interface{}{
				"age": 12,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "integer value for string field",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"text": 12,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "integer value for datetime field",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        time.Unix(int64(12)/1000, 0),
			Document: map[string]interface{}{
				"createdAt": 12,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "string value for datetime field",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"createdAt": "12",
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid datetime field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"createdAt": "1999-10-19T11:45:26.371Z",
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid integer value",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        12,
			Document: map[string]interface{}{
				"exp": 12,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid string value",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        "12",
			Document: map[string]interface{}{
				"text": "12",
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid float value",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        12.5,
			Document: map[string]interface{}{
				"age": 12.5,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "string value for integer field",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"exp": "12",
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "float value for string field",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"text": 12.5,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid boolean value",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        true,
			Document: map[string]interface{}{
				"isMale": true,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "invalid boolean value",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"age": true,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "float value for datetime field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"createdAt": 12.5,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "invalid map string interface",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"exp": map[string]interface{}{"years": 10},
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid map string interface",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"event": map[string]interface{}{"name": "suyash"},
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "float value for integer field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"exp": 12.5,
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid interface value",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"event": []interface{}{1, 2},
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid interface value for matching field (event)",
			IsErrExpected: false,
			IsSkipable:    false,
			result:        map[string]interface{}{"name": "Suyash"},
			Document: map[string]interface{}{
				"event": map[string]interface{}{"name": "Suyash"},
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "invalid interface value",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"text": []interface{}{1, 2},
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "no matching type",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"age": int32(6),
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "valid JSON TYPE",
			IsErrExpected: false,
			result:        map[string]interface{}{"name": "goku", "sage": "cell"},
			Document: map[string]interface{}{
				"spec": map[string]interface{}{"name": "goku", "sage": "cell"},
			},
		},
		{
			coll:          "tweet",
			dbType:        string(model.Mongo),
			dbAlias:       "mongo",
			name:          "in valid JSON TYPE",
			IsErrExpected: true,
			IsSkipable:    true,
			result:        "{\"name\":\"goku\",\"sage\":\"cell\"}",
			Document: map[string]interface{}{
				"spec": 1,
			},
		},
	}

	schemaDoc, err := Parser(Parsedata)
	if err != nil {
		t.Errorf("Unable to genereate test data - (%v)", err)
		return
	}
	r := schemaDoc["mongo"]["tweet"]

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for key, value := range testCase.Document {
				retval, err := checkType(context.Background(), testCase.dbAlias, testCase.dbType, testCase.coll, value, r[key])
				if (err != nil) != testCase.IsErrExpected {
					t.Errorf("\n CheckType() error = Expected error-%v, got-%v)", testCase.IsErrExpected, err)
				}
				if !testCase.IsSkipable {
					if !reflect.DeepEqual(retval, testCase.result) {
						t.Errorf("\n CheckType() error = Expected return value-%v,got-%v)", testCase.result, retval)
					}
				}
			}

		})
	}
}
