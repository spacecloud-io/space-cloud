package schema

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

var queries = `
 type tweet {
 	id: ID @id
 	createdAt: DateTime@createdAt
 	text: String
 	owner: [String]
 	location: location@foreign
	age : Float!
	isMale:Boolean
	exp:Integer
	event:event_logs
	person : sharad @link(table:sharad, from:Name, to:isMale)
   }

   type test {
	id : ID @id
	person : sharad @link(table:sharad, from:Name, to:isMale)
   }

   type location {
 	id: ID! @id
 	latitude: Float
 	longitude: Float

   }
   type sharad {
 	  Name : String!
 	  Surname : String!
 	  age : Integer!
 	  isMale : Boolean!
 	  dob : DateTime@createdAt
   }
   type event_logs {
 	id: Integer
 	name: String
   }
 `
var Parsedata = config.Crud{
	"mongo": &config.CrudStub{
		Collections: map[string]*config.TableRule{
			"tweet": &config.TableRule{
				Schema: queries,
			},
			"test": &config.TableRule{
				Schema: queries,
			},
			"location": &config.TableRule{
				Schema: queries,
			},
		},
	},
}

func TestSchema_ValidateCreateOperation(t *testing.T) {

	tdd := []struct {
		dbName, coll, description string
		value                     model.CreateRequest
		want                      error
	}{
		{
			dbName: "sqlserver",
			coll:   "tweet",
			want:   errors.New("No db was found named sqlserver"),
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"male": true,
				},
			},
		},
		{
			dbName: "mongo",
			coll:   "twee",
			want:   nil,
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"male": true,
				},
			},
		},
		{
			dbName: "mongo",
			coll:   "tweet",
			want:   errors.New("required field age from collection tweet not present in request"),
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"male": true,
				},
			},
		},
		{
			dbName: "mongo",
			coll:   "tweet",
			want:   errors.New("invalid document provided for collection (mongo:tweet)"),
			value: model.CreateRequest{
				Document: []interface{}{
					"text", "12PM",
				},
			},
		},
		{
			dbName: "mongo",
			coll:   "tweet",
			want:   errors.New("required field age from collection tweet not present in request"),
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"isMale": true,
				},
			},
		},
		{
			dbName: "mongo",
			coll:   "location",
			want:   nil,
			value: model.CreateRequest{
				Document: map[string]interface{}{
					"location": "21",
				},
			},
		},
	}

	temp := crud.Module{}
	s := Init(&temp, false)
	err := s.parseSchema(Parsedata)
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range tdd {
		t.Run(val.description, func(t *testing.T) {
			err := s.ValidateCreateOperation(val.dbName, val.coll, &val.value)
			if !reflect.DeepEqual(err, val.want) {
				t.Errorf("\n SchemaValidateCreateOperation() error = (%v,%v)", err, val.want)
			}
		})
	}

}
func TestSchema_SchemaValidate(t *testing.T) {
	td := []struct {
		coll, description string
		Document          map[string]interface{}
		want              error
	}{{
		coll:        "test",
		description: "inserting value for linked field",
		want:        errors.New("cannot insert value for a linked field person"),
		Document: map[string]interface{}{
			"person": "12PM",
		},
	},
		{
			coll:        "tweet",
			description: "required field not present",
			want:        errors.New("required field age from collection tweet not present in request"),
			Document: map[string]interface{}{
				"latitude": "12PM",
			},
		},
		{
			coll:        "tweet",
			description: "field having directive createdAt",
			want:        errors.New("required field age from collection tweet not present in request"),
			Document: map[string]interface{}{
				"createdAt": "12PM",
			},
		},
		{
			coll:        "tweet",
			description: "valid field",
			want:        errors.New("required field age from collection tweet not present in request"),
			Document: map[string]interface{}{
				"text": "12PM",
			},
		},
	}
	temp := crud.Module{}
	s := Init(&temp, false)
	err := s.parseSchema(Parsedata)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range td {
		t.Run(v.description, func(t *testing.T) {
			_, err := s.schemaValidator(v.coll, s.SchemaDoc["mongo"][v.coll], v.Document)
			if !reflect.DeepEqual(err, v.want) {
				t.Errorf("\n SchemaValidateCreateOperation() error = (%v,%v)", err, v.want)
			}
		})
	}
}

