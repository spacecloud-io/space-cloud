package functions

import (
	"context"
	"fmt"
	"sync"
	"text/template"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Module is responsible for functions
type Module struct {
	lock sync.RWMutex

	// Dependencies
	auth           model.AuthFunctionInterface
	manager        *syncman.Manager
	integrationMan integrationManagerInterface

	// Variable configuration
	project    string
	metricHook model.MetricFunctionHook
	config     *config.ServicesModule

	// Templates for body transformation
	templates map[string]*template.Template
}

// Init returns a new instance of the Functions module
func Init(auth model.AuthFunctionInterface, manager *syncman.Manager, integrationMan integrationManagerInterface, hook model.MetricFunctionHook) *Module {
	return &Module{auth: auth, manager: manager, integrationMan: integrationMan, metricHook: hook}
}

const (
	module            string = "remote-services"
	segmentGoTemplate string = "goTemplate"
)

// SetConfig sets the configuration of the functions module
func (m *Module) SetConfig(project string, c *config.ServicesModule) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if c == nil {
		helpers.Logger.LogWarn(helpers.GetRequestID(context.TODO()), "Empty config provided for functions module", map[string]interface{}{"project": project})
		return nil
	}

	m.project = project
	m.config = c

	m.config.InternalServices = config.Services{}

	// Set the go templates
	m.templates = map[string]*template.Template{}
	for serviceID, service := range m.config.Services {
		for endpointID, endpoint := range service.Endpoints {
			// Set the default endpoint kind
			if endpoint.Kind == "" {
				endpoint.Kind = config.EndpointKindInternal
			}

			// Set default templating engine
			if endpoint.Tmpl == "" {
				endpoint.Tmpl = config.EndpointTemplatingEngineGo
			}

			// Set default output format
			if endpoint.OpFormat == "" {
				endpoint.OpFormat = "yaml"
			}

			switch endpoint.Tmpl {
			case config.EndpointTemplatingEngineGo:
				if endpoint.ReqTmpl != "" {
					if err := m.createGoTemplate("request", serviceID, endpointID, endpoint.ReqTmpl); err != nil {
						return err
					}
				}
				if endpoint.ResTmpl != "" {
					if err := m.createGoTemplate("response", serviceID, endpointID, endpoint.ResTmpl); err != nil {
						return err
					}
				}
				if endpoint.GraphTmpl != "" {
					if err := m.createGoTemplate("graph", serviceID, endpointID, endpoint.GraphTmpl); err != nil {
						return err
					}
				}
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid templating engine (%s) provided", endpoint.Tmpl), nil, nil)
			}
		}
	}
	return nil
}
