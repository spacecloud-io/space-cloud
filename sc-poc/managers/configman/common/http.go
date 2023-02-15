package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	for k, v := range config {
		gvr := source.GetResourceGVR(k)

		// Create route of GVR for List operation
		gvrPath := createConfigGVRPath(gvr)
		data := make(map[string]interface{})
		data["gvr"] = gvr

		// Route for List operation
		listRoute := caddyhttp.Route{
			Group:          "config_list",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{gvrPath}, []string{http.MethodGet}),
			HandlersRaw:    utils.GetCaddyHandler("config_list", data),
		}

		// Route for Apply operation
		applyRoute := caddyhttp.Route{
			Group:          "config_apply",
			MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{gvrPath}, []string{http.MethodPut}),
			HandlersRaw:    utils.GetCaddyHandler("config_apply", data),
		}

		configRoutes = append(configRoutes, listRoute, applyRoute)

		// Create route of each instances of GVR for Get, Apply and Delete operations
		for _, unstr := range v {
			name := unstr.GetName()
			namePath := createConfigNamePath(gvrPath, name)
			data["name"] = name

			// Route for Get operation
			getRoute := caddyhttp.Route{
				Group:          "config_get",
				MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{namePath}, []string{http.MethodGet}),
				HandlersRaw:    utils.GetCaddyHandler("config_get", data),
			}

			// Route for Delete operation
			deleteRoute := caddyhttp.Route{
				Group:          "config_delete",
				MatcherSetsRaw: utils.GetCaddyMatcherSet([]string{namePath}, []string{http.MethodDelete}),
				HandlersRaw:    utils.GetCaddyHandler("config_delete", data),
			}

			configRoutes = append(configRoutes, getRoute, deleteRoute)
		}
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

func createConfigGVRPath(gvr schema.GroupVersionResource) string {
	group := gvr.Group
	version := gvr.Version
	resource := gvr.Resource

	return fmt.Sprintf("/v1/config/%s/%s/%s", group, version, resource)
}

func createConfigNamePath(gvrPath string, name string) string {
	return fmt.Sprintf("%s/%s", gvrPath, name)
}
