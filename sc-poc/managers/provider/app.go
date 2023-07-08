package provider

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/source"
)

func init() {
	caddy.RegisterModule(App{})
	apis.RegisterApp("provider", 10)
}

// App describes the provider manager app
type App struct {
	// Internal stuff
	logger *zap.Logger

	workspaces []workspaceset
	apis       apis.APIs
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "provider",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the provider manager.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)

	// Get the source manager
	sourceManT, err := ctx.App("source")
	if err != nil {
		a.logger.Error("Unable to load the source manager", zap.Error(err))
	}
	sourceMan := sourceManT.(*source.App)

	// Get list of registered workspaces
	workspaces := sourceMan.GetWorkspaces()

	// Initialise a Provider set for each workspace
	for _, ws := range workspaces {
		providers, err := a.generateProviders(ctx, ws)
		if err != nil {
			return err
		}
		a.workspaces = append(a.workspaces, workspaceset{name: ws, providers: providers})
	}

	// Add providerset for the main workspace
	mainProviders, err := a.generateProviders(ctx, "main")
	if err != nil {
		return err
	}
	a.workspaces = append(a.workspaces, workspaceset{"main", mainProviders})

	// Prepare all api routes
	a.prepareAPIRoutes()

	return nil
}

// Start begins the provider manager operations
func (a *App) Start() error {
	return nil
}

// Stop ends the provider manager operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
	_ apis.App          = (*App)(nil)
)
