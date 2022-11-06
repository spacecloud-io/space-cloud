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

// App describes the state of the auth app
type App struct {
	HSASecrets  []*v1alpha1.JwtHSASecret
	OPAPolicies []*v1alpha1.OPAPolicy

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

// Provision sets up the auth module.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	return nil
}

// Start begins the auth app operations
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

// Stop ends the auth app operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
