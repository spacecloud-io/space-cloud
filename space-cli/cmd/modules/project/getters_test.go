package project

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

func TestGetProjectConfig(t *testing.T) {
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
			name: "Successful test ",
			args: args{
				project:     "myproject",
				commandName: "project",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"aesKey":             "local-admin",
							"contextTimeGraphQL": 10,
							"DockerRegistry":     "gateway",
							"secrets":            "abcd",
							"name":               "space-cloud",
							"id":                 "local-admin",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}",
					Type: "project",
					Meta: map[string]string{"project": "local-admin"},
					Spec: map[string]interface{}{"aesKey": "local-admin",
						"contextTimeGraphQL": float64(10),
						"DockerRegistry":     "gateway",
						"secrets":            "abcd",
						"name":               "space-cloud",
						"id":                 "local-admin"},
				},
			},
			wantErr: false,
		},
		{
			name: "project not provided",
			args: args{
				project:     "",
				commandName: "project",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/*", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"aesKey":             "local-admin",
							"contextTimeGraphQL": 10,
							"name":               "space-cloud",
							"id":                 "local-admin",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}",
					Type: "project",
					Meta: map[string]string{"project": "local-admin"},
					Spec: map[string]interface{}{"aesKey": "local-admin",
						"contextTimeGraphQL": float64(10),
						"name":               "space-cloud",
						"id":                 "local-admin"},
				},
			},
			wantErr: false,
		},
		{
			name: "project not provided id provided in prams",
			args: args{
				project:     "",
				commandName: "project",
				params:      map[string]string{"id": "myproject"},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject", map[string]string{"id": "myproject"}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"aesKey":             "local-admin",
							"contextTimeGraphQL": 10,
							"name":               "space-cloud",
							"id":                 "local-admin",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}",
					Type: "project",
					Meta: map[string]string{"project": "local-admin"},
					Spec: map[string]interface{}{"aesKey": "local-admin",
						"contextTimeGraphQL": float64(10),
						"name":               "space-cloud",
						"id":                 "local-admin"},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "project",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"aesKey":             "local-admin",
							"contextTimeGraphQL": 10,
							"name":               "space-cloud",
							"id":                 "local-admin",
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
			got, err := GetProjectConfig(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetProjectConfig() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetProjectConfig() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
