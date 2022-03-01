package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// ConfigGetHandler is a module to create config GET handlers
type ConfigGetHandler struct {
	Operation string `json:"operation,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (ConfigGetHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_get_handler",
		New: func() caddy.Module { return new(ConfigGetHandler) },
	}
}

func (rb ConfigGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s,  Operation: %s \n", r.Method, r.URL, rb.Operation)))

	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*ConfigGetHandler)(nil)
