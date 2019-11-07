package schema

import (
	"testing"

	"github.com/spaceuptech/space-cloud/modules/crud"
)

func TestSchema_ValidateUpdateOperation(t *testing.T) {

	type args struct {
		dbType    string
		col       string
		updateDoc map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Default Test case",
			args: args{
				dbType: "mongo",
				col:    "tweet",
				updateDoc: map[string]interface{}{
					"$set": map[string]interface{}{
						"id":        1,
						"createdAt": 986413662654,
						"text":      "heelo",
					},
					"$inc": map[string]interface{}{
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
	}
	temp := crud.Module{}
	s := Init(&temp)
	s.ParseSchema(ParseData)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.ValidateUpdateOperation(tt.args.dbType, tt.args.col, tt.args.updateDoc); err != nil {
				t.Errorf("\n Schema.ValidateUpdateOperation() error = %v", err)
			}
		})
	}
}
