package auth

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

// App describes the state of the auth app
type App struct {
	HSASecrets             []*v1alpha1.JwtHSASecret          `json:"hsaSecrets"`
	OPAPolicies            []*v1alpha1.OPAPolicy             `json:"opaPolicies"`
	CompiledGraphqlSources []*v1alpha1.CompiledGraphqlSource `json:"compiledGraphqlSources"`

	// For internal use
	logger       *zap.Logger
	secrets      []*types.AuthSecret
	regoPolicies map[string]rego.PreparedEvalQuery

	// Dependant apps
	graphqlApp *graphql.App
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "auth",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the auth module.W
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)

	// Load the dependent apps
	graphqlTemp, err := ctx.App("graphql")
	if err != nil {
		a.logger.Error("Unable to load the graphql application", zap.Error(err))
		return err
	}
	a.graphqlApp = graphqlTemp.(*graphql.App)

	// Compile the rego policies
	if err := a.compileRegoPolicies(); err != nil {
		a.logger.Error("Unable to compile rego policies", zap.Error(err))
		return err
	}

	return nil
}

// Start begins the auth app operations
func (a *App) Start() error {
	// TODO: add support of rsa secrets
	// TODO: add support for jwk urls

	secrets := make([]*types.AuthSecret, 0, len(a.HSASecrets))
	for _, s := range a.HSASecrets {
		secrets = append(secrets, &types.AuthSecret{
			AuthSecret: s.Spec.AuthSecret,
			Alg:        types.HS256,
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
