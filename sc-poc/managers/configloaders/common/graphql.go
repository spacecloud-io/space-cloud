package common

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func prepareGraphQLApp(configuration map[string][]*unstructured.Unstructured) []byte {
	fmt.Println("++++", configuration)
	data, _ := json.Marshal(map[string]interface{}{
		"graphqlSources": configuration["GraphqlSource"],
	})
	return data
}
