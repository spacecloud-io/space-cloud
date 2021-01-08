package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSchema_CrudPostProcess(t *testing.T) {

	b, err := json.Marshal(model.ReadRequest{Operation: "hello"})
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.Background()), "err", err, nil)
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.Background()), "err", err, nil)
	}

	type args struct {
		dbAlias   string
		dbType    string
		col       string
		schemaDoc model.Type
		result    interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases for mongo
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CrudPostProcess(context.Background(), tt.args.dbAlias, tt.args.dbType, tt.args.col, tt.args.schemaDoc, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.CrudPostProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.result, tt.want) {
				t.Errorf("Schema.CrudPostProcess() tt.args.result = %v, tt.want %v", tt.args.result, tt.want)
			}
		})
	}
}

func returntime(s string) primitive.DateTime {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("invalid string format of datetime (%s)", s), map[string]interface{}{"error": err})
		return primitive.NewDateTimeFromTime(time.Now())
	}
	return primitive.NewDateTimeFromTime(t)
}
func TestSchema_AdjustWhereClause(t *testing.T) {

	type args struct {
		dbAlias   string
		dbType    model.DBType
		schemaDoc model.Type
		col       string
		find      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "db is not mongo",
			args: args{
				dbAlias:   "mysql",
				dbType:    "sql",
				col:       "table1",
				find:      map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "SchemaDoc not provided",
			args: args{
				dbAlias: "mysql",
				dbType:  "mongo",
				col:     "table1",
				find:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "Col not provided",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
				schemaDoc: model.Type{"mysql": model.Collection{}},
			},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "Tableinfo not provided",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{}}},
			},
			want:    map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
			wantErr: false,
		},
		{
			name: "Using param as string",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": "2014-11-12T11:45:26.371Z"},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": returntime("2014-11-12T11:45:26.371Z")},
			wantErr: false,
		},
		{
			name: "Error string format provided",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": "2014-11-12"},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": "2014-11-12"},
			wantErr: true,
		},
		{
			name: "param as map[string]interface{}",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": map[string]interface{}{"time": "2014-11-12T11:45:26.371Z"}},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": returntime("2014-11-12T11:45:26.371Z")}},
			wantErr: false,
		},
		{
			name: "param with map[string]interface{} having value time.time",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": map[string]interface{}{"time": time.Now().Round(time.Second)}},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": primitive.NewDateTimeFromTime(time.Now().Round(time.Second))}},
			wantErr: false,
		},
		{
			name: "Error foramt provided as value to map[string]interface{} ",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": map[string]interface{}{"time": "string"}},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": "string"}},
			wantErr: true,
		},
		{
			name: "Param as time.time",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": time.Now().Round(time.Second)},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": time.Now().Round(time.Second)},
			wantErr: false,
		},
		{
			name: "Param as default",
			args: args{
				dbAlias:   "mysql",
				dbType:    "mongo",
				col:       "table1",
				find:      map[string]interface{}{"col2": 10},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
			},
			want:    map[string]interface{}{"col2": 10},
			wantErr: true,
		},
		{
			name: "SQL server Using param as string",
			args: args{
				dbAlias:   "mysql",
				dbType:    model.SQLServer,
				col:       "table1",
				find:      map[string]interface{}{"col2": true},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeBoolean}}}},
			},
			want:    map[string]interface{}{"col2": 1},
			wantErr: false,
		},
		{
			name: "SQL server param as map[string]interface{}",
			args: args{
				dbAlias:   "mysql",
				dbType:    model.SQLServer,
				col:       "table1",
				find:      map[string]interface{}{"col2": map[string]interface{}{"time": false}},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeBoolean}}}},
			},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": 0}},
			wantErr: false,
		},
		{
			name: "SQL server Error format provided as value to map[string]interface{} ",
			args: args{
				dbAlias:   "mysql",
				dbType:    model.SQLServer,
				col:       "table1",
				find:      map[string]interface{}{"col2": map[string]interface{}{"time": "string"}},
				schemaDoc: model.Type{"mysql": model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeBoolean}}}},
			},
			want:    map[string]interface{}{"col2": map[string]interface{}{"time": "string"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AdjustWhereClause(context.Background(), tt.args.dbAlias, tt.args.dbType, tt.args.col, tt.args.schemaDoc, tt.args.find)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.AdjustWhereClause() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, tt.args.find) {
				t.Errorf("Schema.AdjustWhereClause() find = %v, want %v", tt.args.find, tt.want)
			}
		})
	}
}

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
							FieldName:      "ID",
							Kind:           model.TypeID,
							TypeIDSize:     model.SQLTypeIDSize,
							IsPrimary:      true,
							PrimaryKeyInfo: &model.TableProperties{},
						},
						"age": &model.FieldType{
							FieldName:  "age",
							TypeIDSize: model.SQLTypeIDSize,
							Kind:       model.TypeFloat,
							Args: &model.FieldArgs{
								Precision: model.DefaultPrecision,
								Scale:     model.DefaultScale,
							},
						},
						"spec": &model.FieldType{
							FieldName:  "spec",
							TypeIDSize: model.SQLTypeIDSize,
							Kind:       model.TypeJSON,
						},
						"customer_id": &model.FieldType{
							FieldName:           "customer_id",
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
							FieldName:      "ID",
							Kind:           model.TypeID,
							TypeIDSize:     model.SQLTypeIDSize,
							IsPrimary:      true,
							PrimaryKeyInfo: &model.TableProperties{},
						},
						"age": &model.FieldType{
							FieldName:  "age",
							Kind:       model.TypeFloat,
							TypeIDSize: model.SQLTypeIDSize,
							Args: &model.FieldArgs{
								Precision: model.DefaultPrecision,
								Scale:     model.DefaultScale,
							},
						},
						"role": &model.FieldType{
							FieldName:           "role",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
							IsDefault:           true,
							Default:             "user",
						},
						"spec": &model.FieldType{
							FieldName:  "spec",
							Kind:       model.TypeJSON,
							TypeIDSize: model.SQLTypeIDSize,
						},
						"createdAt": &model.FieldType{
							FieldName:   "createdAt",
							Kind:        model.TypeDateTime,
							TypeIDSize:  model.SQLTypeIDSize,
							IsCreatedAt: true,
							Args: &model.FieldArgs{
								Precision: model.DefaultDateTimePrecision,
							},
						},
						"updatedAt": &model.FieldType{
							FieldName:   "updatedAt",
							Kind:        model.TypeDateTime,
							TypeIDSize:  model.SQLTypeIDSize,
							IsUpdatedAt: true,
							Args: &model.FieldArgs{
								Precision: model.DefaultDateTimePrecision,
							},
						},
						"first_name": &model.FieldType{
							FieldName:           "first_name",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.SQLTypeIDSize,
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
							TypeIDSize:          model.SQLTypeIDSize,
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
							FieldName:  "order_dates",
							IsList:     true,
							Kind:       model.TypeDateTime,
							TypeIDSize: model.SQLTypeIDSize,
							IsLinked:   true,
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

func TestSchema_ValidateUpdateOperation(t *testing.T) {

	var Query = `type tweet {
		id: ID! @primary
		createdAt: DateTime! @createdAt
		text: String
		spec: JSON
		owner: String!
		age : Integer!
		cpi: Float!
		diplomastudent: Boolean! @foreign(table:"shreyas",field:"diploma")
		friends:[String!]!
		update:DateTime @updatedAt
		mentor: shreyas
	}
	type shreyas {
		name:String!
		surname:String!
		diploma:Boolean
	}`

	var dbSchemas = config.DatabaseSchemas{
		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
			Table:   "tweet",
			DbAlias: "mongo",
			Schema:  Query,
		},
	}

	type args struct {
		dbAlias   string
		dbType    string
		col       string
		updateDoc map[string]interface{}
	}
	tests := []struct {
		name          string
		args          args
		IsErrExpected bool
	}{
		{
			name:          "Successful Test case",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id":        "1234",
						"createdAt": 986413662654,
						"text":      "heelo",
						"spec": map[string]interface{}{
							"name": "goku",
							"sage": "boo",
						},
					},
					"$inc": map[string]interface{}{
						"age": 1999,
					},
					"$min": map[string]interface{}{
						"age": 1999,
					},
					"$max": map[string]interface{}{
						"age": 1999,
					},
					"$mul": map[string]interface{}{
						"age": 1999,
					},
					"$push": map[string]interface{}{
						"owner": []interface{}{"hello", "go", "java"},
					},
					"$currentDate": map[string]interface{}{
						"createdAt": 16641894861,
					},
				},
			},
		},
		{
			name:          "Invalid Test case got integer wanted object for json type",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"spec": 123,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded ID got integer",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id": 123,
					},
				},
			},
		},
		{
			name:          "Test case-Nothing to Update",
			IsErrExpected: false,
			args: args{
				dbAlias:   "mongo",
				dbType:    string(model.Mongo),
				col:       "tweet",
				updateDoc: nil,
			},
		},
		{
			name:          "Invalid Test case-$createdAt update operator unsupported",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$createdAt": map[string]interface{}{
						"age": 45,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-expected ID",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"id": "123",
					},
				},
			},
		},
		{
			name:          "Valid Test case-increment operation",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "suyash",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"age": 1234567890,
					},
				},
			},
		},
		{
			name:          "Valid Test case- increment float but kind is integer type",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"age": 6.34,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-document not of type object",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": "suyash",
				},
			},
		},
		{
			name:          "Valid Test case-createdAt",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{

						"createdAt": "2015-11-22T13:57:31.123ZIDE",
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded ID(currentDate)",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{
						"id": 123,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-field not defined in schema",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"location": []interface{}{"hello", "go", "java"},
						"cpi":      7.25,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded string got integer",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"owner": []interface{}{123, 45.64, "java"},
						"cpi":   7.22,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-invalid type for field owner",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"owner": 123,
					},
				},
			},
		},
		{
			name:          "Test Case-Float value",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": 12.33,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded ID got integer",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"id": 721,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-invalid datetime format",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"createdAt": "2015-11-22T13:57:31.123ZI",
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded Integer got String",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": "12",
					},
				},
			},
		},
		{
			name:          "Float value for field createdAt",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"createdAt": 12.13,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded String got Float",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"text": 12.13,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-IsErrExpecteded float got boolean",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": true,
					},
				},
			},
		},
		{
			name:          "Valid Test Case-Boolean",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"diplomastudent": false,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-invalid map string interface",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": map[string]interface{}{"1": 7.2, "2": 8.5, "3": 9.3},
					},
				},
			},
		},
		{
			name:          "Invalid Test case-invalid array interface",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": []interface{}{7.2, "8", 9},
					},
				},
			},
		},
		{
			name:          "set array type for field friends",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friends": []interface{}{"7.2", "8", "9"},
					},
				},
			},
		},
		{
			name:          "Invalid Test case-field not defined in schema",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friend": []interface{}{"7.2", "8", "9"},
					},
				},
			},
		},
		{
			name:          "Invalid Test case-Wanted Object got integer",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"mentor": []interface{}{1, 2},
					},
				},
			},
		},
		{
			name:          "Invalid Test case-no matching type found",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": int32(2),
					},
				},
			},
		},
		{
			name:          "Valid Test Case-set operation",
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": 2,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-field not present in schema",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friend": []map[string]interface{}{{"7.2": "8"}, {"1": 2}},
					},
				},
			},
		},
		{
			name:          "Invalid Test case-invalid boolean field",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"diplomastudent": []interface{}{1, 2, 3},
					},
				},
			},
		},
		{
			name:          "Invalid Test case-unsupported operator",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push1": map[string]interface{}{
						"friends": 4,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-field not present in schema(math)",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"friends1": 4,
					},
				},
			},
		},
		{
			name:          "Invalid Test case-field not present in schema(date)",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{
						"friends1": "4/12/2019",
					},
				},
			},
		},
		{
			name:          "Invalid Test case-document not of type object(math)",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": "age",
				},
			},
		},
		{
			name:          "Invalid Test case-document not of type object(set)",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": "age",
				},
			},
		},
		{
			name:          "Invalid Test case-document not of type object(Date)",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": "15/10/2019",
				},
			},
		},
		{
			name:          "Valid Test case-updatedAt directive involved",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{"update": "15/10/2019"},
				},
			},
		},
		{
			name:          "Invalid Test case-invalid field type in push operation",
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"friends": nil,
					},
				},
			},
		},
		{
			name:          "Invalid Test Case-DB name not present in schema",
			IsErrExpected: true,
			args: args{
				dbAlias: "mysql",
				dbType:  string(model.Mongo),
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id":        123,
						"createdAt": 986413662654,
						"text":      456,
					},
				},
			},
		},
	}

	schemaDoc, err := Parser(dbSchemas)
	if err != nil {
		t.Errorf("unable to genereate test cases - (%v)", err)
		return
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			err := ValidateUpdateOperation(context.Background(), testcase.args.dbAlias, testcase.args.dbType, testcase.args.col, utils.All, testcase.args.updateDoc, nil, schemaDoc)
			if (err != nil) != testcase.IsErrExpected {
				t.Errorf("\n ValidateUpdateOperation() error = expected error-%v, got-%v)", testcase.IsErrExpected, err)
			}

		})
	}
}

