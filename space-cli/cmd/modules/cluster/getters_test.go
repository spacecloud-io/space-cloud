package cluster

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func TestGetClusterConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name              string
		transportMockArgs []mockArgs
		want              []*model.SpecObject
		wantErr           bool
	}{
		{
			name: "Successful test",
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/cluster", map[string]string{}, new(resp)},
					paramsReturned: []interface{}{nil, resp{
						Result: map[string]interface{}{
							"letsEncryptEmail": "info@gmail.com",
							"enableTelemetry":  false,
						},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/cluster",
					Type: "cluster-config",
					Meta: map[string]string{},
					Spec: map[string]interface{}{"clusterConfig": map[string]interface{}{"letsEncryptEmail": "info@gmail.com", "enableTelemetry": false}},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function return Error",
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/cluster", map[string]string{}, new(resp)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{
							"letsEncryptEmail": "",
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
			got, err := GetClusterConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetClusterConfig() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetClusterConfig() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}

func TestGetIntegration(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name              string
		transportMockArgs []mockArgs
		want              []*model.SpecObject
		wantErr           bool
	}{
		{
			name: "Successful test",
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/integrations", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{nil, model.Response{
						Result: []interface{}{map[string]interface{}{
							"id":      "integration1",
							"name":    "name1",
							"version": "v1",
						}},
					}},
				},
			},
			want: []*model.SpecObject{
				{
					API:  "/v1/config/integrations",
					Type: "integrations",
					Meta: map[string]string{},
					Spec: map[string]interface{}{"integration": map[string]interface{}{"id": "integration1", "name": "name1", "version": "v1"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Get function return Error",
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args:   []interface{}{"GET", "/v1/config/integrations", map[string]string{}, new(model.Response)},
					paramsReturned: []interface{}{fmt.Errorf("cannot unmarshal"), model.Response{
						Result: []interface{}{map[string]interface{}{}},
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
			got, err := GetIntegration()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntegration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetIntegration() len= %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf("GetIntegration() v = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
