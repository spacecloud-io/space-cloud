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
	interval := viper.GetInt64("loading-interval")
	loadingInterval := caddy.Duration(time.Second) * caddy.Duration(interval)
	if isInitialLoad {
		loadingInterval = 0
	}

	// Selecting store-type
	var loader json.RawMessage
	switch v := viper.GetString("store-type"); v {
	case "file":
		loader = prepareFileLoaderConfig()
	case "kube":
		loader = prepareKubeLoaderConfig()
	default:
		return nil, fmt.Errorf("store-type (%s) is not suppoerted", v)
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
					Level: viper.GetString("log-level"),
				},
			},
		},
	}, nil
}

func prepareFileLoaderConfig() json.RawMessage {
	config := map[string]interface{}{
		"module": "file",
		"path":   viper.GetString("config-path"),
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