var testQueries = `
 type tweet {
 	id: ID @primary
 	createdAt: DateTime @createdAt
 	text: String
 	owner: [String]
 	location: location @foreign
	age : Float!
	isMale:Boolean
	exp:Integer
	spec: JSON
	event: event_logs
	person : sharad @link(table:sharad, from:Name, to:isMale)
   }

   type test {
	id : ID @primary
	person : sharad @link(table:sharad, from:Name, to:isMale)
   }

   type location {
 	id: ID! @primary
 	latitude: Float
 	longitude: Float
   }
   type sharad {
 	  Name : String!
 	  Surname : String!
 	  age : Integer!
 	  isMale : Boolean!
 	  dob : DateTime @createdAt
   }
   type event_logs {
 	name: String
   }
 `
var Parsedata = config.DatabaseSchemas{
	config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "tweet"): &config.DatabaseSchema{
		Table:   "tweet",
		DbAlias: "mongo",
		Schema:  testQueries,
	},
	config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "test"): &config.DatabaseSchema{
		Table:   "test",
		DbAlias: "mongo",
		Schema:  testQueries,
	},
	config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "mongo", "location"): &config.DatabaseSchema{
		Table:   "location",
		DbAlias: "mongo",
		Schema:  testQueries,
	},
}

func TestSchema_ValidateCreateOperation(t *testing.T) {

	testCases := []struct {
		dbAlias, dbType, coll, name string
		schemaDoc                   model.Type
		value                       model.CreateRequest
		IsErrExpected               bool
	}{
		{
			dbAlias:       "sqlserver",
			coll:          "tweet",
			name:          "No db was found named sqlserver",
			IsErrExpected: true,
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"male": true,
				},
			},
		},
		{
			dbAlias:       "mongo",
			coll:          "twee",
			name:          "Collection which does not exist",
			IsErrExpected: false,
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"male": true,
				},
			},
		},
		{
			dbAlias:       "mongo",
			coll:          "tweet",
			name:          "required field age from collection tweet not present in request",
			IsErrExpected: true,
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"male": true,
				},
			},
		},
		{
			dbAlias:       "mongo",
			coll:          "tweet",
			name:          "invalid document provided for collection (mongo:tweet)",
			IsErrExpected: true,
			value: model.CreateRequest{
				Document: []interface{}{
					"text", "12PM",
				},
			},
		},
		{
			dbAlias:       "mongo",
			coll:          "tweet",
			name:          "required field age from collection tweet not present in request",
			IsErrExpected: true,
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"isMale": true,
				},
			},
		},
		{
			dbAlias:       "mongo",
			coll:          "location",
			IsErrExpected: true,
			name:          "Invalid Test Case-document gives extra params",
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"location": 21.5,
					"age":      12.5,
				},
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
			err := ValidateCreateOperation(context.Background(), testCase.dbAlias, testCase.dbType, testCase.coll, schemaDoc, &testCase.value)
			if (err != nil) != testCase.IsErrExpected {
				t.Errorf("\n ValidateCreateOperation() error = expected error-%v,got-%v)", testCase.IsErrExpected, err)
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
			result, err := SchemaValidator(context.Background(), testCase.dbAlias, testCase.dbType, testCase.coll, schemaDoc["mongo"][testCase.coll], testCase.Document)
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
