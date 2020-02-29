package syncman

import (
	"context"

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
func (s *Manager) SetProjectRoute(ctx context.Context, project string, c *config.Route) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	doesExist := false
	for _, route := range projectConfig.Modules.Routes {
		if route.ID == c.ID {
			route.Source = c.Source
			route.Targets = c.Targets
			doesExist = true
		}
	}
	if !doesExist {
		projectConfig.Modules.Routes = append(projectConfig.Modules.Routes, c)
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

	routes := projectConfig.Modules.Routes
	for index, route := range routes {
		if route.ID == routeID {
			// delete the route at specified index
			routes[index] = routes[len(routes)-1]
			projectConfig.Modules.Routes = routes[:len(routes)-1]
			// update the config
			return s.setProject(ctx, projectConfig)
		}
	}
	return nil
}
