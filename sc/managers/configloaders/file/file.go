package file

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/configloaders/common"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func init() {
	caddy.RegisterModule(Loader{})
}

// Loader loads the sc configuration from a YAML/JSON file
type Loader struct {
	Module string `json:"module"`
	Path   string `json:"path"`

	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Loader) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.config_loaders.file",
		New: func() caddy.Module { return new(Loader) },
	}
}

// Provision sets up the file loader module.
func (l *Loader) Provision(ctx caddy.Context) error {
	l.logger = ctx.Logger(l)
	return nil
}

// LoadConfig returns the final caddy config from the store.
func (l *Loader) LoadConfig(ctx caddy.Context) ([]byte, error) {
	// Load SC config file from file system
	fileConfig := new(model.SCConfig)
	if err := utils.LoadFile(l.Path, fileConfig); err != nil {
		l.logger.Error("Unable to load SpaceCloud config file", zap.Error(err))
		return nil, err
	}

	// Load the new caddy config
	config, err := common.PrepareConfig(fileConfig)
	if err != nil {
		l.logger.Error("Unable to prepare caddy config", zap.Error(err))
		return nil, err
	}

	return json.MarshalIndent(config, "", "  ")
}

// Interface guards
var (
	_ caddy.Provisioner  = (*Loader)(nil)
	_ caddy.ConfigLoader = (*Loader)(nil)
)
