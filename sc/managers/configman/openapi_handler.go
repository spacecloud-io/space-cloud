package configman

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// OpenAPIHandler is a module to create config Delete handlers
type OpenAPIHandler struct {
	logger *zap.Logger
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
	return nil
}

// ServeHTTP handles the http request
func (h *OpenAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	helpers.Logger.LogInfo(helpers.GetRequestID(r.Context()), "Response", map[string]interface{}{"statusCode": http.StatusOK})

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
