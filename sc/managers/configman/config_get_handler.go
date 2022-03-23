package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// ConfigGetHandler is a module to create config GET handlers
type ConfigGetHandler struct {
	logger    *zap.Logger
	appLoader loadApp
}

// CaddyModule returns the Caddy module information.
func (ConfigGetHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_get_handler",
		New: func() caddy.Module { return new(ConfigGetHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *ConfigGetHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	h.appLoader = ctx.App
	return nil
}

// ServeHTTP handles the http request
func (h *ConfigGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL)))
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigGetHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigGetHandler)(nil)
