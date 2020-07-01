package routing

import (
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetProjectRoutes adds a project's routes to the global list of routes
func (r *Routing) SetProjectRoutes(project string, routes config.Routes) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Delete all templates for this project
	for k := range r.goTemplates {
		if strings.HasPrefix(k, project) {
			delete(r.goTemplates, k)
		}
	}

	// Add projects to the routes object and generate go templates
	for _, route := range routes {
		route.Project = project
		route.Modify.Tmpl = config.EndpointTemplatingEngineGo

		// Parse request template
		if route.Modify.ReqTmpl != "" {
			if err := r.createGoTemplate("request", project, route.ID, route.Modify.ReqTmpl); err != nil {
				return err
			}
		}

		// Parse response template
		if route.Modify.ResTmpl != "" {
			if err := r.createGoTemplate("response", project, route.ID, route.Modify.ResTmpl); err != nil {
				return err
			}
		}
	}

	r.routes.addProjectRoutes(project, routes)
	return nil
}

// DeleteProjectRoutes deletes a project's routes from the global list or routes
func (r *Routing) DeleteProjectRoutes(project string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.routes.deleteProjectRoutes(project)
}
