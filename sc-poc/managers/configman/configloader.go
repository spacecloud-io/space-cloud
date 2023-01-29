package configman

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bep/debounce"
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/configman/adapter"
	"github.com/spacecloud-io/space-cloud/managers/configman/adapter/file"
	"github.com/spacecloud-io/space-cloud/managers/configman/adapter/k8s"
	"github.com/spf13/viper"
)

// ConfigLoader defines a space cloud config loader.
type ConfigLoader struct {
	adapter          adapter.Adapter
	debounceInterval time.Duration
}

// InitializeConfigLoader initializes the config loader with the given adapter.
func InitializeConfigLoader() (*ConfigLoader, error) {
	configAdapter := viper.GetString("config.adapter")
	path := viper.GetString("config.path")
	debounceInterval, err := time.ParseDuration(viper.GetString("config.debounce-interval"))
	if err != nil {
		return nil, err
	}

	configloader := &ConfigLoader{
		debounceInterval: debounceInterval,
	}
	switch configAdapter {
	case "file":
		configloader.adapter = file.MakeFileAdapter(path)
	case "k8s":
		configloader.adapter, err = k8s.MakeK8sAdapter()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid adapter specified")
	}
	return configloader, nil
}

// GetCaddyConfig reads the config from the adapter and
// converts from bytes to caddy Config.
func (configloader *ConfigLoader) GetCaddyConfig() (*caddy.Config, error) {
	raw, err := configloader.adapter.GetRawConfig()
	if err != nil {
		return nil, err
	}

	c := &caddy.Config{}
	err = json.Unmarshal(raw, c)
	return c, err
}

// WatchChanges continuously watches the config objects and reloads caddy
// if it detects any changes
func (configloader *ConfigLoader) WatchChanges(ctx context.Context) {
	debounced := debounce.New(configloader.debounceInterval)
	cfgChan, err := configloader.adapter.Run(ctx)
	if err != nil {
		fmt.Println("Error watching changes: ", err)
		return
	}
	for cfgJSON := range cfgChan {
		debounced(func() { loadConfig(cfgJSON) })
	}
}

func loadConfig(cfgJSON []byte) {
	_ = caddy.Load(cfgJSON, false)
}
