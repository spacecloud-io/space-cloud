package k8s

import (
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
)

func TestAddOrUpdateConfig(t *testing.T) {
	// Source objects
	config1 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "core.space-cloud.io/v1alpha1",
			"kind":       "CompiledGraphqlSource",
			"metadata": map[string]interface{}{
				"name":      "get-all-todos",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"graphql": map[string]interface{}{
					"query": "query MyQuery($_eq: Int = \"\") {\n  hasura_persons(where: {id: {_eq: $_eq}}) {\n    age\n    id\n    name\n  }\n}\n",
				},
				"http": map[string]interface{}{
					"method": "GET",
					"url":    "/v1/todos",
				},
			},
		},
	}

	config2 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "core.space-cloud.io/v1alpha1",
			"kind":       "CompiledGraphqlSource",
			"metadata": map[string]interface{}{
				"name":      "opa-test",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"graphql": map[string]interface{}{
					"query": "query MyQuery @auth {\n  hasura_persons {\n    name\n    id\n    age\n  }\n}\n",
				},
				"http": map[string]interface{}{
					"method": "GET",
					"url":    "/test-opa",
				},
				"plugins": []interface{}{
					map[string]interface{}{
						"driver": "auth_opa",
						"name":   "only-admins",
						"params": map[string]interface{}{
							"ref": map[string]interface{}{
								"name": "basicrule",
							},
						},
					},
				},
			},
		},
	}

	config3 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "core.space-cloud.io/v1alpha1",
			"kind":       "GraphqlSource",
			"metadata": map[string]interface{}{
				"name":      "hasura",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"source": map[string]interface{}{
					"url": "http://localhost:8080/v1/graphql",
				},
			},
		},
	}

	config4 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "core.space-cloud.io/v1alpha1",
			"kind":       "CompiledGraphqlSource",
			"metadata": map[string]interface{}{
				"name":      "get-all-todos",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"graphql": map[string]interface{}{
					"query": "query MyQuery($_eq: Int = \"\") {\n  hasura_persons(where: {id: {_eq: $_eq}}) {\n    age\n    id\n    name\n  }\n}\n",
				},
				"http": map[string]interface{}{
					"method": "GET",
					"url":    "/v1/todos123",
				},
			},
		},
	}

	// Tests
	tests := []struct {
		name    string
		k       K8s
		configs []*unstructured.Unstructured
		want    common.ConfigType
	}{
		{
			name: "add a compiledgraphqlsource",
			k: K8s{
				configuration: common.ConfigType{},
			},
			configs: []*unstructured.Unstructured{
				config1,
			},
			want: common.ConfigType{
				"source.core---space-cloud---io--v1alpha1--compiledgraphqlsources": []*unstructured.Unstructured{
					config1,
				},
			},
		},
		{
			name: "add 2 compiledgraphqlsources and 1 graphqlsource",
			k: K8s{
				configuration: common.ConfigType{},
			},
			configs: []*unstructured.Unstructured{
				config1, config2, config3,
			},
			want: common.ConfigType{
				"source.core---space-cloud---io--v1alpha1--compiledgraphqlsources": []*unstructured.Unstructured{
					config1, config2,
				},
				"source.core---space-cloud---io--v1alpha1--graphqlsources": []*unstructured.Unstructured{
					config3,
				},
			},
		},
		{
			name: "add 2 compiledgraphqlsources, 1 graphqlsource and update 1 compiledgraphqlsources",
			k: K8s{
				configuration: common.ConfigType{},
			},
			configs: []*unstructured.Unstructured{
				config1, config2, config3, config4,
			},
			want: common.ConfigType{
				"source.core---space-cloud---io--v1alpha1--compiledgraphqlsources": []*unstructured.Unstructured{
					config4, config2,
				},
				"source.core---space-cloud---io--v1alpha1--graphqlsources": []*unstructured.Unstructured{
					config3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, config := range tt.configs {
				tt.k.addOrUpdateConfig(config)
			}

			if !reflect.DeepEqual(tt.k.configuration, tt.want) {
				t.Errorf("configuration mismatched for test %s\n", tt.name)
				fmt.Println("want:")
				fmt.Println(tt.want)
				fmt.Println("got:")
				fmt.Println(tt.k.configuration)
				return
			}
		})
	}
}
