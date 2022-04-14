package common

import (
	"encoding/json"

	"github.com/spf13/viper"
)

func prepareStoreApp() json.RawMessage {

	storeType := viper.GetString("store-type")
	configPath := viper.GetString("config-path")

	config := map[string]interface{}{
		"type": storeType,
		"path": configPath,
	}

	raw, _ := json.Marshal(config)
	return raw
}
