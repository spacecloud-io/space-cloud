package configman

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// OperationDeleteHandler is a module to create operation Delete handlers
type OperationDeleteHandler struct {
	Operation string `json:"operation,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (OperationDeleteHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_operation_delete_handler",
		New: func() caddy.Module { return new(OperationDeleteHandler) },
	}
}

// ServeHTTP handles the http request
func (rb OperationDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte(fmt.Sprintf("Method: %s, Path: %s,  Operation: %s \n", r.Method, r.URL, rb.Operation)))

	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*OperationDeleteHandler)(nil)
