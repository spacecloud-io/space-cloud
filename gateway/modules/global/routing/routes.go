package routing

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (r *Routing) addProjectRoutes(project string, routes config.Routes) {
	r.deleteProjectRoutes(project)
	r.routes = append(r.routes, routes...)
	sort.Stable(r.routes) // This will sort the array in place
}

func (r *Routing) deleteProjectRoutes(project string) {
	newRoutes := make(config.Routes, 0)
	for _, route := range r.routes {
		if route.Project != project {
			newRoutes = append(newRoutes, route)
		}
	}
	r.routes = newRoutes
}

func (r *Routing) selectRoute(ctx context.Context, host, method, url string) (*config.Route, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	// Iterate over each route
	for _, route := range r.routes {
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
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type (%s) provided for url matching", route.Source.Type), nil, nil)
		}
	}

	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Route not found for provided host (%s), method (%s) and url (%s)", host, method, url), nil, nil)
}
