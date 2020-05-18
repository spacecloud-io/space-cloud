package auth

import (
	"reflect"
	"testing"

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
				mockArgs{
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
			want:    []*model.SpecObject{},
			wantErr: false,
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthProviders() = %v-%v-%v-%v, want %v-%v-%v-%v", got, len(got), reflect.TypeOf(got), cap(got), tt.want, len(tt.want), reflect.TypeOf(tt.want), cap(tt.want))
			}
		})
	}
}
