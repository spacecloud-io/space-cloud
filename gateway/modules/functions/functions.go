package functions

import (
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
			// Add go templates only
			if endpoint.Kind == config.EndpointKindTransform && endpoint.Tmpl != "" {
				key := getGoTemplateKey(serviceID, endpointID)

				// Create a new template object
				t := template.New(key)
				t = t.Funcs(m.createGoFuncMaps())
				tmpl, err := t.Parse(endpoint.Tmpl)
				if err != nil {
					return utils.LogError("Invalid golang template provided", module, segmentSetConfig, err)
				}

				// Save it for later use
				m.templates[getGoTemplateKey(serviceID, endpointID)] = tmpl
			}
		}
	}
	return nil
}
