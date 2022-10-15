package apis

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/rs/cors"
)

// CorsHandler is responsible to handler cors
type CorsHandler struct {
}

// CaddyModule returns the Caddy module information.
func (CorsHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_cors_handler",
		New: func() caddy.Module { return new(CorsHandler) },
	}
}

// ServeHTTP handles the http request
func (h *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc: func(s string) bool {
			return true
		},
		Debug: true,
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{},
	})

	c.HandlerFunc(w, r)
	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*RootAPIHandler)(nil)
