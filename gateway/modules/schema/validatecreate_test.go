package schema

import (
	"reflect"
	"testing"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

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
var Parsedata = config.Crud{
	"mongo": &config.CrudStub{
		Collections: map[string]*config.TableRule{
			"tweet": {
				Schema: testQueries,
			},
			"test": {
				Schema: testQueries,
			},
			"location": {
				Schema: testQueries,
			},
		},
	},
}

func TestSchema_ValidateCreateOperation(t *testing.T) {

	testCases := []struct {
		dbName, coll, name string
		value              model.CreateRequest
		IsErrExpected      bool
	}{
		{
			dbName:        "sqlserver",
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
			dbName:        "mongo",
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
			dbName:        "mongo",
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
			dbName:        "mongo",
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
			dbName:        "mongo",
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
			dbName:        "mongo",
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

	s := Init(&crud.Module{})
	err := s.parseSchema(Parsedata)
	if err != nil {
		t.Errorf("Error while parsing schema-%v", err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := s.ValidateCreateOperation(testCase.dbName, testCase.coll, &testCase.value)
			if (err != nil) != testCase.IsErrExpected {
				t.Errorf("\n ValidateCreateOperation() error = expected error-%v,got-%v)", testCase.IsErrExpected, err)
			}
		})
	}
}
func TestSchema_SchemaValidate(t *testing.T) {
	testCases := []struct {
		coll, name    string
		Document      map[string]interface{}
		IsErrExpected bool
		IsSkipable    bool
	}{{
		coll:          "test",
		name:          "inserting value for linked field",
		IsErrExpected: true,
		IsSkipable:    true,
		Document: map[string]interface{}{
			"person": "12PM",
		},
	},
		{
			coll:          "tweet",
			name:          "required field not present",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"latitude": "12PM",
			},
		},
		{
			coll:          "tweet",
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
	s := Init(&crud.Module{})
	err := s.parseSchema(Parsedata)
	if err != nil {
		t.Errorf("Error while parsing schema:%v", err)
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := s.SchemaValidator(testCase.coll, s.SchemaDoc["mongo"][testCase.coll], testCase.Document)
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
		coll, name    string
		Document      map[string]interface{}
		result        interface{}
		IsErrExpected bool
		IsSkipable    bool
	}{{
		coll:          "tweet",
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
			name:          "integer value for string field",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"text": 12,
			},
		},
		{
			coll:          "tweet",
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
			name:          "string value for datetime field",
			IsErrExpected: true,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"createdAt": "12",
			},
		},
		{
			coll:          "tweet",
			name:          "valid datetime field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"createdAt": "1999-10-19T11:45:26.371Z",
			},
		},
		{
			coll:          "tweet",
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
			name:          "string value for integer field",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"exp": "12",
			},
		},
		{
			coll:          "tweet",
			name:          "float value for string field",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"text": 12.5,
			},
		},
		{
			coll:          "tweet",
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
			name:          "invalid boolean value",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"age": true,
			},
		},
		{
			coll:          "tweet",
			name:          "float value for datetime field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"createdAt": 12.5,
			},
		},
		{
			coll:          "tweet",
			name:          "invalid map string interface",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"exp": map[string]interface{}{"years": 10},
			},
		},
		{
			coll:          "tweet",
			name:          "valid map string interface",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"event": map[string]interface{}{"name": "suyash"},
			},
		},
		{
			coll:          "tweet",
			name:          "float value for integer field",
			IsErrExpected: false,
			IsSkipable:    true,
			Document: map[string]interface{}{
				"exp": 12.5,
			},
		},
		{
			coll:          "tweet",
			name:          "valid interface value",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"event": []interface{}{1, 2},
			},
		},
		{
			coll:          "tweet",
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
			name:          "invalid interface value",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"text": []interface{}{1, 2},
			},
		},
		{
			coll:          "tweet",
			name:          "no matching type",
			IsErrExpected: true,
			Document: map[string]interface{}{
				"age": int32(6),
			},
		},
		{
			coll:          "tweet",
			name:          "valid JSON TYPE",
			IsErrExpected: false,
			result:        "{\"name\":\"goku\",\"sage\":\"cell\"}",
			Document: map[string]interface{}{
				"spec": map[string]interface{}{"name": "goku", "sage": "cell"},
			},
		},
		{
			coll:          "tweet",
			name:          "in valid JSON TYPE",
			IsErrExpected: true,
			IsSkipable:    true,
			result:        "{\"name\":\"goku\",\"sage\":\"cell\"}",
			Document: map[string]interface{}{
				"spec": 1,
			},
		},
	}
	s := Init(&crud.Module{})
	err := s.parseSchema(Parsedata)
	if err != nil {
		t.Errorf("Error while parsing schema:%v", err)
	}
	r := s.SchemaDoc["mongo"]["tweet"]
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for key, value := range testCase.Document {
				retval, err := s.checkType(testCase.coll, value, r[key])
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
