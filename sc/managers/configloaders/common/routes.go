package common

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func getRootRoutes() caddyhttp.RouteList {
	return caddyhttp.RouteList{
		// Config routes
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("GET", "/v1/config/*"),
			HandlersRaw:    []json.RawMessage{getHandler("config_get")},
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("DELETE", "/v1/config/*"),
			HandlersRaw:    []json.RawMessage{getHandler("config_delete")},
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("POST", "/v1/config/*"),
			HandlersRaw:    []json.RawMessage{getHandler("config_apply")},
		},

		// Operation routes
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("GET", "/v1/operation/*"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_get")},
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("DELETE", "/v1/operation/*"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_delete")},
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("POST", "/v1/operation/*"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_apply")},
		},
	}
}
