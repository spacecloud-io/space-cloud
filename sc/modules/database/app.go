package database

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(DatabaseApp{})
}

// DatabaseApp manages all the database modules
type DatabaseApp struct {
	// The config this app needs
	ConnectorsRaw caddy.ModuleMap `json:"connectors,omitempty" caddy:"namespace=database.connectors inline_key=type"`

	// For internal usage
	logger     *zap.Logger
	connectors map[string]Connector
}

// CaddyModule returns the Caddy module information.
func (DatabaseApp) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "database",
		New: func() caddy.Module { return new(DatabaseApp) },
	}
}

// Provision sets up the file loader module.
func (l *DatabaseApp) Provision(ctx caddy.Context) error {
	l.logger = ctx.Logger(l)

	return nil
}

// Start begins the database app operations
func (l *DatabaseApp) Start() error {
	l.logger.Info("Starting app")
	return nil
}

// Stop ends the database app operations
func (l *DatabaseApp) Stop() error {
	l.logger.Info("Stopped app")
	return nil
}

// LoadConfig returns the final caddy config from the store.
func (l *DatabaseApp) LoadConfig(ctx caddy.Context) ([]byte, error) {
	c := utils.LoadAdminConfig(false)
	return json.Marshal(c)
}

// Interface guards
var (
	_ caddy.Provisioner  = (*DatabaseApp)(nil)
	_ caddy.ConfigLoader = (*DatabaseApp)(nil)
	_ caddy.App          = (*DatabaseApp)(nil)
)
