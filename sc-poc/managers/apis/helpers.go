package apis

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/utils"
)

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
		for _, api := range app.GetAPIRoutes() {
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
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{path}, methods),
			HandlersRaw: utils.GetCaddyHandler("api", map[string]interface{}{
				"path":    path,
				"indexes": indexes,
				"app":     api.app,
				"name":    api.Name,
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

func getMethods(api *API) []string {
	var methods []string

	if api.PathDef.Post != nil {
		methods = append(methods, http.MethodPost)
	}

	if api.PathDef.Put != nil {
		methods = append(methods, http.MethodPut)
	}

	if api.PathDef.Patch != nil {
		methods = append(methods, http.MethodPatch)
	}

	if api.PathDef.Get != nil {
		methods = append(methods, http.MethodGet)
	}

	if api.PathDef.Options != nil {
		methods = append(methods, http.MethodOptions)
	}

	if api.PathDef.Delete != nil {
		methods = append(methods, http.MethodDelete)
	}

	return methods
}

func sanitizeURL(api *API) (string, []string) {
	path := api.Path

	// Make an index map. A map is used since we replace one param at a time in a random order.
	indexMap := make(map[string]string)

	// Loop over each path param and replace them with an `*` one by one
	for _, param := range api.PathDef.Parameters {
		if param.Value.In == "path" {
			path, indexMap = replacePathParam(path, param.Value.Name, indexMap)
		}
	}

	// Convert the index map to an array
	finalIndexes := make([]string, len(indexMap))
	for k, v := range indexMap {
		i, _ := strconv.Atoi(k)
		finalIndexes[i] = v
	}

	return path, finalIndexes
}

func replacePathParam(path, param string, indexes map[string]string) (string, map[string]string) {
	var index, start int
	var subtract int
	newPath := path
	for i, r := range path {
		// Add the index whenever a path param is encountered
		if r == '*' || r == '{' {
			index++
		}

		// Update the start index on encountering an '{'
		if r == '{' {
			start = i + 1
		}

		// Check if the param matches with what was specified. We will replace the param with an '*'
		// and store the index location for future reference
		if r == '}' {
			start2 := start - subtract
			stop2 := i - subtract
			if newPath[start2:stop2] == param {
				newPath = newPath[:start2-1] + "*" + newPath[stop2+1:]
				indexes[strconv.Itoa(index-1)] = param
				subtract += i - start + 1
			}
		}
	}

	return newPath, indexes
}

func getPathParams(configuredURL, receivedURL string, indexes []string) map[string]string {
	var index int
	m := map[string]string{}

	arr1 := strings.Split(configuredURL, "/")
	arr2 := strings.Split(receivedURL, "/")
	for i, segment := range arr1 {
		if segment == "*" {
			m[indexes[index]] = arr2[i]
			index++
		}
	}

	return m
}

func mergePaths(paths openapi3.Paths, url string, pathDef *openapi3.PathItem) {
	existingPathDef, p := paths[url]
	if !p {
		paths[url] = pathDef
		return
	}

	if pathDef.Description != "" {
		existingPathDef.Description = pathDef.Description
	}
	if len(pathDef.Parameters) > 0 {
		existingPathDef.Parameters = pathDef.Parameters
	}
	if pathDef.Connect != nil {
		existingPathDef.Connect = pathDef.Connect
	}
	if pathDef.Post != nil {
		existingPathDef.Post = pathDef.Post
	}
	if pathDef.Put != nil {
		existingPathDef.Put = pathDef.Put
	}
	if pathDef.Patch != nil {
		existingPathDef.Patch = pathDef.Patch
	}
	if pathDef.Get != nil {
		existingPathDef.Get = pathDef.Get
	}
	if pathDef.Options != nil {
		existingPathDef.Options = pathDef.Options
	}
	if pathDef.Head != nil {
		existingPathDef.Head = pathDef.Head
	}
	if pathDef.Delete != nil {
		existingPathDef.Delete = pathDef.Delete
	}
}

var emptyHandler caddyhttp.Handler = caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Route not found for provided host (%s), method (%s) and url (%s)", r.Host, r.Method, r.URL.Path)
	return nil
})
