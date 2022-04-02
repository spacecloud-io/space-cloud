package common

import (
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

		// Open API for the config and operation endpoints
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/config/openapi.json", []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("config_openapi", nil),
		},

		// Config routes
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/config/*", []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("config_get", nil),
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/config/*", []string{http.MethodDelete}),
			HandlersRaw:    utils.GetCaddyHandler("config_delete", nil),
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/config/*", []string{http.MethodPost}),
			HandlersRaw:    utils.GetCaddyHandler("config_apply", nil),
		},

		// Operation routes
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/operation/*", []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("operation_get", nil),
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/operation/*", []string{http.MethodDelete}),
			HandlersRaw:    utils.GetCaddyHandler("operation_delete", nil),
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/v1/operation/*", []string{http.MethodPost}),
			HandlersRaw:    utils.GetCaddyHandler("operation_apply", nil),
		},

		// API route handler
		caddyhttp.Route{
			Group:          "api",
			MatcherSetsRaw: utils.GetCaddyMatcherSet("/*", nil),
			HandlersRaw:    utils.GetCaddyHandler("root_api", nil),
		},
	}
}
