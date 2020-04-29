package manager

import (
	"github.com/spaceuptech/space-cloud/runner/model"
)

// SetServiceRoutes sets the routes, saves config in manager & adjusts the ports as required
func (m *Manager) SetServiceRoutes(projectID, serviceID string, r model.Routes) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.serviceRoutes[getConfigKey(projectID, serviceID)] = r

	if err := m.writeConfigToFile(); err != nil {
		return err
	}

	m.adjustProxyServers()
	return nil
}

// SetServiceRouteIfNotExists is used to set the service routes if there exists no services (on start)
func (m *Manager) SetServiceRouteIfNotExists(projectID, serviceID, version string, ports []model.Port) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := getConfigKey(projectID, serviceID)
	if _, p := m.serviceRoutes[key]; p {
		// Simply return if the key already exists. We only want to do this for new services.
		return nil
	}

	routes := make(model.Routes, len(ports))
	for i, port := range ports {
		routes[i] = &model.Route{
			Source: model.RouteSource{Port: port.Port},
			Targets: []model.RouteTarget{{
				Type:    model.RouteTargetVersion,
				Version: version,
				Port:    port.Port,
				Weight:  100,
			}},
		}
	}
	m.serviceRoutes[key] = routes

	if err := m.writeConfigToFile(); err != nil {
		return err
	}

	m.adjustProxyServers()
	return nil
}

// GetServiceRoutes returns a map of routes for the required projectID
func (m *Manager) GetServiceRoutes(projectID string) (map[string]model.Routes, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	serviceConfig := map[string]model.Routes{}
	for k, routes := range m.serviceRoutes {
		pID, serviceID := getProjectAndServiceIDFromKey(k)
		if pID == projectID {
			// Don't forget to set the service id of the routes
			for _, r := range routes {
				r.ID = serviceID
			}
			serviceConfig[serviceID] = routes
		}
	}

	return serviceConfig, nil
}

// DeleteServiceRoutes deletes a particular service based on projectID and serviceID
func (m *Manager) DeleteServiceRoutes(projectID, serviceID string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.serviceRoutes, getConfigKey(projectID, serviceID))

	if err := m.writeConfigToFile(); err != nil {
		return err
	}

	m.adjustProxyServers()
	return nil
}
