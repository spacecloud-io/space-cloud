package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
)

func prepareHTTPHandlerApp(config ConfigType) []byte {
	port := viper.GetInt("caddy.port")

	httpsPort := 0
	listen := []string{":" + strconv.Itoa(port)}
	// if sshCert != "none" && sshKey != "none" {
	// 	httpsPort = port + 4
	// 	listen = []string{":" + strconv.Itoa(httpsPort)}
	// 	port = 0
	// }

	httpConfig := caddyhttp.App{
		HTTPPort:  port,
		HTTPSPort: httpsPort,
		Servers: map[string]*caddyhttp.Server{
			"default": {
				Listen: listen,
				Routes: getRootRoutes(config),
			},
		},
	}

	data, _ := json.Marshal(httpConfig)
	return data
}

func getRootRoutes(config ConfigType) caddyhttp.RouteList {
	return caddyhttp.RouteList{
		// Routes for CORS
		caddyhttp.Route{
			Group:       "cors",
			HandlersRaw: utils.GetCaddyHandler("cors", nil),
		},

		// TODO: Fix this
		// // Root middleware for all routes
		// caddyhttp.Route{
		// 	Group:       "req-params",
		// 	HandlersRaw: utils.GetCaddyHandler("req_params", nil),
		// },

		// Config route handlers
		getConfigRoutes(config),

		// API route handler
		getAPIRoutes(),
	}
}

func getAPIRoutes() caddyhttp.Route {
	// Make route list for the sub router
	routeList := caddyhttp.RouteList{
		caddyhttp.Route{
			Group:       "api_auth",
			HandlersRaw: utils.GetCaddyHandler("auth_verify", nil),
		},
		caddyhttp.Route{
			Group:       "api_route",
			HandlersRaw: utils.GetCaddyHandler("root_api", nil),
		},
	}

	// Create matcher and handler for subroute
	handler := map[string]interface{}{
		"handler": "subroute",
		"routes":  routeList,
	}
	handlerRaw, _ := json.Marshal(handler)

	return caddyhttp.Route{
		Group:       "api",
		HandlersRaw: []json.RawMessage{handlerRaw},
	}
}

func getConfigRoutes(config ConfigType) caddyhttp.Route {
	// Make route list for the sub router
	configRoutes := caddyhttp.RouteList{}
	for k := range config {
		gvr := source.GetResourceGVR(k)

		// Create route of GVR for List operation
		path := createConfigPath(gvr)
		data := make(map[string]interface{})
		data["gvr"] = gvr

		// Route for List/Get operation
		getRoute := caddyhttp.Route{
			Group:          "config_get",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{path}, []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("config_get", data),
		}

		// Route for Apply operation
		applyRoute := caddyhttp.Route{
			Group:          "config_apply",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{path}, []string{http.MethodPut}),
			HandlersRaw:    utils.GetCaddyHandler("config_apply", data),
		}

		// Route for Delete operation
		deleteRoute := caddyhttp.Route{
			Group:          "config_delete",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{path}, []string{http.MethodDelete}),
			HandlersRaw:    utils.GetCaddyHandler("config_delete", data),
		}

		configRoutes = append(configRoutes, getRoute, applyRoute, deleteRoute)
	}

	// Make handler for subroute
	handler := map[string]interface{}{
		"handler": "subroute",
		"routes":  configRoutes,
	}
	handlerRaw, _ := json.Marshal(handler)

	return caddyhttp.Route{
		Group:       "config",
		HandlersRaw: []json.RawMessage{handlerRaw},
	}
}

func createConfigPath(gvr schema.GroupVersionResource) string {
	group := gvr.Group
	version := gvr.Version
	resource := gvr.Resource

	return fmt.Sprintf("/v1/config/%s/%s/%s/*", group, version, resource)
}