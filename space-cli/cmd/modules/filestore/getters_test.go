package filestore

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

func TestGetFileStoreConfig(t *testing.T) {
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
				commandName: "filestore-config",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/file-storage/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"enabled":   true,
							"storeType": "filestore",
							"conn":      "http",
							"endpoint":  "/v1/config/projects/myproject/file-storage/config",
							"bucket":    "abcd",
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/file-storage/config/{id}",
					Type: "filestore-config",
					Meta: map[string]string{"project": "myproject", "id": "filestore-config"},
					Spec: map[string]interface{}{"enabled": true, "storeType": "filestore", "conn": "http", "endpoint": "/v1/config/projects/myproject/file-storage/config", "bucket": "abcd"},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns error",
			args: args{
				project:     "myproject",
				commandName: "filestore-config",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/file-storage/config", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"enabled":   true,
							"storeType": "filestore",
							"conn":      "http",
							"endpoint":  "/v1/config/projects/myproject/file-storage/config",
							"bucket":    "abcd",
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
			got, err := GetFileStoreConfig(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileStoreConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetFileStoreConfig() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetFileStoreConfig() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetFileStoreRule(t *testing.T) {
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
				commandName: "filestore-rule",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/file-storage/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":     "local-admin",
							"prefix": "hello",
							"rule": map[string]interface{}{
								"rule": "allow",
								"eval": "==",
							},
						},
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/projects/{project}/file-storage/rules/{id}",
					Type: "filestore-rule",
					Meta: map[string]string{"project": "myproject", "id": "local-admin"},
					Spec: map[string]interface{}{"prefix": "hello", "rule": map[string]interface{}{"rule": "allow", "eval": "=="}},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function returns error",
			args: args{
				project:     "myproject",
				commandName: "filestore-rule",
				params:      map[string]string{},
			},
			transportMockArgs: []mockArgs{
				{
					method: "Get",
					args:   []interface{}{"GET", "/v1/config/projects/myproject/file-storage/rules", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":     "local-admin",
							"prefix": "hello",
							"rule": map[string]interface{}{
								"rule": "allow",
								"eval": "==",
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
			got, err := GetFileStoreRule(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileStoreRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetFileStoreRule() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetFileStoreRule() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
