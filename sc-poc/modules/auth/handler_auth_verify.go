package auth

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/provider"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/utils"
)

// JWTAuthVerifyHandler is responsible to authenticate the incoming request
// on a best effort basis
type JWTAuthVerifyHandler struct {
	logger      *zap.Logger
	authApp     *App
	providerMan *provider.App
}

// CaddyModule returns the Caddy module information.
func (JWTAuthVerifyHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_jwt_auth_verify_handler",
		New: func() caddy.Module { return new(JWTAuthVerifyHandler) },
	}
}

// Provision sets up the auth verify module.
func (h *JWTAuthVerifyHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Get the provider manager
	providerManT, err := ctx.App("provider")
	if err != nil {
		return err
	}
	h.providerMan = providerManT.(*provider.App)
	return nil
}

// ServeHTTP handles the http request
func (h *JWTAuthVerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get the auth app
	ws := source.GetWorkspaceNameFromHeaders(r)
	app, err := h.providerMan.GetProvider(ws, "auth")
	if err != nil {
		h.logger.Error("Unable to load auth module", zap.String("workspace", ws), zap.Error(err))
		return err
	}
	authApp := app.(*App)

	// Prepare authentication object
	result := types.AuthResult{}
	// Check if token is present in the header
	token, p := getTokenFromHeader(r)
	if p {
		claims, err := authApp.Verify(token)
		if err != nil {
			SendUnauthenticatedResponse(w, r, h.logger, err)
			return nil
		}

		result.IsAuthenticated = true
		result.Claims = claims
	}

	// Add the result in the context object
	r = utils.StoreAuthenticationResult(r, &result)
	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddy.Provisioner = (*JWTAuthVerifyHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*JWTAuthVerifyHandler)(nil)
