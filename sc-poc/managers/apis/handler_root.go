package apis

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

// RootAPIHandler is a module to inject REST APIs into caddy http app
type RootAPIHandler struct {
	// For internal use
	logger     *zap.Logger
	handler    caddyhttp.Handler
	openapiDoc *openapi3.T
}

// CaddyModule returns the Caddy module information.
func (RootAPIHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_root_api_handler",
		New: func() caddy.Module { return new(RootAPIHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *RootAPIHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Load the routes of each registered app
	doc, allAPIs, err := generateOpenAPIDocAndAPIs(ctx)
	if err != nil {
		h.logger.Error("Unable to generate OpenAPI doc for SpaceCloud APIs", zap.Error(err))
		return err
	}

	// Validate the openapi doc
	if err := doc.Validate(ctx.Context); err != nil {
		h.logger.Error("Invalid OpenAPI doc generated for SpaceCloud APIs", zap.Error(err))
		return err
	}

	// Create a handler
	handler, err := makeSubRouter(ctx, allAPIs)
	if err != nil {
		h.logger.Error("Unable to create root handler for SpaceCloud APIs", zap.Error(err))
		return err
	}

	h.handler = handler
	h.openapiDoc = doc
	return nil
}

// ServeHTTP handles the http request
func (h *RootAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Check if request was to send the OpenAPI document
	// TODO: Make this endoint optional
	if r.Method == http.MethodGet && r.URL.Path == "/v1/api/openapi.json" {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(h.openapiDoc)
	}

	// Invoke the handler
	return h.handler.ServeHTTP(w, r)
}

// Interface guard
var _ caddy.Provisioner = (*RootAPIHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*RootAPIHandler)(nil)
