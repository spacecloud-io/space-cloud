package rpc

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/provider"
	"github.com/spacecloud-io/space-cloud/managers/source"
)

func init() {
	caddy.RegisterModule(App{})
	provider.Register("rpc", 0)
}

// App describes the state of the auth app
type App struct {
	Workspace string `json:"workspace"`

	// For internal use
	logger *zap.Logger

	// APIs
	apis apis.APIs
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "provider.rpc",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the auth module.
func (a *App) Provision(ctx caddy.Context) error {
	// Get the logger
	a.logger = ctx.Logger(a)

	// Get all the dependencies
	sourceManT, _ := ctx.App("source")
	sourceMan := sourceManT.(*source.App)

	for _, s := range sourceMan.GetSources(a.Workspace, "rpc") {
		rpcSource, ok := s.(Source)
		if ok {
			a.prepareAPIs(rpcSource)
		}
	}

	// // Prepare all the rest endpoints
	// if err := a.prepareCompilesGraphqlEndpoints(); err != nil {
	// 	a.logger.Error("Unable to compile provided graphql queries", zap.Error(err))
	// 	return err
	// }

	return nil
}

// Start begins the rest app operations
func (a *App) Start() error {
	return nil
}

// Stop ends the rest app operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
	_ apis.App          = (*App)(nil)
)
