package apis

import (
	"fmt"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/utils"
)

// The necesary global objects to hold all registered apps
var (
	appsLock sync.RWMutex

	registeredApps apps
)

// RegisterApp marks the app as having routers
func RegisterApp(name string, priority int) {
	appsLock.Lock()
	defer appsLock.Unlock()

	registeredApps = append(registeredApps, app{name, priority})
	registeredApps.sort()
}

func generateOpenAPIDocAndAPIs(ctx caddy.Context) (*openapi3.T, []*API, error) {
	var allAPIs []*API

	// Load the paths of each app
	paths := make(openapi3.Paths)
	for _, a := range registeredApps {
		// Get the app
		appTmp, err := ctx.App(a.name)
		if err != nil {
			return nil, nil, err
		}

		// See if the app implements the required methods
		app, ok := appTmp.(App)
		if !ok {
			return nil, nil, fmt.Errorf("app '%s' does not implement the method 'GetRoutes'", a.name)
		}

		// Merge the paths returned by the app
		for _, api := range app.GetRoutes() {
			mergePaths(paths, api.Path, api.PathDef)

			api.app = a.name
			allAPIs = append(allAPIs, api)
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "SpaceCloud exposed APIs",
			Description: "Specification of all the APIs exposed by the various modules of SpaceCloud",
			Version:     "v0.22.0",
		},
		Paths:      paths,
		Components: openapi3.NewComponents(),
	}, allAPIs, nil
}

func makeSubRouter(ctx caddy.Context, allAPIs []*API) (caddyhttp.Handler, error) {
	routeList := make(caddyhttp.RouteList, 0)

	for _, api := range allAPIs {
		// Substitute all path parameters with '*'
		path, indexes := sanitizeURL(api)

		// Get the methods to be used
		methods := getMethods(api)

		handlerObj := caddyhttp.Route{
			Group:          api.app,
			MatcherSetsRaw: utils.GetCaddyMatcherSet(path, methods),
			HandlersRaw: utils.GetCaddyHandler("api", map[string]interface{}{
				"path":    path,
				"indexes": indexes,
				"app":     api.app,
				"op":      api.Op,
			}),
		}
		routeList = append(routeList, handlerObj)
	}

	// Provision all the handlers
	if err := routeList.Provision(ctx); err != nil {
		return nil, err
	}

	return routeList.Compile(emptyHandler), nil
}
