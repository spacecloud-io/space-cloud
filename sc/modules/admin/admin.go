package admin

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils/auth"
)

func init() {
	caddy.RegisterModule(App{})
	caddy.RegisterModule(AuthHandler{})
	configman.RegisterOperationController("adminman")
	configman.RegisterConfigController("adminman")
}

// App manages all the admin actions
type App struct {
	// Admin credentials
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Secret string `json:"secret"`

	// Cluster parameters
	IsDev bool `json:"isDev"`

	// Project config
	Projects map[string]*Project `json:"projects"` // Key is project id

	// Internal stuff
	logger *zap.Logger
	auth   *auth.Module
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "adminman",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the app.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	a.auth = auth.New([]*config.Secret{
		{IsPrimary: true, Alg: config.HS256, KID: "sc-admin", Secret: a.Secret},
	})
	return nil
}

// Start begins the app's operation
func (a *App) Start() error {
	// Initialise the auth module for each project
	for _, p := range a.Projects {
		p.auth = auth.New(p.Secrets)
	}
	return nil
}

// Stop shuts down the app's operation
func (a *App) Stop() error {
	return nil
}

// GetOperationTypes returns all the operation types returned by this model.
func (a *App) GetOperationTypes() model.OperationTypes {
	types := model.OperationTypes{}
	for k, v := range a.getGlobalTypes() {
		types[k] = v
	}
	for k, v := range a.getProjectOperationTypes() {
		types[k] = v
	}

	return types
}

// GetConfigTypes returns all the operation types returned by this model.
func (a *App) GetConfigTypes() model.ConfigTypes {
	types := model.ConfigTypes{}
	for k, v := range getProjectConfigTypes() {
		types[k] = v
	}

	return types
}

// Auth returns the admin auth module
func (a *App) Auth() *auth.Module {
	return a.auth
}

// Interface guards
var (
	_ caddy.Provisioner   = (*App)(nil)
	_ caddy.App           = (*App)(nil)
	_ model.OperationCtrl = (*App)(nil)
	_ model.ConfigCtrl    = (*App)(nil)
)
