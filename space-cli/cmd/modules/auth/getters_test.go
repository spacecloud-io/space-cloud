package auth

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func TestGetAuthProviders(t *testing.T) {

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
	type provider struct {
		ID      string `json:"id" yaml:"id"`
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Secret  string `json:"secret" yaml:"secret"`
	}
	tests := []struct {
		name           string
		args           args
		schemaMockArgs []mockArgs
		url            map[string]interface{}
		want           []*model.SpecObject
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				project:     "myproject",
				commandName: "auth-providers",
				params:      map[string]string{},
			},
			schemaMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/user-management/provider", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{provider{
							ID:      "local-admin",
							Enabled: true,
							Secret:  "hello",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/user-management/provider/{id}",
					Type: "auth-providers",
					Meta: map[string]string{"id": "local-admin", "project": "myproject"},
					Spec: map[string]interface{}{"enabled": true, "secret": "hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "test-2",
			args: args{
				project:     "myproject",
				commandName: "auth-providers",
				params:      map[string]string{},
			},
			schemaMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/user-management/provider", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{provider{
							ID:      "local-admin",
							Enabled: true,
							Secret:  "hello",
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
			mockSchema := utils.MocketAuthProviders{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			utils.Client = &mockSchema
			got, err := GetAuthProviders(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetAuthProviders() len= %v, want %v", len(got), len(tt.want))
			} else if len(got) != 0 {
				i := 0
				for i < len(got) {
					if cmp.Equal(*got[i], *tt.want[0]) {
						return
					}
					i = i + 1
				}
				t.Errorf("GetAuthProviders() = %v, want %v", got, tt.want)
			}
			// for {
			// 	i := 0
			// 	t.Errorf("i=%v and len=%v", i, len(got))
			// 	for i < len(got) {
			// 		if cmp.Equal(*got[i], tt.want) {
			// 			t.Errorf("for if")
			// 			i = i + 1
			// 			return
			// 		}
			// 	}
			// 	t.Errorf("GetAuthProviders() = %v, want %v", got, tt.want)
			// 	return
			// }
		})
	}
}
