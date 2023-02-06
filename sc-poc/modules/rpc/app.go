package rpc

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/source"
)

func init() {
	caddy.RegisterModule(App{})
	apis.RegisterApp("rpc", 200)
}

// App describes the state of the auth app
type App struct {
	// For internal use
	logger *zap.Logger

	// APIs
	apis apis.APIs
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "rpc",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the auth module.
func (a *App) Provision(ctx caddy.Context) error {
	// Get the logger
	a.logger = ctx.Logger(a)

	// Get all the dependencies
	sourceAppT, _ := ctx.App("source")
	sourceApp := sourceAppT.(*source.App)

	for _, s := range sourceApp.GetSources("rpc") {
		// First resolve the source's dependencies
		if err := source.ResolveDependencies(ctx, "rpc", s); err != nil {
			a.logger.Error("Unable to resolve source's dependency", zap.String("source", s.GetName()), zap.Error(err))
			return err
		}

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
