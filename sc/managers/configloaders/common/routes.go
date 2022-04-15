package common

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/utils"
)

func getRootRoutes() caddyhttp.RouteList {
	return caddyhttp.RouteList{
		// Routes for CORS
		caddyhttp.Route{
			Group:       "cors",
			HandlersRaw: utils.GetCaddyHandler("cors", nil),
		},

		// Root middleware for all routes
		caddyhttp.Route{
			Group:       "req-params",
			HandlersRaw: utils.GetCaddyHandler("req_params", nil),
		},

		// Config & Operation handlers
		getConfigRoutes(),

		// API route handler
		caddyhttp.Route{
			Group:          "api",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/*"}, nil),
			HandlersRaw:    utils.GetCaddyHandler("root_api", nil),
		},
	}
}

func getConfigRoutes() caddyhttp.Route {
	// Make route list for the sub router
	routeList := caddyhttp.RouteList{
		// Open API for the config and operation endpoints
		caddyhttp.Route{
			Group:          "openapi",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/v1/config/openapi.json"}, []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("config_openapi", nil),
		},

		// Admin auth middleware
		caddyhttp.Route{
			Group:       "auth",
			HandlersRaw: utils.GetCaddyHandler("admin_auth", nil),
		},

		// Config routes
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/v1/config/*"}, []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("config_get", nil),
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/v1/config/*"}, []string{http.MethodDelete}),
			HandlersRaw:    utils.GetCaddyHandler("config_delete", nil),
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/v1/config/*"}, []string{http.MethodPost}),
			HandlersRaw:    utils.GetCaddyHandler("config_apply", nil),
		},

		// Operation routes
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/v1/operation/*"}, []string{}),
			HandlersRaw:    utils.GetCaddyHandler("operation", nil),
		},
	}

	// Create matcher for subroute

	// Make handler for subroute
	handler := map[string]interface{}{
		"handler": "subroute",
		"routes":  routeList,
	}
	handlerRaw, _ := json.Marshal(handler)

	return caddyhttp.Route{
		Group:          "config",
		MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{"/v1/config/*", "/v1/operation/*"}, nil),
		HandlersRaw:    []json.RawMessage{handlerRaw},
	}
}
