package ingress

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

func TestGetIngressRoutes(t *testing.T) {
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
				commandName: "ingress-routes",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": map[string]interface{}{
								"version": "v0.18.0",
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/routing/ingress/{id}",
					Type: "ingress-routes",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{"source": map[string]interface{}{"url": "/v1/config/projects/myproject/routing/ingress"}, "targets": map[string]interface{}{"version": "v0.18.0"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "ingress-routes",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/routing/ingress", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "local-admin",
							"source": map[string]interface{}{
								"url": "/v1/config/projects/myproject/routing/ingress",
							},
							"targets": map[string]interface{}{
								"version": "v0.18.0",
							},
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
			got, err := GetIngressRoutes(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIngressRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetIngressRoutes() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetIngressRoutes() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
