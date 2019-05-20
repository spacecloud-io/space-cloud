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
	config  *config.Auth
	crud    *crud.Module
	auth    *auth.Module
	project string
}

// Init creates a new instance of the user management object
func Init(crud *crud.Module, auth *auth.Module) *Module {
	return &Module{crud: crud, auth: auth}
}

// SetConfig set the config required by the user management module
func (m *Module) SetConfig(project string, config *config.Auth) {
	m.Lock()
	defer m.Unlock()

	m.project = project
	m.config = config
}

func (m *Module) getOAuth(method string) (*config.OAuth, bool) {
	m.RLock()
	defer m.RUnlock()

	if m.config == nil {
		return nil, false
	}

	switch method {
	case "google":
		if m.config.Google != nil && m.config.Google.Enabled {
			return m.config.Google, true
		}
	case "twitter":
		if m.config.Twitter != nil && m.config.Twitter.Enabled {
			return m.config.Twitter, true
		}
	case "fb":
		if m.config.Facebook != nil && m.config.Facebook.Enabled {
			return m.config.Facebook, true
		}
	case "github":
		if m.config.Github != nil && m.config.Github.Enabled {
			return m.config.Github, true
		}
	}

	return nil, false
}

// IsActive checks if given sign on method is enabled
func (m *Module) IsActive(method string) bool {
	m.RLock()
	defer m.RUnlock()

	if m.config == nil {
		return false
	}

	switch method {
	case "email":
		if m.config.Email != nil {
			return m.config.Email.Enabled
		}
	case "google":
		if m.config.Google != nil {
			return m.config.Google.Enabled
		}
	case "twitter":
		if m.config.Twitter != nil {
			return m.config.Twitter.Enabled
		}
	case "fb":
		if m.config.Facebook != nil {
			return m.config.Facebook.Enabled
		}
	case "github":
		if m.config.Github != nil {
			return m.config.Github.Enabled
		}
	}

	return false
}

func (m *Module) isEnabled() bool {
	methods := []string{"email", "google", "twitter", "fb", "github"}
	for _, method := range methods {
		if p := m.IsActive(method); p {
			return p
		}
	}
	return false
}

func (m *Module) getDefaultRole() string {
	m.RLock()
	defer m.RUnlock()

	if m.config == nil {
		return "default"
	}
	if m.config.DefaultRole != "" {
		return m.config.DefaultRole
	}
	return "default"
}
