package schema

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

var Query = `type tweet {
 	id: ID! @id
 	createdAt: DateTime! @createdAt
 	text: String
 	owner: String!
	age : Integer!
	cpi: Float!
	diplomastudent: Boolean!@foreign(table:"shreyas",field:"diploma")
	friends:[String!]!
	mentor: shreyas
   }
   type shreyas{
	   name:String!
	   surname:String!
	   diploma:Boolean
   }
  `

var ParseData = config.Crud{
	"mongo": &config.CrudStub{
		Collections: map[string]*config.TableRule{
			"tweet": &config.TableRule{
				Schema: Query,
			},
		},
	},
}

func TestSchema_ValidateUpdateOperation(t *testing.T) {

	type args struct {
		dbType    string
		col       string
		updateDoc map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		// TODO: Add test cases.
		{
			name: "Successful Test case",
			want: nil,
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id":        "1234",
						"createdAt": 986413662654,
						"text":      "heelo",
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
			name: "Invalid Test case 1",
			want: errors.New("invalid type received for field id in collection tweet - wanted ID got Integer"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id": 123,
					},
				},
			},
		},
		{
			name: "Invalid Test case 2",
			want: errors.New("Nothing to update"),
			args: args{
				dbType:    "mongo",
				col:       "tweet",
				updateDoc: nil,
			},
		},
		{
			name: "Invalid Test case 3",
			want: errors.New("$createdAt update operator is not supported"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$createdAt": map[string]interface{}{
						"age": 45,
					},
				},
			},
		},
		{
			name: "Invalid Test case 4",
			want: errors.New("invalid type received for field id in collection tweet - wanted ID"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"id": "123",
					},
				},
			},
		},
		{
			name: "Valid Test case 2",
			want: nil,
			args: args{
				dbType: "mongo",
				col:    "suyash",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"age": 1234567890,
					},
				},
			},
		},
		{
			name: "Invalid Test case 5",
			want: errors.New("invalid type received for field age in collection tweet - wanted Integer got Float"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"age": 6.34,
					},
				},
			},
		},
		{
			name: "Invalid Test case 6",
			want: errors.New("document not of type object in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push": "suyash",
				},
			},
		},
		{
			name: "Valid Test case 3",
			want: nil,
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{

						"createdAt": "2015-11-22T13:57:31.123ZIDE",
					},
				},
			},
		},
		{
			name: "Invalid Test case 7",
			want: errors.New("invalid type received for field id in collection tweet - wanted ID"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$currentDate": map[string]interface{}{
						"id": 123,
					},
				},
			},
		},
		{
			name: "Invalid Test case 8",
			want: errors.New("field location from collection tweet is not defined in the schema"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"location": []interface{}{"hello", "go", "java"},
						"cpi":      7.25,
					},
				},
			},
		},
		{
			name: "Invalid Test case 9",
			want: errors.New("invalid type received for field owner in collection tweet - wanted String got Integer"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"owner": []interface{}{123, 45.64, "java"},
						"cpi":   7.22,
					},
				},
			},
		},
		{
			name: "Invalid Test case 10",
			want: errors.New("invalid type received for field owner in collection tweet - wanted String got Integer"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"owner": 123,
					},
				},
			},
		},
		{
			name: "Invalid Test case 11",
			want: errors.New("Invalid Type for field age in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": 12.33,
					},
				},
			},
		},
		{
			name: "Invalid Test case 12",
			want: errors.New("invalid type received for field id in collection tweet - wanted ID got Integer"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$inc": map[string]interface{}{
						"id": 721,
					},
				},
			},
		},
		{
			name: "Invalid Test case 13",
			want: errors.New("invalid datetime format recieved for field createdAt in collection tweet - use RFC3339 fromat"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"createdAt": "2015-11-22T13:57:31.123ZI",
					},
				},
			},
		},
		{
			name: "Invalid Test case 14",
			want: errors.New("invalid type received for field age in collection tweet - wanted Integer got String"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": "12",
					},
				},
			},
		},
		{
			name: "Invalid Test case 15",
			want: errors.New("Invalid Type for field createdAt in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"createdAt": 12.13,
					},
				},
			},
		},
		{
			name: "Invalid Test case 16",
			want: errors.New("invalid type received for field text in collection tweet - wanted String got Float"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"text": 12.13,
					},
				},
			},
		},
		{
			name: "Invalid Test case 17",
			want: errors.New("invalid type received for field cpi in collection tweet - wanted Float got Bool"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": true,
					},
				},
			},
		},
		{
			name: "Valid Test Case 4",
			want: nil,
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"diplomastudent": false,
					},
				},
			},
		},
		{
			name: "Invalid Test case 18",
			want: errors.New("invalid type received for field cpi in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": map[string]interface{}{"1": 7.2, "2": 8.5, "3": 9.3},
					},
				},
			},
		},
		{
			name: "Invalid Test case 19",
			want: errors.New("invalid type received for field cpi in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"cpi": []interface{}{7.2, "8", 9},
					},
				},
			},
		},
		{
			name: "Invalid Test Case 20",
			want: errors.New("invalid type received for field friends in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friends": []interface{}{"7.2", "8", "9"},
					},
				},
			},
		},
		{
			name: "Invalid Test case 21",
			want: errors.New("field friend from collection tweet is not defined in the schema"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friend": []interface{}{"7.2", "8", "9"},
					},
				},
			},
		},
		{
			name: "Invalid Test case 22",
			want: errors.New("invalid type received for field mentor in collection tweet - wanted Object got Integer"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"mentor": []interface{}{1, 2},
					},
				},
			},
		},
		{
			name: "Invalid Test case 23",
			want: errors.New("no matching type found for field age in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": int32(2),
					},
				},
			},
		},
		{
			name: "Valid Test Case 5",
			want: nil,
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"age": 2,
					},
				},
			},
		},
		{
			name: "Invalid Test case 24",
			want: errors.New("field friend from collection tweet is not defined in the schema"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"friend": []map[string]interface{}{{"7.2": "8"}, {"1": 2}},
					},
				},
			},
		},
		{
			name: "Invalid Test case 25",
			want: errors.New("invalid type provided for field diplomastudent in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"diplomastudent": []interface{}{1, 2, 3},
					},
				},
			},
		},
		{
			name: "Invalid Test case 26",
			want: errors.New("$push1 update operator is not supported"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push1": map[string]interface{}{
						"friends": 4,
					},
				},
			},
		},
		{
			name: "Invalid Test case 27",
			want: errors.New("invalid type provided for field friends in collection tweet"),
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$push": map[string]interface{}{
						"friends": nil,
					},
				},
			},
		},
		{
			name: "Invalid Test Case 28",
			want: errors.New("mysql is not present in schema"),
			args: args{
				dbType: "mysql",
				col:    "tweet",
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
	temp := crud.Module{}
	s := Init(&temp, false)
	if err := s.parseSchema(ParseData); err != nil {
		t.Errorf("error parsing scheam :: %v", err)
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if err := s.ValidateUpdateOperation(v.args.dbType, v.args.col, v.args.updateDoc); err != nil {
				if !reflect.DeepEqual(err, v.want) {
					t.Errorf("\n ValidateUpdateOperation() error = (%v,%v,%v)", v.name, err, v.want)
				}
			}

		})
	}
}
