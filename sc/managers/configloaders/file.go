package configloaders

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

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
	l.logger.Info("File module called")

	appRaw, _ := json.Marshal(map[string]string{})
	c := utils.LoadAdminConfig(false)
	c.AppsRaw = caddy.ModuleMap{"database": appRaw}
	return json.Marshal(c)
}

// Interface guards
var (
	_ caddy.Provisioner  = (*FileLoader)(nil)
	_ caddy.ConfigLoader = (*FileLoader)(nil)
)
