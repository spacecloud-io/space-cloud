package common

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func prepareGraphQLApp(configuration map[string][]*unstructured.Unstructured) []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"graphqlSources":  configuration["GraphqlSource"],
		"compiledQueries": configuration["CompiledGraphqlSource"],
	})
	return data
}
