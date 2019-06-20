package static

import (
	"strings"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

// Module is responsible for static
type Module struct {
	sync.RWMutex
	Enabled bool
	routes  []*config.StaticRoute
}

// Init returns a new instance of the static module wit default values
func Init() *Module {
	return &Module{Enabled: false}
}

// SetConfig set the config required by the static module
func (m *Module) SetConfig(s *config.Static) error {
	m.Lock()
	defer m.Unlock()

	if s == nil || !s.Enabled {
		m.Enabled = false
		return nil
	}
	
	m.routes = s.Routes
	m.Enabled = true

	return nil
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.Enabled
}

// SelectRoute select the rules for a given request
func (m *Module) SelectRoute(host, url string) (*config.StaticRoute, bool) {
	m.RLock()
	defer m.RUnlock()

	for _, route := range m.routes {
		if strings.HasPrefix(url, route.URLPrefix) {
			if route.Host != "" && route.Host != host {
				continue
			}
			return route, true
		}
	}

	return nil, false
}
