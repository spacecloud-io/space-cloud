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
	caching        cachingInterface

	// Variable configuration
	project    string
	metricHook model.MetricFunctionHook
	config     config.Services

	clusterID string
	// Templates for body transformation
	templates map[string]*template.Template
}

// Init returns a new instance of the Functions module
func Init(clusterID string, auth model.AuthFunctionInterface, manager *syncman.Manager, integrationMan integrationManagerInterface, hook model.MetricFunctionHook) *Module {
	return &Module{clusterID: clusterID, auth: auth, manager: manager, integrationMan: integrationMan, metricHook: hook}
}

// SetConfig sets the configuration of the functions module
func (m *Module) SetConfig(project string, c config.Services) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.project = project
	m.config = c

	// Set the go templates
	m.templates = map[string]*template.Template{}
	for _, service := range m.config {
		for endpointID, endpoint := range service.Endpoints {
			// Set the default endpoint kind
			if endpoint.Kind == "" {
				endpoint.Kind = config.EndpointKindInternal
			}

			// Set default templating engine
			if endpoint.Tmpl == "" {
				endpoint.Tmpl = config.TemplatingEngineGo
			}

			// Set default output format
			if endpoint.OpFormat == "" {
				endpoint.OpFormat = "yaml"
			}

			if endpoint.Timeout == 0 {
				endpoint.Timeout = 60
			}

			switch endpoint.Tmpl {
			case config.TemplatingEngineGo:
				if endpoint.ReqTmpl != "" {
					if err := m.createGoTemplate("request", service.ID, endpointID, endpoint.ReqTmpl); err != nil {
						return err
					}
				}
				if endpoint.ResTmpl != "" {
					if err := m.createGoTemplate("response", service.ID, endpointID, endpoint.ResTmpl); err != nil {
						return err
					}
				}
				if endpoint.GraphTmpl != "" {
					if err := m.createGoTemplate("graph", service.ID, endpointID, endpoint.GraphTmpl); err != nil {
						return err
					}
				}
				if endpoint.Claims != "" {
					if err := m.createGoTemplate("claim", service.ID, endpointID, endpoint.Claims); err != nil {
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

// SetCachingModule sets caching module
func (m *Module) SetCachingModule(c cachingInterface) {
	m.caching = c
}
