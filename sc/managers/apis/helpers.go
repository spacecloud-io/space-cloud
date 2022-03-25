package apis

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/getkin/kin-openapi/openapi3"
)

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

func sanitizeURL(api *API) (string, map[string]string) {
	path := api.Path
	indexes := make(map[string]string)

	for _, param := range api.PathDef.Parameters {
		if param.Value.In == "path" {
			path, indexes = replacePathParam(path, param.Value.Name, indexes)
		}
	}

	return path, indexes
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

func getPathParams(ogURL, receivedURL string, indexes map[string]string) map[string]string {
	var index int
	m := map[string]string{}

	arr1 := strings.Split(ogURL, "/")
	arr2 := strings.Split(receivedURL, "/")
	for i, segment := range arr1 {
		if segment == "*" {
			m[indexes[strconv.Itoa(index)]] = arr2[i]
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
