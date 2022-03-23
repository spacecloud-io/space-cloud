package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
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

var emptyHandler caddyhttp.Handler = caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Route not found for provided host (%s), method (%s) and url (%s)", r.Host, r.Method, r.URL.Path)
	return nil
})
