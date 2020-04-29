package userman

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Module is responsible for user management
type Module struct {
	sync.RWMutex
	methods map[string]*config.AuthStub
	crud    model.CrudUserInterface
	auth    model.AuthUserInterface
}

// Init creates a new instance of the user management object
func Init(crud model.CrudUserInterface, auth model.AuthUserInterface) *Module {
	return &Module{crud: crud, auth: auth}
}

// SetConfig sets the config required by the user management module
func (m *Module) SetConfig(auth config.Auth) {
	m.Lock()
	defer m.Unlock()

	m.methods = make(map[string]*config.AuthStub, len(auth))

	for k, v := range auth {
		m.methods[k] = v
	}
}

// IsActive shows if a given method is active
func (m *Module) IsActive(method string) bool {
	m.RLock()
	defer m.RUnlock()

	s, p := m.methods[method]
	return p && s.Enabled
}

// IsEnabled shows if the user management module is enabled
func (m *Module) IsEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return len(m.methods) > 0
}
