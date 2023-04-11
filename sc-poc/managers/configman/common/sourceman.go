package common

import (
	"encoding/json"
)

func prepareSourceManagerApp(configuration ConfigType) []byte {
	data, _ := json.Marshal(map[string]any{"config": configuration})
	return data
}
