package auth

import (
	"context"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// AuthHandler is responsible to authenticate the incoming request
// on a best effort basis
type AuthHandler struct {
	logger  *zap.Logger
	authApp *App
}

// CaddyModule returns the Caddy module information.
func (AuthHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_auth_handler",
		New: func() caddy.Module { return new(AuthHandler) },
	}
}

// Provision sets up the graphql module.
func (h *AuthHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Get the auth app
	app, _ := ctx.App("auth")
	h.authApp = app.(*App)
	return nil
}

// ServeHTTP handles the http request
func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Prepare authentication object
	result := AuthResult{}
	// Check if token is present in the header
	token, p := getTokenFromHeader(r)
	if p {
		claims, err := h.authApp.Verify(token)
		if err != nil {
			SendUnauthenticatedResponse(w, r, h.logger, err)
			return nil
		}

		result.IsAuthenticated = true
		result.Claims = claims
	}

	// Add the result in the context object
	r = r.WithContext(context.WithValue(r.Context(), authenticationResultKey, &result))
	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*AuthHandler)(nil)
