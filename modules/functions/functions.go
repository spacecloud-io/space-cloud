package functions

import (
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// Module is responsible for functions
type Module struct {
	lock sync.RWMutex

	// Dependencies
	auth    *auth.Module
	manager *syncman.Manager

	// Variable configuration
	project string
	config  *config.ServicesModule
}

// Init returns a new instance of the Functions module
func Init(auth *auth.Module, manager *syncman.Manager) *Module {
	return &Module{auth: auth, manager: manager}
}

// SetConfig sets the configuration of the functions module
func (m *Module) SetConfig(project string, c *config.ServicesModule) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if c == nil {
		return
	}

	m.project = project
	m.config = c

	m.config.InternalServices = config.Services{}
}
