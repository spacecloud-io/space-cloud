package configloaders

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/managers/configloaders/common"
	"github.com/spacecloud-io/space-cloud/utils"
)

func init() {
	caddy.RegisterModule(FileLoader{})
}

// FileLoader loads the sc configuration from a YAML/JSON file
type FileLoader struct {
	Module string `json:"module"`
	Path   string `json:"path"`

	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (FileLoader) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.config_loaders.file",
		New: func() caddy.Module { return new(FileLoader) },
	}
}

// Provision sets up the file loader module.
func (l *FileLoader) Provision(ctx caddy.Context) error {
	l.logger = ctx.Logger(l)
	return nil
}

// LoadConfig returns the final caddy config from the store.
func (l *FileLoader) LoadConfig(ctx caddy.Context) ([]byte, error) {
	// Load SC config file from file system
	scConfig := new(config.Config)
	if err := utils.LoadFile(l.Path, scConfig); err != nil {
		l.logger.Error("Unable to load SpaceCloud config file", zap.Error(err))
		return nil, err
	}

	// Load the new caddy config
	return json.Marshal(common.PrepareConfig(scConfig))
}

// Interface guards
var (
	_ caddy.Provisioner  = (*FileLoader)(nil)
	_ caddy.ConfigLoader = (*FileLoader)(nil)
)
