package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// OperationGetHandler is a module to create operation GET handlers
type OperationGetHandler struct {
	Operation string `json:"operation,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (OperationGetHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_operation_get_handler",
		New: func() caddy.Module { return new(OperationGetHandler) },
	}
}

// ServeHTTP handles the http request
func (rb OperationGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s,  Operation: %s \n", r.Method, r.URL, rb.Operation)))

	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*OperationGetHandler)(nil)
