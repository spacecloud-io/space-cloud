package database

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

func TestGetDbRule(t *testing.T) {
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
				commandName: "db-rules",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres-event": map[string]interface{}{
								"isRealtimeEnabled": true,
								"rules": map[string]interface{}{
									"id":   "local-admin",
									"rule": "allow",
								},
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules",
					Type: "db-rules",
					Meta: map[string]string{"project": "myproject", "col": "event", "dbAlias": "postgres"},
					Spec: map[string]interface{}{"isRealtimeEnabled": true, "rules": map[string]interface{}{"id": "local-admin", "rule": "allow"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "db-rules",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres-event": map[string]interface{}{
								"isRealtimeEnabled": true,
								"rules": map[string]interface{}{
									"id":   "local-admin",
									"rule": "allow",
								},
							},
						},
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: true,
		},
		{
			name: "Event_logs as col name",
			args: args{
				project:     "myproject",
				commandName: "db-rules",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres-event_logs": map[string]interface{}{
								"isRealtimeEnabled": true,
								"rules": map[string]interface{}{
									"id":   "local-admin",
									"rule": "allow",
								},
							},
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
			mockSchema := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockSchema
			got, err := GetDbRule(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetDbRule() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetDbRule() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetDbConfig(t *testing.T) {
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
				commandName: "db-config",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres": map[string]interface{}{
								"enabled": true,
								"conn":    "connected",
								"type":    "tcp",
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/database/{dbAlias}/config/{id}",
					Type: "db-config",
					Meta: map[string]string{"project": "myproject", "dbAlias": "postgres", "id": "postgres-config"},
					Spec: map[string]interface{}{"enabled": true, "conn": "connected", "type": "tcp"},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "db-config",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres": map[string]interface{}{
								"enabled": true,
								"conn":    "connected",
								"type":    "tcp",
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
		mockSchema := transport.MocketAuthProviders{}

		for _, m := range tt.transportMockArgs {
			mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
		}

		transport.Client = &mockSchema
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDbConfig(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetDbConfig() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetDbConfig() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetDbSchema(t *testing.T) {
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
				commandName: "db-schema",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/schema/mutate", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres-event": map[string]interface{}{
								"schema": "mongodb",
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate",
					Type: "db-schema",
					Meta: map[string]string{"project": "myproject", "col": "event", "dbAlias": "postgres"},
					Spec: map[string]interface{}{"schema": "mongodb"},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "db-schema",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/schema/mutate", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres-event": map[string]interface{}{
								"schema": "mongodb",
							},
						},
						},
					}},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: true,
		},
		{
			name: "Event_llogs in col",
			args: args{
				project:     "myproject",
				commandName: "db-schema",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/schema/mutate", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"postgres-event_logs": map[string]interface{}{
								"schema": "mongodb",
							},
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
			mockSchema := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockSchema
			got, err := GetDbSchema(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetDbSchema() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetDbSchema() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
