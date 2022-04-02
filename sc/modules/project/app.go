package project

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/configman"
)

func init() {
	caddy.RegisterModule(App{})
	_ = configman.RegisterConfigController("project", getTypeDefinitions())
}

// App manages all the database modules
type App struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "project",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the file loader module.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	return nil
}

// Start begins the app's operation
func (a *App) Start() error {
	return nil
}

// Stop shuts down the app's operation
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner  = (*App)(nil)
	_ caddy.App          = (*App)(nil)
	_ configman.HookImpl = (*App)(nil)
)
