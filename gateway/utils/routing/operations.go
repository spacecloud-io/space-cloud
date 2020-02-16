package routing

import "github.com/spaceuptech/space-cloud/gateway/config"

// SetProjectRoutes adds a project's routes to the global list of routes
func (r *Routing) SetProjectRoutes(project string, routes config.Routes) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.routes.addProjectRoutes(project, routes)
}

// DeleteProjectRoutes deletes a project's routes from the global list or routes
func (r *Routing) DeleteProjectRoutes(project string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.routes.deleteProjectRoutes(project)
}
