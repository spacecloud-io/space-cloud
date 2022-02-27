package handlers

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Handler{})
}

// Handler is a middleware for manipulating the request body.
type Handler struct {
}

// CaddyModule returns the Caddy module information.
func (Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_handler",
		New: func() caddy.Module { return new(Handler) },
	}
}

func (rb Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	_, _ = w.Write([]byte("Hello World!!!"))

	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*Handler)(nil)
