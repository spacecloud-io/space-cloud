package auth

import (
	"context"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/source"
)

// App describes the state of the auth app
type App struct {
	// For internal use
	logger   *zap.Logger
	secrets  []SecretSource
	policies map[string]PolicySource
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "auth",
		New: func() caddy.Module { return new(App) },
	}
}

func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	sourceManT, err := ctx.App("source")
	if err != nil {
		a.logger.Error("Unable to load the source manager", zap.Error(err))
	}
	sourceMan := sourceManT.(*source.App)

	// Get all relevant sources
	sources := sourceMan.GetSources("auth")
	a.secrets = []SecretSource{}
	a.policies = make(map[string]PolicySource)
	for _, src := range sources {
		name := src.GetName()

		// First resolve the source's dependencies
		if err := source.ResolveDependencies(ctx, "auth", src); err != nil {
			a.logger.Error("Unable to resolve source's dependency", zap.String("source", src.GetName()), zap.Error(err))
			return err
		}

		// Get sources with secret
		s, ok := src.(SecretSource)
		if ok {
			a.secrets = append(a.secrets, s)
		}

		// Get sources with policy
		p, ok := src.(PolicySource)
		if ok {
			a.policies[name] = p
		}
	}

	return nil
}

func (a *App) EvaluatePolicy(ctx context.Context, name string, input interface{}) (bool, string, error) {
	if policy, ok := a.policies[name]; ok {
		return policy.Evaluate(ctx, input)
	}

	return false, "", fmt.Errorf("policy with name %s not found", name)
}

// Start begins the auth app operations
func (a *App) Start() error {
	return nil
}

// Stop ends the auth app operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
