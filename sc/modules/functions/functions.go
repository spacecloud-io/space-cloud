package functions

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/caddyserver/caddy/v2"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/model"
)

func init() {
	caddy.RegisterModule(App{})
	configman.RegisterConfigController("remote_service")
	apis.RegisterApp("remote_service", 1)
}

// App manages all the admin actions
type App struct {
	// Remote services
	Services map[string]*config.Service `json:"services"`

	// Internal stuff
	templates map[string]*template.Template
	logger    *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "remote_service",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the app.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	// Set the go templates
	a.templates = map[string]*template.Template{}
	for name, service := range a.Services {
		projectID := strings.Split(name, "---")[0]
		serviceName := strings.Split(name, "---")[1]
		for endpointID, endpoint := range service.Endpoints {
			tmpl, err := a.getStringOutputFromPlugins(endpoint, config.PluginTmpl)
			if err != nil {
				a.logger.Error("Unable to load plugins app", zap.Error(err), zap.String("app", "remote_service"))
				return err
			}

			reqTmpl, err := a.getStringOutputFromPlugins(endpoint, config.PluginRequestTemplate)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote_service"))
				return err
			}

			resTmpl, err := a.getStringOutputFromPlugins(endpoint, config.PluginResponseTemplate)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote_service"))
				return err
			}

			graphTmpl, err := a.getStringOutputFromPlugins(endpoint, config.PluginGraphTemplate)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote_service"))
				return err
			}

			claims, err := a.getStringOutputFromPlugins(endpoint, config.PluginClaims)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote_service"))
				return err
			}

			switch tmpl {
			case string(config.TemplatingEngineGo):
				if reqTmpl != "" {
					if err := a.createGoTemplate("request", projectID, serviceName, endpointID, reqTmpl); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote_service"))
						return err
					}
				}
				if resTmpl != "" {
					if err := a.createGoTemplate("response", projectID, serviceName, endpointID, resTmpl); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote_service"))
						return err
					}
				}
				if graphTmpl != "" {
					if err := a.createGoTemplate("graph", projectID, serviceName, endpointID, graphTmpl); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote_service"))
						return err
					}
				}
				if claims != "" {
					if err := a.createGoTemplate("claim", projectID, serviceName, endpointID, claims); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote_service"))
						return err
					}
				}
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid templating engine (%s) provided", tmpl), nil, nil)
			}
		}
	}

	return nil
}

// Start begins the app's operation
func (a *App) Start() error {
	return nil
}

// Stop shuts down the app's operation
func (a *App) Stop() error {
	return nil
}

// GetConfigTypes returns all the operation types returned by this model.
func (a *App) GetConfigTypes() model.ConfigTypes {
	types := model.ConfigTypes{}
	for k, v := range getServiceConfigTypes() {
		types[k] = v
	}

	return types
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
	_ model.ConfigCtrl  = (*App)(nil)
	_ apis.App          = (*App)(nil)
)
