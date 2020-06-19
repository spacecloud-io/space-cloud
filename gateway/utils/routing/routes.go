package routing

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type routeMapping map[string]config.Routes // The key here is the project name

const (
	module         string = "ingress-route"
	handleRequest  string = "handle-request"
	handleResponse string = "handle-response"
)

func (r routeMapping) addProjectRoutes(project string, routes config.Routes) {
	sort.Stable(routes) // This will sort the array in place

	// Store the projects
	r[project] = routes
}

func (r routeMapping) deleteProjectRoutes(project string) {
	delete(r, project)
}

func (r routeMapping) selectRoute(host, method, url string) (*config.Route, error) {
	// Iterate over each project
	for _, routes := range r {
		// Iterate over each route of the project
		for _, route := range routes {
			// Skip if the hosts isn't present in the rule and hosts doesn't contain `*`
			if !utils.StringExists(route.Source.Hosts, host) && !utils.StringExists(route.Source.Hosts, "*") {
				continue
			}

			// Skip if the method doesn't match
			if len(route.Source.Methods) > 0 && !utils.StringExists(route.Source.Methods, "*") && !utils.StringExists(route.Source.Methods, method) {
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
				return nil, fmt.Errorf("invalid type (%s) provided for url matching", route.Source.Type)
			}
		}
	}

	return nil, fmt.Errorf("route not found for provided host (%s), method (%s) and url (%s)", host, method, url)
}

func (r *Routing) selectRoute(host, method, url string) (*config.Route, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.routes.selectRoute(host, method, url)
}
