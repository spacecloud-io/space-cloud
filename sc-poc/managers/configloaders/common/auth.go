package common

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func prepareAuthApp(configuration map[string][]*unstructured.Unstructured) []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"hsaSecrets": configuration["HSASecret"],
	})
	return data
}
