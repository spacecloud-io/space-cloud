package common

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func prepareRestApp(configuration map[string][]*unstructured.Unstructured) []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"compiledGraphqlQueries": configuration["CompiledGraphqlSource"],
	})
	return data
}
