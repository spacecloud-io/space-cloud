package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// ConfigPostHandler is a module to create config POST handlers
type ConfigPostHandler struct {
	Operation string `json:"operation,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (ConfigPostHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_post_handler",
		New: func() caddy.Module { return new(ConfigPostHandler) },
	}
}

func (rb ConfigPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s,  Operation: %s \n", r.Method, r.URL, rb.Operation)))

	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*ConfigPostHandler)(nil)
