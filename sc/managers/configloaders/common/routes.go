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
			MatcherSetsRaw: getMatcherSet("GET", "/v1/config/*/*/*/single"),
			HandlersRaw:    []json.RawMessage{getHandler("config_get", "one")},
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("GET", "/v1/config/*/*/*/many"),
			HandlersRaw:    []json.RawMessage{getHandler("config_get", "many")},
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("DELETE", "/v1/config/*/*/*/single"),
			HandlersRaw:    []json.RawMessage{getHandler("config_delete", "one")},
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("DELETE", "/v1/config/*/*/*/many"),
			HandlersRaw:    []json.RawMessage{getHandler("config_delete", "many")},
		},
		caddyhttp.Route{
			Group:          "config",
			MatcherSetsRaw: getMatcherSet("POST", "/v1/config/*/*/*/single"),
			HandlersRaw:    []json.RawMessage{getHandler("config_post", "single")},
		},

		// Config routes
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("GET", "/v1/operation/*/*/*/single"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_get", "one")},
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("GET", "/v1/operation/*/*/*/many"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_get", "many")},
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("DELETE", "/v1/operation/*/*/*/single"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_delete", "one")},
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("DELETE", "/v1/operation/*/*/*/many"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_delete", "many")},
		},
		caddyhttp.Route{
			Group:          "operation",
			MatcherSetsRaw: getMatcherSet("POST", "/v1/operation/*/*/*/single"),
			HandlersRaw:    []json.RawMessage{getHandler("operation_post", "single")},
		},
	}
}
