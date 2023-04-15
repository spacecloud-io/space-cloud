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
	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spf13/viper"
)

// ConfigLoader defines a space cloud config loader.
type ConfigLoader struct {
	adapter          adapter.Adapter
	debounceInterval time.Duration
}

// configLoader is a globally initialized config loader.
var configLoader ConfigLoader = ConfigLoader{}

// InitializeConfigLoader initializes the config loader with the given adapter.
func InitializeConfigLoader() error {
	configAdapter := viper.GetString("config.adapter")
	path := viper.GetString("config.path")
	debounceInterval, err := time.ParseDuration(viper.GetString("config.debounce-interval"))
	if err != nil {
		return err
	}

	configLoader.debounceInterval = debounceInterval
	switch configAdapter {
	case "file":
		configLoader.adapter = file.MakeFileAdapter(path)
	case "k8s":
		configLoader.adapter, err = k8s.MakeK8sAdapter()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid adapter specified")
	}
	return nil
}

// GetCaddyConfig reads the config from the adapter and
// converts from bytes to caddy Config.
func GetCaddyConfig() (*caddy.Config, error) {
	cfg, err := configLoader.adapter.GetRawConfig()
	if err != nil {
		return nil, err
	}

	return common.PrepareConfig(cfg)
}

// WatchChanges continuously watches the config objects and reloads caddy
// if it detects any changes
func WatchChanges(ctx context.Context) {
	debounced := debounce.New(configLoader.debounceInterval)
	cfgChan, err := configLoader.adapter.Run(ctx)
	if err != nil {
		fmt.Println("Error watching changes: ", err)
		return
	}
	for cfg := range cfgChan {
		debounced(func() { loadConfig(cfg) })
	}
}

func loadConfig(cfg common.ConfigType) {
	caddyCfg, err := common.PrepareConfig(cfg)
	if err != nil {
		fmt.Println("Error watching changes: ", err)
		return
	}

	raw, err := json.MarshalIndent(caddyCfg, "", "  ")
	if err != nil {
		fmt.Println("Error watching changes: ", err)
		return
	}
	_ = caddy.Load(raw, false)
}