func TestSchema_CheckType(t *testing.T) {
	td := []struct {
		coll, description string
		Document          map[string]interface{}
		want              error
	}{{
		coll:        "tweet",
		description: "integer value for float field",
		want:        nil,
		Document: map[string]interface{}{
			"age": 12,
		},
	},
		{
			coll:        "tweet",
			description: "integer value for string field",
			want:        errors.New("invalid type received for field text in collection tweet - wanted String got Integer"),
			Document: map[string]interface{}{
				"text": 12,
			},
		},
		{
			coll:        "tweet",
			description: "integer value for datetime field",
			want:        nil,
			Document: map[string]interface{}{
				"createdAt": 12,
			},
		},
		{
			coll:        "tweet",
			description: "string value for datetime field",
			want:        errors.New("invalid datetime format recieved for field createdAt in collection tweet - use RFC3339 fromat"),
			Document: map[string]interface{}{
				"createdAt": "12",
			},
		},
		{
			coll:        "tweet",
			description: "valid datetime field",
			want:        nil,
			Document: map[string]interface{}{
				"createdAt": "1999-10-19T11:45:26.371Z",
			},
		},
		{
			coll:        "tweet",
			description: "valid integer value",
			want:        nil,
			Document: map[string]interface{}{
				"exp": 12,
			},
		},
		{
			coll:        "tweet",
			description: "valid string value",
			want:        nil,
			Document: map[string]interface{}{
				"text": "12",
			},
		},
		{
			coll:        "tweet",
			description: "valid float value",
			want:        nil,
			Document: map[string]interface{}{
				"age": 12.5,
			},
		},
		{
			coll:        "tweet",
			description: "string value for integer field",
			want:        errors.New("invalid type received for field exp in collection tweet - wanted Integer got String"),
			Document: map[string]interface{}{
				"exp": "12",
			},
		},
		{
			coll:        "tweet",
			description: "float value for string field",
			want:        errors.New("invalid type received for field text in collection tweet - wanted String got Float"),
			Document: map[string]interface{}{
				"text": 12.5,
			},
		},
		{
			coll:        "tweet",
			description: "valid boolean value",
			want:        nil,
			Document: map[string]interface{}{
				"isMale": true,
			},
		},
		{
			coll:        "tweet",
			description: "invalid boolean value",
			want:        errors.New("invalid type received for field age in collection tweet - wanted Float got Bool"),
			Document: map[string]interface{}{
				"age": true,
			},
		},
		{
			coll:        "tweet",
			description: "float value for datetime field",
			want:        nil,
			Document: map[string]interface{}{
				"createdAt": 12.5,
			},
		},
		{
			coll:        "tweet",
			description: "invalid map string interface",
			want:        errors.New("invalid type received for field exp in collection tweet"),
			Document: map[string]interface{}{
				"exp": map[string]interface{}{"years": 10},
			},
		},
		{
			coll:        "tweet",
			description: "valid map string interface",
			want:        nil,
			Document: map[string]interface{}{
				"event": map[string]interface{}{"name": "suyash"},
			},
		},
		{
			coll:        "tweet",
			description: "float value for integer field",
			want:        nil,
			Document: map[string]interface{}{
				"exp": 12.5,
			},
		},
		{
			coll:        "tweet",
			description: "valid interface value",
			want:        errors.New("invalid type received for field event in collection tweet - wanted Object got Integer"),
			Document: map[string]interface{}{
				"event": []interface{}{1, 2},
			},
		},
		{
			coll:        "tweet",
			description: "invalid interface value",
			want:        errors.New("invalid type received for field text in collection tweet"),
			Document: map[string]interface{}{
				"text": []interface{}{1, 2},
			},
		},
		{
			coll:        "tweet",
			description: "no matching type",
			want:        errors.New("no matching type found for field age in collection tweet"),
			Document: map[string]interface{}{
				"age": int32(6),
			},
		},
	}
	temp := crud.Module{}
	s := Init(&temp, false)
	err := s.parseSchema(Parsedata)
	if err != nil {
		t.Fatal(err)
	}
	r := s.SchemaDoc["mongo"]["tweet"]
	for _, v := range td {
		t.Run(v.description, func(t *testing.T) {
			for key, value := range v.Document {
				if _, err := s.checkType(v.coll, value, r[key]); err != nil {
					if !reflect.DeepEqual(err, v.want) {
						t.Errorf("\n CheckType() error = (%v,%v,%v)", v.description, err, v.want)
					}
				}
			}

		})
	}
}
