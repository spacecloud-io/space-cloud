package apis

import (
	"encoding/json"
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
	schemas := openapi3.Schemas{}

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
			api.app = a.name
			allAPIs = append(allAPIs, api)

			if api.OpenAPI != nil {
				mergePaths(paths, api.Path, api.OpenAPI.PathDef)

				// Check if any schemas were exposed
				for k, v := range api.OpenAPI.Schemas {
					schemas[k] = v
				}
			}
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "SpaceCloud exposed APIs",
			Description: "Specification of all the APIs exposed by the various modules of SpaceCloud",
			Version:     "v0.22.0",
		},
		Paths: paths,
		Components: openapi3.Components{
			Schemas: schemas,
			SecuritySchemes: openapi3.SecuritySchemes{
				"bearerAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:         "http",
						Scheme:       "bearer",
						BearerFormat: "JWT",
					},
				},
			},
		},
	}, allAPIs, nil
}

func makeSubRouter(ctx caddy.Context, allAPIs []*API) (caddyhttp.Handler, error) {
	routeList := make(caddyhttp.RouteList, 0)

	for _, api := range allAPIs {
		// Substitute all path parameters with '*'
		path, indexes := sanitizeURL(api)

		// Get the methods to be used
		methods := getMethods(api)

		handlerObj := prepareRoute(api, path, methods, indexes)
		routeList = append(routeList, handlerObj)
	}

	// Provision all the handlers
	if err := routeList.Provision(ctx); err != nil {
		return nil, err
	}

	return routeList.Compile(emptyHandler), nil
}

func prepareRoute(api *API, path string, methods, indexes []string) caddyhttp.Route {
	// Create the route for this api
	apiRoute := caddyhttp.Route{
		Group:          api.app,
		MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{path}, methods),
		HandlersRaw:    make([]json.RawMessage, 0, len(api.Plugins)+1),
	}

	// TODO: add a handler to extract all the variables
	// TODO: add a handler to validate json schema of extracted variables

	// First we add the handlers for all the plugins
	for _, p := range api.Plugins {
		var params map[string]interface{}
		if len(p.Params.Raw) > 0 {
			_ = json.Unmarshal(p.Params.Raw, &params)
		}
		apiRoute.HandlersRaw = append(apiRoute.HandlersRaw, utils.GetCaddyHandler(p.Driver, params)...)
	}

	// Finally comes the main api route
	apiHandler := utils.GetCaddyHandler("api", map[string]interface{}{
		"path":    path,
		"indexes": indexes,
		"app":     api.app,
		"name":    api.Name,
	})
	apiRoute.HandlersRaw = append(apiRoute.HandlersRaw, apiHandler...)

	return apiRoute
}

func getMethods(api *API) []string {
	var methods []string

	if api.OpenAPI == nil {
		methods = append(methods, http.MethodGet, http.MethodPost)
		return methods
	}

	if api.OpenAPI.PathDef.Post != nil {
		methods = append(methods, http.MethodPost)
	}

	if api.OpenAPI.PathDef.Put != nil {
		methods = append(methods, http.MethodPut)
	}

	if api.OpenAPI.PathDef.Patch != nil {
		methods = append(methods, http.MethodPatch)
	}

	if api.OpenAPI.PathDef.Get != nil {
		methods = append(methods, http.MethodGet)
	}

	if api.OpenAPI.PathDef.Options != nil {
		methods = append(methods, http.MethodOptions)
	}

	if api.OpenAPI.PathDef.Delete != nil {
		methods = append(methods, http.MethodDelete)
	}

	return methods
}

func sanitizeURL(api *API) (string, []string) {
	// Simply return if openapi definition isn't provided
	if api.OpenAPI == nil {
		return api.Path, []string{}
	}

	path := api.Path

	// Make an index map. A map is used since we replace one param at a time in a random order.
	indexMap := make(map[string]string)

	// Loop over each path param and replace them with an `*` one by one
	for _, param := range api.OpenAPI.PathDef.Parameters {
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
		newPathDef := new(openapi3.PathItem)
		d, _ := pathDef.MarshalJSON()
		_ = newPathDef.UnmarshalJSON(d)
		paths[url] = newPathDef
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
