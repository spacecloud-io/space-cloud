package functions

import (
	"fmt"
	"sync"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Module is responsible for functions
type Module struct {
	lock sync.RWMutex

	// Dependencies
	auth    model.AuthFunctionInterface
	manager *syncman.Manager

	// Variable configuration
	project    string
	metricHook model.MetricFunctionHook
	config     *config.ServicesModule

	// Templates for body transformation
	templates map[string]*template.Template
}

// Init returns a new instance of the Functions module
func Init(auth model.AuthFunctionInterface, manager *syncman.Manager, hook model.MetricFunctionHook) *Module {
	return &Module{auth: auth, manager: manager, metricHook: hook}
}

const (
	module            string = "remote-services"
	segmentSetConfig  string = "set-config"
	segmentCall       string = "call"
	segmentGoTemplate string = "goTemplate"
)

// SetConfig sets the configuration of the functions module
func (m *Module) SetConfig(project string, c *config.ServicesModule) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if c == nil {
		utils.LogWarn("Empty config module provided", module, segmentSetConfig)
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
				return utils.LogError(fmt.Sprintf("Invalid templating engine (%s) provided", endpoint.Tmpl), module, segmentSetConfig, nil)
			}
		}
	}
	return nil
}
