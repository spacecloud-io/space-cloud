package schema

import (
	"testing"
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
					"set": map[string]interface{}{
						"id":        1,
						"createdAt": 986413662654,
						"text":      "heelo",
					},
					"inc": map[string]interface{}{
						"age": 19,
					},
					"push": map[string]interface{}{
						"owner": []interface{}{1, 2, 3},
					},
					"currentDate": map[string]interface{}{
						"createdAt": 16641894861,
					},
				},
			},
		},
	}
	s := Init()
	s.ParseSchema(ParseData)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.ValidateUpdateOperation(tt.args.dbType, tt.args.col, tt.args.updateDoc); err != nil {
				t.Errorf("Schema.ValidateUpdateOperation() error = %v", err)
			}
		})
	}
}
