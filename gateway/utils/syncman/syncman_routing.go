package syncman

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetProjectRoutes sets a projects routes
func (s *Manager) SetProjectRoutes(ctx context.Context, project string, c config.Routes) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// Update the project's routes
	projectConfig.Modules.Routes = c
	if err := s.routing.SetProjectRoutes(project, c); err != nil {
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// GetProjectRoutes gets all the routes for specified project config
func (s *Manager) GetProjectRoutes(ctx context.Context, project string) (config.Routes, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}

	return projectConfig.Modules.Routes, nil
}

// SetProjectRoute adds a route in specified project config
func (s *Manager) SetProjectRoute(ctx context.Context, project, id string, c *config.Route) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	c.ID = id
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
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

	if err := s.routing.SetProjectRoutes(project, projectConfig.Modules.Routes); err != nil {
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// DeleteProjectRoute deletes a route from specified project config
func (s *Manager) DeleteProjectRoute(ctx context.Context, project, routeID string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	for index, route := range projectConfig.Modules.Routes {
		if route.ID == routeID {
			// delete the route at specified index
			projectConfig.Modules.Routes[index] = projectConfig.Modules.Routes[len(projectConfig.Modules.Routes)-1]
			projectConfig.Modules.Routes = projectConfig.Modules.Routes[:len(projectConfig.Modules.Routes)-1]

			// update the config
			if err := s.routing.SetProjectRoutes(project, projectConfig.Modules.Routes); err != nil {
				return err
			}

			return s.setProject(ctx, projectConfig)
		}
	}
	return nil
}

// GetIngressRouting gets ingress routing from config
func (s *Manager) GetIngressRouting(ctx context.Context, project, routeID string) ([]interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return nil, err
	}
	if routeID != "*" {
		for _, value := range projectConfig.Modules.Routes {
			if routeID == value.ID {
				return []interface{}{value}, nil
			}
		}
		return nil, fmt.Errorf("route id (%s) not present in config", routeID)
	}

	routes := []interface{}{}
	for _, value := range projectConfig.Modules.Routes {
		routes = append(routes, value)
	}
	return routes, nil
}
