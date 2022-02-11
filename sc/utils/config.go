package utils

import (
	"encoding/json"
	"time"

	"github.com/caddyserver/caddy/v2"
)

// LoadAdminConfig creates a caddy config from the viper config provided. This only contains the admin
// and logging portion of the configuration. The config loaders will be responsible to load the
// configuration of the applications.
func LoadAdminConfig(isInitialLoad bool) *caddy.Config {
	// TODO: Read base config from viper over here.
	loadingInterval := caddy.Duration(time.Second) * 5
	if isInitialLoad {
		loadingInterval = 0
	}

	return &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
			Config: &caddy.ConfigSettings{
				LoadInterval: loadingInterval,

				// TODO: Choose the right loader based on the flags
				LoadRaw: prepareFileLoaderConfig(),
			},
		},
		// TODO: Configure logging as well
	}
}

func prepareFileLoaderConfig() json.RawMessage {
	config := map[string]interface{}{
		"module": "file",
		"path":   "./config.yaml",
	}

	raw, _ := json.Marshal(config)
	return raw
}
