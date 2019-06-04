package static

import (
	"strings"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

// Module is responsible for Static
type Module struct {
	sync.RWMutex
	Enabled bool
	routes  []*config.StaticRoute
	Gzip    bool
}

// Init returns a new instance of the Static module wit default values
func Init() *Module {
	return &Module{Enabled: false, Gzip: false}
}

// SetConfig set the config required by the Static module
func (m *Module) SetConfig(s *config.Static) error {
	m.Lock()
	defer m.Unlock()

	if s == nil || !s.Enabled {
		m.Enabled = false
		return nil
	}

	m.Gzip = s.Gzip
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
