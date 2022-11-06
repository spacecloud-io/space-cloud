package rest

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func init() {
	caddy.RegisterModule(App{})
	apis.RegisterApp("rest", 200)
}

// App describes the state of the auth app
type App struct {
	CompiledGraphqlQueries []*v1alpha1.CompiledGraphqlSource `json:"compiledGraphqlQueries"`

	// For internal use
	logger *zap.Logger

	// APIs
	apis apis.APIs

	// Graphql dependencies
	graphqlApp *graphql.App
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "rest",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the auth module.
func (a *App) Provision(ctx caddy.Context) error {
	// Get the logger
	a.logger = ctx.Logger(a)

	// Get all the dependencies
	gAppTemp, _ := ctx.App("graphql")
	a.graphqlApp = gAppTemp.(*graphql.App)

	// Prepare all the rest endpoints
	if err := a.prepareCompilesGraphqlEndpoints(); err != nil {
		a.logger.Error("Unable to compile provided graphql queries", zap.Error(err))
		return err
	}

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
