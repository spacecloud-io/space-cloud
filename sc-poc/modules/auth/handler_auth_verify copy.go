package auth

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/utils"
)

// AuthKratosVerifyHandler is responsible to authenticate the incoming request
// using Kratos
type AuthKratosVerifyHandler struct {
	logger *zap.Logger
	//authApp *App
}

// CaddyModule returns the Caddy module information.
func (AuthKratosVerifyHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_auth__kratos_verify_handler",
		New: func() caddy.Module { return new(AuthKratosVerifyHandler) },
	}
}

// Provision sets up the auth verify module.
func (h *AuthKratosVerifyHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Get the auth app
	// app, err := ctx.App("auth")
	// if err != nil {
	// 	h.logger.Error("Unable to load the auth provider", zap.Error(err))
	// 	return err
	// }

	// h.authApp = app.(*App)
	return nil
}

// ServeHTTP handles the http request
func (h *AuthKratosVerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Prepare authentication object
	result := types.AuthResult{}
	// Check if token is present in the header
	token, p := getTokenFromHeader(r)
	if p {
		// claims, err := h.authApp.Verify(token)
		// if err != nil {
		// 	SendUnauthenticatedResponse(w, r, h.logger, err)
		// 	return nil
		// }

		// result.IsAuthenticated = true
		// result.Claims = claims
	}

	// Add the result in the context object
	r = utils.StoreAuthenticationResult(r, &result)
	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddy.Provisioner = (*AuthVerifyHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*AuthVerifyHandler)(nil)
