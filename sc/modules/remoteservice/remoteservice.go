package remoteservice

import (
	"context"
	"fmt"
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
	configman.RegisterConfigController("remote-service")
	apis.RegisterApp("remote-service", 1)
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
		ID:  "remote-service",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the app.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	// Set the go templates
	a.templates = map[string]*template.Template{}
	for name, service := range a.Services {
		projectID, serviceName := SplitPojectRemoteServiceName(name)
		for endpointID, endpoint := range service.Endpoints {

			reqTmpl, _, reqTmplEngine, err := a.getRequestTemplatePlugin(endpoint)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote-service"))
				return err
			}

			resTmpl, _, resTmplEngine, err := a.getResponseTemplatePlugin(endpoint)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote-service"))
				return err
			}

			graphTmpl, _, graphTmplEngine, err := a.getGraphTemplatePlugin(endpoint)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote-service"))
				return err
			}

			claims, _, claimsTmplEngine, err := a.getClaimsPlugin(endpoint)
			if err != nil {
				a.logger.Error("Unable to load plugins", zap.Error(err), zap.String("app", "remote-service"))
				return err
			}

			switch reqTmplEngine {
			case string(config.TemplatingEngineGo):
				if reqTmpl != "" {
					if err := a.createGoTemplate("request", projectID, serviceName, endpointID, reqTmpl); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote-service"))
						return err
					}
				}
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid templating engine (%s) provided", reqTmplEngine), nil, nil)
			}

			switch resTmplEngine {
			case string(config.TemplatingEngineGo):
				if resTmpl != "" {
					if err := a.createGoTemplate("response", projectID, serviceName, endpointID, resTmpl); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote-service"))
						return err
					}
				}
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid templating engine (%s) provided", resTmplEngine), nil, nil)
			}

			switch graphTmplEngine {
			case string(config.TemplatingEngineGo):
				if graphTmpl != "" {
					if err := a.createGoTemplate("graph", projectID, serviceName, endpointID, graphTmpl); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote-service"))
						return err
					}
				}
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid templating engine (%s) provided", graphTmplEngine), nil, nil)
			}

			switch claimsTmplEngine {
			case string(config.TemplatingEngineGo):
				if claims != "" {
					if err := a.createGoTemplate("claim", projectID, serviceName, endpointID, claims); err != nil {
						a.logger.Error("Unable to load create go template", zap.Error(err), zap.String("app", "remote-service"))
						return err
					}
				}
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid templating engine (%s) provided", claimsTmplEngine), nil, nil)
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
	return getServiceConfigTypes()
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
	_ model.ConfigCtrl  = (*App)(nil)
	_ apis.App          = (*App)(nil)
)
