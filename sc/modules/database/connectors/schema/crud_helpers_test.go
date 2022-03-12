package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
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
		schemaDoc model.DBSchemas
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
		schemaDoc model.DBSchemas
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeDateTime}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeBoolean}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeBoolean}}}},
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
				schemaDoc: model.DBSchemas{"mysql": model.CollectionSchemas{"table1": model.FieldSchemas{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeBoolean}}}},
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
		schemaDoc                   model.DBSchemas
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
