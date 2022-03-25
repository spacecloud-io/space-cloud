package common

import "encoding/json"

func prepareGraphQLApp() []byte {
	data, _ := json.Marshal(map[string]interface{}{})
	return data
}
