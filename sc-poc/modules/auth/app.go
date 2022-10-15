package auth

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func init() {
	caddy.RegisterModule(App{})
	caddy.RegisterModule(AuthHandler{})
}

type App struct {
	HSASecrets []*v1alpha1.HSASecret

	// For internal use
	logger  *zap.Logger
	secrets []*AuthSecret
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "auth",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the graphql module.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	return nil
}

// Start begins the graphql app operations
func (a *App) Start() error {
	// TODO: add support of rsa secrets
	// TODO: add support for jwk urls

	secrets := make([]*AuthSecret, 0, len(a.HSASecrets))
	for _, s := range a.HSASecrets {
		secrets = append(secrets, &AuthSecret{
			AuthSecret: s.Spec.AuthSecret,
			Alg:        HS256,
			Value:      s.Spec.Value,
		})
	}

	a.secrets = secrets
	return nil
}

// Stop ends the graphql app operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
