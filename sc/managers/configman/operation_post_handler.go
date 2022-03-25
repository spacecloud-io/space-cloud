package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// OperationPostHandler is a module to create operation POST handlers
type OperationPostHandler struct {
	Operation string `json:"operation,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (OperationPostHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_operation_apply_handler",
		New: func() caddy.Module { return new(OperationPostHandler) },
	}
}

// ServeHTTP handles the http request
func (rb OperationPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s,  Operation: %s \n", r.Method, r.URL, rb.Operation)))

	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*OperationPostHandler)(nil)
