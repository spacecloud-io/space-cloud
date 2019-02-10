package userman

import (
	"github.com/spaceuptech/space-cloud/auth"
	"github.com/spaceuptech/space-cloud/crud"
)

// Module is responsible for user management
type Module struct {
	methods map[string]struct{}
	crud    *crud.Module
	auth    *auth.Module
}

// Init creates a new instance of the user management object
func Init(crud *crud.Module, auth *auth.Module) *Module {
	return &Module{crud: crud, auth: auth}
}

func (m *Module) isActive(method string) bool {
	_, p := m.methods[method]
	return p
}
