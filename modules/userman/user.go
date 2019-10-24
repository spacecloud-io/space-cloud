package userman

import (
	"sync"

	"github.com/spaceuptech/space-cloud/config"

	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
)

// Module is responsible for user management
type Module struct {
	sync.RWMutex
	methods map[string]*config.AuthStub
	crud    *crud.Module
	auth    *auth.Module
}

// Init creates a new instance of the user management object
func Init(crud *crud.Module, auth *auth.Module) *Module {
	return &Module{crud: crud, auth: auth}
}

// SetConfig set the config required by the user management module
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
