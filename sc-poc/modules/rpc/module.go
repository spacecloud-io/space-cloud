package rpc

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/provider"
	"github.com/spacecloud-io/space-cloud/managers/source"
)

func init() {
	caddy.RegisterModule(Module{})
	provider.Register("rpc", 0)
}

// Module describes the state of the auth app
type Module struct {
	Workspace string `json:"workspace"`

	// For internal use
	logger *zap.Logger

	// APIs
	apis apis.APIs
}

// CaddyModule returns the Caddy module information.
func (Module) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "provider.rpc",
		New: func() caddy.Module { return new(Module) },
	}
}

// Provision sets up the auth module.
func (m *Module) Provision(ctx caddy.Context) error {
	// Get the logger
	m.logger = ctx.Logger(m)

	// Get all the dependencies
	sourceManT, _ := ctx.App("source")
	sourceMan := sourceManT.(*source.App)

	for _, s := range sourceMan.GetSources(m.Workspace, "rpc") {
		rpcSource, ok := s.(Source)
		if ok {
			m.prepareAPIs(rpcSource)
		}
	}

	// // Prepare all the rest endpoints
	// if err := a.prepareCompilesGraphqlEndpoints(); err != nil {
	// 	a.logger.Error("Unable to compile provided graphql queries", zap.Error(err))
	// 	return err
	// }

	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*Module)(nil)
	_ apis.App          = (*Module)(nil)
)
