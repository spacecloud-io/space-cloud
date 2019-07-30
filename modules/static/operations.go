package static

import "github.com/spaceuptech/space-cloud/config"

// CompareAndAddInternalRoutes adds a internal route in the static module
func (m *Module) CompareAndAddInternalRoutes(routes []*config.StaticRoute) []*config.StaticRoute {
	m.Lock()
	defer m.Unlock()

	// Return if no routes are present
	if len(routes) == 0 {
		return m.internalRoutes
	}

	m.deleteRoutesWithID(routes[0].ID)

	for _, r := range routes {
		m.internalRoutes = append(m.internalRoutes, r)
	}

	return m.internalRoutes
}

// SetInternalRoutes sets the internal routes
func (m *Module) SetInternalRoutes(conf *config.Static) {
	m.Lock()
	defer m.Unlock()

	m.internalRoutes = conf.InternalRoutes
}

func (m *Module) deleteRoutesWithID(id string) {
	// Filter out those routes whose ids don't match
	routes := []*config.StaticRoute{}
	for _, r := range m.routes {
		if r.ID != id {
			routes = append(routes, r)
		}
	}
	m.routes = routes
}
