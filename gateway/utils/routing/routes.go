package routing

import (
	"fmt"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type routeMapping map[string]config.Routes // The key here is the project name

func (r routeMapping) addProjectRoutes(project string, routes config.Routes) {
	r[project] = routes
}

func (r routeMapping) deleteProjectRoutes(project string) {
	delete(r, project)
}

func (r routeMapping) selectRoute(host, url string) (config.Route, error) {
	// Iterate over each project
	for _, routes := range r {
		// Iterate over each route of the project
		for _, route := range routes {
			// Skip if the hosts isn't present in the rule and hosts doesn't contain `*`
			if !utils.StringExists(route.Source.Hosts, host) && !utils.StringExists(route.Source.Hosts, "*") {
				continue
			}

			// TODO: add support for path parameters in routes
			switch route.Source.Type {
			case config.RoutePrefix:
				if strings.HasPrefix(url, route.Source.URL) {
					return route, nil
				}
			case config.RouteExact:
				if url == route.Source.URL {
					return route, nil
				}
			default:
				return config.Route{}, fmt.Errorf("invalid type (%s) provided for url matching", route.Source.Type)
			}
		}
	}

	return config.Route{}, fmt.Errorf("route not found for provided host (%s) and url (%s)", host, url)
}

func (r *Routing) selectRoute(host, url string) (config.Route, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.routes.selectRoute(host, url)
}
