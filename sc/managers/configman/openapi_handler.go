package configman

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

// OpenAPIHandler is a module to create config Delete handlers
type OpenAPIHandler struct {
	logger    *zap.Logger
	configMan *ConfigMan
}

// CaddyModule returns the Caddy module information.
func (OpenAPIHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_openapi_handler",
		New: func() caddy.Module { return new(OpenAPIHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *OpenAPIHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Get the configman app
	app, _ := ctx.App("configman")
	h.configMan = app.(*ConfigMan)
	return nil
}

// ServeHTTP handles the http request
func (h *OpenAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	// Load all the type definitions
	operationTypeDefs := h.configMan.GetOperationTypes()
	configTypeDefs := h.configMan.GetConfigTypes()

	// Create the open api doc
	openapiDoc := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "SpaceCloud config and operation APIs",
			Description: "Specification of all the config and operation APIs exposed by SpaceCloud",
			Version:     "v0.22.0",
		},
		Components: openapi3.NewComponents(),
		Paths:      make(openapi3.Paths),
	}

	// Add operation paths to openapi doc
	for module, types := range operationTypeDefs {
		_ = addOperationToOpenAPIDoc(openapiDoc, module, types)
	}

	// Add config paths to openapi doc
	for module, types := range configTypeDefs {
		addConfigToOpenAPIDoc(openapiDoc, module, types)
	}

	// Set the headers and status code
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the response body
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(openapiDoc)
}

// Interface guard
var _ caddy.Provisioner = (*OpenAPIHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*OpenAPIHandler)(nil)
