package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetProjectRoutes sets a projects routes
func (s *Manager) SetProjectRoutes(ctx context.Context, project string, c config.Routes) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	ingressRoutes := make(config.IngressRoutes)
	for _, route := range c {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceIngressRoute, route.ID)
		ingressRoutes[resourceID] = route
	}

	// Update the project's routes
	projectConfig.IngressRoutes = ingressRoutes
	if err := s.modules.Routing().SetProjectRoutes(project, ingressRoutes); err != nil {
		return http.StatusInternalServerError, err
	}

	for resourceID, route := range ingressRoutes {
		if err := s.store.SetResource(ctx, resourceID, route); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

// GetProjectRoutes gets all the routes for specified project config
func (s *Manager) GetProjectRoutes(ctx context.Context, project string) (int, interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	ingressRoutes := make(config.Routes, 0)
	for _, route := range projectConfig.IngressRoutes {
		ingressRoutes = append(ingressRoutes, route)
	}
	return http.StatusOK, ingressRoutes, nil
}

// SetProjectRoute adds a route in specified project config
func (s *Manager) SetProjectRoute(ctx context.Context, project, id string, c *config.Route, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	c.ID = id
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceIngressRoute, id)
	if projectConfig.IngressRoutes == nil {
		projectConfig.IngressRoutes = config.IngressRoutes{resourceID: c}
	} else {
		projectConfig.IngressRoutes[resourceID] = c
	}

	if err := s.modules.Routing().SetProjectRoutes(project, projectConfig.IngressRoutes); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetResource(ctx, resourceID, c); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteProjectRoute deletes a route from specified project config
func (s *Manager) DeleteProjectRoute(ctx context.Context, project, routeID string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceIngressRoute, routeID)
	_, ok := projectConfig.IngressRoutes[resourceID]
	if ok {
		delete(projectConfig.IngressRoutes, resourceID)

		// update the config
		if err := s.modules.Routing().SetProjectRoutes(project, projectConfig.IngressRoutes); err != nil {
			return http.StatusInternalServerError, err
		}

		if err := s.store.DeleteResource(ctx, resourceID); err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	}
	return http.StatusNotFound, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Route (%s) not found in config", routeID), nil, map[string]interface{}{})
}

// GetIngressRouting gets ingress routing from config
func (s *Manager) GetIngressRouting(ctx context.Context, project, routeID string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if routeID != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceIngressRoute, routeID)
		value, ok := projectConfig.IngressRoutes[resourceID]
		if ok {
			return http.StatusOK, []interface{}{value}, nil
		}
		return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("route with id (%s) does not exists in ingress routing", routeID), nil, nil)
	}

	routes := []interface{}{}
	for _, value := range projectConfig.IngressRoutes {
		routes = append(routes, value)
	}
	return http.StatusOK, routes, nil
}

// SetGlobalRouteConfig sets the project level ingress routing config
func (s *Manager) SetGlobalRouteConfig(ctx context.Context, project string, globalConfig *config.GlobalRoutesConfig, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the provided project's config
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Set config in project config object
	projectConfig.IngressGlobal = globalConfig

	// Update the routing module
	s.modules.Routing().SetGlobalConfig(globalConfig)

	// Finally lets store the config
	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceIngressGlobal, "global")
	if err := s.store.SetResource(ctx, resourceID, globalConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteGlobalRouteConfig sets the project level ingress routing config
func (s *Manager) DeleteGlobalRouteConfig(ctx context.Context, project string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the provided project's config
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Delete config in project config object
	projectConfig.IngressGlobal = &config.GlobalRoutesConfig{}

	// Update the routing module
	s.modules.Routing().SetGlobalConfig(projectConfig.IngressGlobal)

	// Finally lets store the config
	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceIngressGlobal, "global")
	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetGlobalRouteConfig returns the project level ingress routing config
func (s *Manager) GetGlobalRouteConfig(ctx context.Context, project string, params model.RequestParams) (int, interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the provided project's config
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, projectConfig.IngressGlobal, nil
}
