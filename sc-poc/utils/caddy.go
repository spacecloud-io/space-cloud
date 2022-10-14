package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/spf13/viper"
)

// LoadAdminConfig creates a caddy config from the viper config provided. This only contains the admin
// and logging portion of the configuration. The config loaders (`manager/configloaders`) will be responsible to load the
// configuration of the applications.
func LoadAdminConfig(isInitialLoad bool) (*caddy.Config, error) {
	logLevel := viper.GetString("caddy.log-level")
	loadTime := viper.GetString("caddy.loading-interval")

	interval, err := time.ParseDuration(loadTime)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config loading interval (%s), error: %v", loadTime, err)
	}

	loadingDelay := caddy.Duration(interval)
	if isInitialLoad {
		loadingDelay = 0
	}

	persist := false
	return &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
			Config: &caddy.ConfigSettings{
				LoadDelay: loadingDelay,
				LoadRaw:   prepareFileLoaderConfig(viper.GetString("config.path")),
				Persist:   &persist,
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

// GetCaddyMatcherSet returns a caddy matcher set
func GetCaddyMatcherSet(path []string, methods []string) []caddy.ModuleMap {
	// We will always need to match based on the path
	set := map[string]json.RawMessage{
		"path": GetByteStringArray(path...),
	}

	// Match on method if provided
	if len(methods) > 0 {
		set["method"] = GetByteStringArray(methods...)
	}

	// Return the match set
	return []caddy.ModuleMap{set}
}

// GetByteStringArray returns an array of string in json form
func GetByteStringArray(val ...string) []byte {
	data, _ := json.Marshal(val)
	return data
}

// GetCaddyHandler returns a marshaled caddy handler config
func GetCaddyHandler(handlerName string, params map[string]interface{}) []json.RawMessage {
	handler := make(map[string]interface{})

	// Add the handler name / identifier
	handler["handler"] = fmt.Sprintf("sc_%s_handler", handlerName)

	// Add the params the handler needs
	for k, v := range params {
		handler[k] = v
	}

	data, _ := json.Marshal(handler)
	return []json.RawMessage{data}
}

func prepareFileLoaderConfig(path string) json.RawMessage {
	config := map[string]interface{}{
		"module": "file",
		"path":   path,
	}

	raw, _ := json.Marshal(config)
	return raw
}
