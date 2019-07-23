package static

import "github.com/spaceuptech/space-cloud/config"

// AddProxyRoute adds a proxy route in the static module
func (m *Module) AddProxyRoute(id, host, prefix, proxy string) {
	m.Lock()
	defer m.Unlock()

	route := &config.StaticRoute{ID: id, Host: host, URLPrefix: prefix, Proxy: proxy}
	m.routes = append(m.routes, route)
}

// DeleteRoutesWithID removes all routes of particular id
func (m *Module) DeleteRoutesWithID(id string) {
	m.Lock()
	defer m.Unlock()

	// Filter out those routes whose ids don't match
	routes := []*config.StaticRoute{}
	for _, r := range m.routes {
		if r.ID != id {
			routes = append(routes, r)
		}
	}
	m.routes = routes
}
