package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// ConfigDeleteHandler is a module to create config Delete handlers
type ConfigDeleteHandler struct {
	logger    *zap.Logger
	appLoader loadApp
}

// CaddyModule returns the Caddy module information.
func (ConfigDeleteHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_delete_handler",
		New: func() caddy.Module { return new(ConfigDeleteHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *ConfigDeleteHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	h.appLoader = ctx.App
	return nil
}

// ServeHTTP handles the http request
func (h *ConfigDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL)))

	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigDeleteHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigDeleteHandler)(nil)
