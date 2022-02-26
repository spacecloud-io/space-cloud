package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/spf13/viper"
)

// LoadAdminConfig creates a caddy config from the viper config provided. This only contains the admin
// and logging portion of the configuration. The config loaders will be responsible to load the
// configuration of the applications.
func LoadAdminConfig(isInitialLoad bool) (*caddy.Config, error) {
	logLevel := viper.GetString("log-level")
	loadTime := viper.GetString("loading-interval")

	storeType := viper.GetString("store-type")
	configPath := viper.GetString("config-path")

	interval, err := time.ParseDuration(loadTime)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config loading interval (%s), error: %v", loadTime, err)
	}

	loadingInterval := caddy.Duration(interval)
	if isInitialLoad {
		loadingInterval = 0
	}

	// Selecting store-type
	var loader json.RawMessage
	switch storeType {
	case "file":
		loader = prepareFileLoaderConfig(configPath)
	case "kube":
		loader = prepareKubeLoaderConfig()
	default:
		return nil, fmt.Errorf("store-type (%s) is not suppoerted", storeType)
	}

	return &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
			Config: &caddy.ConfigSettings{
				LoadInterval: loadingInterval,
				LoadRaw:      loader,
			},
		},
		Logging: &caddy.Logging{
			Logs: map[string]*caddy.CustomLog{
				"default": {
					Level: logLevel,
				},
			},
		},
	}, nil
}

func prepareFileLoaderConfig(path string) json.RawMessage {
	config := map[string]interface{}{
		"module": "file",
		"path":   path,
	}

	raw, _ := json.Marshal(config)
	return raw
}

func prepareKubeLoaderConfig() json.RawMessage {
	config := map[string]interface{}{
		"module": "kube",
	}

	raw, _ := json.Marshal(config)
	return raw
}
