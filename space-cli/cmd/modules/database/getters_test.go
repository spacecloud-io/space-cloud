package database

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
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
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"dbAlias":           "postgres",
							"col":               "event",
							"isRealtimeEnabled": true,
							"rules": map[string]interface{}{
								"id":   "local-admin",
								"rule": "allow",
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
					method: "MakeHTTPRequest",
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
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"dbAlias":           "postgres",
							"col":               "event_logs",
							"isRealtimeEnabled": true,
							"rules": map[string]interface{}{
								"id":   "local-admin",
								"rule": "allow",
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
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"dbAlias": "postgres",
							"enabled": true,
							"conn":    "connected",
							"type":    "tcp",
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
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"dbAlias": "postgres",
							"enabled": true,
							"conn":    "connected",
							"type":    "tcp",
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
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/schema/mutate", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{
							map[string]interface{}{
								"col":     "subscribers",
								"dbAlias": "db",
								"schema":  "type subscribers { id: ID! @primary name: String!}",
							},
							map[string]interface{}{
								"col":     "default",
								"dbAlias": "db",
								"schema":  "",
							},
							map[string]interface{}{
								"col":     "genres",
								"dbAlias": "db",
								"schema":  "type genres { id: ID! @primary name: String!}",
							},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate",
					Type: "db-schema",
					Meta: map[string]string{"project": "myproject", "col": "subscribers", "dbAlias": "db"},
					Spec: map[string]interface{}{"schema": "type subscribers { id: ID! @primary name: String!}"},
				},
				{
					API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate",
					Type: "db-schema",
					Meta: map[string]string{"project": "myproject", "col": "genres", "dbAlias": "db"},
					Spec: map[string]interface{}{"schema": "type genres { id: ID! @primary name: String!}"},
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
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/database/collections/schema/mutate", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{nil},
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
			got, err := GetDbSchema(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbSchema() error = %v,\n wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetDbSchema() len= %v,\n want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetDbSchema() v = %v,\n want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetDbPreparedQuery(t *testing.T) {
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
		{
			name: "unable to get response",
			args: args{commandName: "db-prepared-query", project: "project", params: map[string]string{"dbAlias": "dbAlias", "id": "prep"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/project/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "prep"}, new(model.Response)},
					paramsReturned: []interface{}{errors.New("unable to unmarshall"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"id": "prep",
						}},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "Got prepared query",
			args: args{commandName: "db-prepared-query", project: "project", params: map[string]string{"dbAlias": "dbAlias", "id": "prep"}},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/projects/project/database/prepared-queries", map[string]string{"dbAlias": "dbAlias", "id": "prep"}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":      "prep",
							"dbAlias": "dbAlias",
							"args":    []interface{}{"args.id"},
							"sql":     "select * from users",
							"rule": map[string]interface{}{
								"rule": "allow",
							},
						}},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/database/{db}/prepared-queries/{id}",
					Type: "db-prepared-query",
					Meta: map[string]string{"project": "project", "db": "dbAlias", "id": "prep"},
					Spec: map[string]interface{}{
						"args": []interface{}{"args.id"},
						"sql":  "select * from users",
						"rule": map[string]interface{}{
							"rule": "allow",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockTransport

			got, err := GetDbPreparedQuery(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbPreparedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetDbPreparedQuery() = %v, want %v", got, tt.want)
				}
			}

			mockTransport.AssertExpectations(t)
		})
	}
}
