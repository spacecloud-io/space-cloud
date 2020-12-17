package schema

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

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
		col       string
		updateDoc map[string]interface{}
	}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name          string
		args          args
		crudMockArgs  []mockArgs
		IsErrExpected bool
	}{
		{
			name: "Successful Test case",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
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
			name: "Invalid Test case got integer wanted object for json type",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"spec": 123,
					},
				},
			},
		},
		{
			name: "Invalid Test case-IsErrExpecteded ID got integer",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id": 123,
					},
				},
			},
		},
		{
			name: "Test case-Nothing to Update",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias:   "mongo",
				col:       "tweet",
				updateDoc: nil,
			},
		},
		{
			name: "Invalid Test case-$createdAt update operator unsupported",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$createdAt": map[string]interface{}{
						"age": 45,
					},
				},
			},
		},
		{
			name: "Invalid Test case-expected ID",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"id": "123",
					},
				},
			},
		},
		{
			name: "Valid Test case-increment operation",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "suyash",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"age": 1234567890,
					},
				},
			},
		},
		{
			name: "Valid Test case- increment float but kind is integer type",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"age": 6.34,
					},
				},
			},
		},
		{
			name: "Invalid Test case-document not of type object",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": "suyash",
				},
			},
		},
		{
			name: "Valid Test case-createdAt",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{

						"createdAt": "2015-11-22T13:57:31.123ZIDE",
					},
				},
			},
		},
		{
			name: "Invalid Test case-IsErrExpecteded ID(currentDate)",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{
						"id": 123,
					},
				},
			},
		},
		{
			name: "Invalid Test case-field not defined in schema",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
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
			name: "Invalid Test case-IsErrExpecteded string got integer",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
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
			name: "Invalid Test case-invalid type for field owner",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"owner": 123,
					},
				},
			},
		},
		{
			name: "Test Case-Float value",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": 12.33,
					},
				},
			},
		},
		{
			name: "Invalid Test case-IsErrExpecteded ID got integer",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"id": 721,
					},
				},
			},
		},
		{
			name: "Invalid Test case-invalid datetime format",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"createdAt": "2015-11-22T13:57:31.123ZI",
					},
				},
			},
		},
		{
			name: "Invalid Test case-IsErrExpecteded Integer got String",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": "12",
					},
				},
			},
		},
		{
			name: "Float value for field createdAt",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"createdAt": 12.13,
					},
				},
			},
		},
		{
			name: "Invalid Test case-IsErrExpecteded String got Float",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"text": 12.13,
					},
				},
			},
		},
		{
			name: "Invalid Test case-IsErrExpecteded float got boolean",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": true,
					},
				},
			},
		},
		{
			name: "Valid Test Case-Boolean",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"diplomastudent": false,
					},
				},
			},
		},
		{
			name: "Invalid Test case-invalid map string interface",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": map[string]interface{}{"1": 7.2, "2": 8.5, "3": 9.3},
					},
				},
			},
		},
		{
			name: "Invalid Test case-invalid array interface",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": []interface{}{7.2, "8", 9},
					},
				},
			},
		},
		{
			name: "set array type for field friends",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friends": []interface{}{"7.2", "8", "9"},
					},
				},
			},
		},
		{
			name: "Invalid Test case-field not defined in schema",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friend": []interface{}{"7.2", "8", "9"},
					},
				},
			},
		},
		{
			name: "Invalid Test case-Wanted Object got integer",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"mentor": []interface{}{1, 2},
					},
				},
			},
		},
		{
			name: "Invalid Test case-no matching type found",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": int32(2),
					},
				},
			},
		},
		{
			name: "Valid Test Case-set operation",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: false,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": 2,
					},
				},
			},
		},
		{
			name: "Invalid Test case-field not present in schema",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friend": []map[string]interface{}{{"7.2": "8"}, {"1": 2}},
					},
				},
			},
		},
		{
			name: "Invalid Test case-invalid boolean field",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"diplomastudent": []interface{}{1, 2, 3},
					},
				},
			},
		},
		{
			name: "Invalid Test case-unsupported operator",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push1": map[string]interface{}{
						"friends": 4,
					},
				},
			},
		},
		{
			name: "Invalid Test case-field not present in schema(math)",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"friends1": 4,
					},
				},
			},
		},
		{
			name: "Invalid Test case-field not present in schema(date)",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{
						"friends1": "4/12/2019",
					},
				},
			},
		},
		{
			name: "Invalid Test case-document not of type object(math)",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$inc": "age",
				},
			},
		},
		{
			name: "Invalid Test case-document not of type object(set)",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": "age",
				},
			},
		},
		{
			name: "Invalid Test case-document not of type object(Date)",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": "15/10/2019",
				},
			},
		},
		{
			name: "Valid Test case-updatedAt directive involved",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{"update": "15/10/2019"},
				},
			},
		},
		{
			name: "Invalid Test case-invalid field type in push operation",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mongo",
				col:     "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"friends": nil,
					},
				},
			},
		},
		{
			name: "Invalid Test Case-DB name not present in schema",
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"mongo"},
					paramsReturned: []interface{}{"mongo"},
				},
			},
			IsErrExpected: true,
			args: args{
				dbAlias: "mysql",
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

	for _, testcase := range tests {
		mockCrud := &mockCrudSchemaInterface{}
		for _, m := range testcase.crudMockArgs {
			mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
		}
		s := Init("chicago", mockCrud)

		err := s.SetConfig(dbSchemas, "myproject")
		if err != nil {
			t.Errorf("Error while parsing schema:%v", err)
		}

		t.Run(testcase.name, func(t *testing.T) {
			err := s.ValidateUpdateOperation(context.Background(), testcase.args.dbAlias, testcase.args.col, utils.All, testcase.args.updateDoc, nil)
			if (err != nil) != testcase.IsErrExpected {
				t.Errorf("\n ValidateUpdateOperation() error = expected error-%v, got-%v)", testcase.IsErrExpected, err)
			}

		})
	}
}
