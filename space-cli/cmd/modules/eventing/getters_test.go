package eventing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

func TestGetEventingTrigger(t *testing.T) {
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
				commandName: "eventing-triggers",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/triggers", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"type":    "mongodb",
							"retries": 2,
							"timeout": 10,
							"id":      "local-admin",
							"url":     "/v1/config/projects/myproject/eventing/triggers",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/eventing/triggers/{id}",
					Type: "eventing-triggers",
					Meta: map[string]string{"id": "local-admin", "project": "myproject"},
					Spec: map[string]interface{}{"type": "mongodb", "retries": float64(2), "timeout": float64(10), "url": "/v1/config/projects/myproject/eventing/triggers"},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "eventing-triggers",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/triggers", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"type":    "mongodb",
							"retries": 2,
							"timeout": 10,
							"id":      "local-admin",
							"url":     "/v1/config/projects/myproject/eventing/triggers",
							"options": map[string]string{
								"and": "if",
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
			got, err := GetEventingTrigger(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventingTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetEventingTrigger() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !cmp.Equal(v, tt.want[i]) {
					t.Errorf("GetEventingTrigger() v = %v want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetEventingConfig(t *testing.T) {
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
				commandName: "eventing-config",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/config", map[string]string{}, new(model.Response)},
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
					API:  "/v1/config/projects/{project}/eventing/config/{id}",
					Type: "eventing-config",
					Meta: map[string]string{"project": "myproject", "id": "eventing-config"},
					Spec: map[string]interface{}{"aesKey": "mongodb", "id": "local-admin", "name": "abcd", "dockerRegistry": "space-cloud", "ContextTimeGraphQL": float64(10)},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "eventing-config",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/config", map[string]string{}, new(model.Response)},
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
			got, err := GetEventingConfig(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventingConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetEventingConfig() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !cmp.Equal(*v, *tt.want[i]) {
					t.Errorf("GetEventingConfig() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetEventingSchema(t *testing.T) {
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
				commandName: "eventing-schema",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/schema", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":     "local-admin",
							"schema": "mongodb",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/eventing/schema/{id}",
					Type: "eventing-schema",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{"schema": "mongodb"},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns Error",
			args: args{
				project:     "myproject",
				commandName: "eventing-schema",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/schema", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarchal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":     "local-admin",
							"schema": "mongodb",
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
			got, err := GetEventingSchema(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventingSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetEventingSchema() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetEventingSchema() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetEventingSecurityRule(t *testing.T) {
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
				commandName: "eventing-rule",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":   "local-admin",
							"rule": "date",
							"eval": "==",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/eventing/rules/{id}",
					Type: "eventing-rule",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{"rule": "date", "eval": "=="},
				},
			},
			wantErr: false,
		},
		{
			name: "Get Function returns Error",
			args: args{
				project:     "myproject",
				commandName: "eventing-rule",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/eventing/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":   "local-admin",
							"rule": "date",
							"eval": "==",
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
			got, err := GetEventingSecurityRule(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventingSecurityRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetEventingSecurityRule() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetEventingSecurityRule() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
