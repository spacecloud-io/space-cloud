package letsencrypt

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

func TestGetLetsEncryptDomain(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project     string
		commandName string
		params      map[string]string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		want              []*model.SpecObject
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			name: "Successful test",
			args: args{
				project:     "myproject",
				commandName: "letsencrypt",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/letsencrypt/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"aesKey":             "mongodb",
							"id":                 "local-admin",
							"name":               "abcd",
							"dockerRegistry":     "space-cloud",
							"ContextTimeGraphQL": 10,
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/letsencrypt/config/{id}",
					Type: "letsencrypt",
					Meta: map[string]string{"project": "myproject", "id": "letsencrypt"},
					Spec: map[string]interface{}{"aesKey": "mongodb", "name": "abcd", "dockerRegistry": "space-cloud", "ContextTimeGraphQL": float64(10)},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "letsencrypt",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/letsencrypt/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"aesKey":             "mongodb",
							"id":                 "local-admin",
							"name":               "abcd",
							"dockerRegistry":     "space-cloud",
							"ContextTimeGraphQL": 10,
						},
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSchema := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockSchema
			got, err := GetLetsEncryptDomain(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLetsEncryptDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetLetsEncryptDomain() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetLetsEncryptDomain() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
