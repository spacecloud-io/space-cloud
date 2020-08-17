package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetProjectRoutes sets a projects routes
func (s *Manager) SetProjectRoutes(ctx context.Context, project string, c config.Routes) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Update the project's routes
	projectConfig.Modules.Routes = c
	if err := s.modules.Routing().SetProjectRoutes(project, c); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetProjectRoutes gets all the routes for specified project config
func (s *Manager) GetProjectRoutes(ctx context.Context, project string) (int, interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, projectConfig.Modules.Routes, nil
}

// SetProjectRoute adds a route in specified project config
func (s *Manager) SetProjectRoute(ctx context.Context, project, id string, c *config.Route, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	c.ID = id
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	doesExist := false
	for _, route := range projectConfig.Modules.Routes {
		if id == route.ID {
			route.Source = c.Source
			route.Targets = c.Targets
			route.Rule = c.Rule
			route.Modify = c.Modify
			doesExist = true
		}
	}
	if !doesExist {
		projectConfig.Modules.Routes = append(projectConfig.Modules.Routes, c)
	}

	if err := s.modules.Routing().SetProjectRoutes(project, projectConfig.Modules.Routes); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteProjectRoute deletes a route from specified project config
func (s *Manager) DeleteProjectRoute(ctx context.Context, project, routeID string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	for index, route := range projectConfig.Modules.Routes {
		if route.ID == routeID {
			// delete the route at specified index
			projectConfig.Modules.Routes[index] = projectConfig.Modules.Routes[len(projectConfig.Modules.Routes)-1]
			projectConfig.Modules.Routes = projectConfig.Modules.Routes[:len(projectConfig.Modules.Routes)-1]

			// update the config
			if err := s.modules.Routing().SetProjectRoutes(project, projectConfig.Modules.Routes); err != nil {
				return http.StatusInternalServerError, err
			}

			if err := s.setProject(ctx, projectConfig); err != nil {
				return http.StatusInternalServerError, err
			}

			return http.StatusOK, nil
		}
	}
	return http.StatusNotFound, utils.LogError(fmt.Sprintf("Route (%s) not found", routeID), "syncman", "ingres-route-delete", nil)
}

// GetIngressRouting gets ingress routing from config
func (s *Manager) GetIngressRouting(ctx context.Context, project, routeID string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if routeID != "*" {
		for _, value := range projectConfig.Modules.Routes {
			if routeID == value.ID {
				return http.StatusOK, []interface{}{value}, nil
			}
		}
		return http.StatusBadRequest, nil, fmt.Errorf("route id (%s) not present in config", routeID)
	}

	routes := []interface{}{}
	for _, value := range projectConfig.Modules.Routes {
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
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Set config in project config object
	projectConfig.Modules.GlobalRoutes = globalConfig

	// Update the routing module
	s.modules.Routing().SetGlobalConfig(globalConfig)

	// Finally lets store the config
	if err := s.setProject(ctx, projectConfig); err != nil {
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
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, projectConfig.Modules.GlobalRoutes, nil
}
