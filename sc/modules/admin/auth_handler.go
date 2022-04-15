package admin

import (
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/modules/middlewares"
	"github.com/spacecloud-io/space-cloud/utils"
)

// AuthHandler is a module to create admin auth handler
type AuthHandler struct {
	logger *zap.Logger
	admin  *App
}

// CaddyModule returns the Caddy module information.
func (AuthHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_admin_auth_handler",
		New: func() caddy.Module { return new(AuthHandler) },
	}
}

// Provision sets up the handler
func (h *AuthHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Get the admin auth
	app, err := ctx.App("admin")
	if err != nil {
		h.logger.Error("Unable to load admin module", zap.Error(err))
		return err
	}
	h.admin = app.(*App)
	return nil
}

// ServeHTTP handles the http request
func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get the admin token
	token := utils.GetTokenFromHeader(r)

	// Attemp to parse the token if present
	if token != "" {
		claims, err := h.admin.Auth().Verify(token)
		if err != nil {
			h.logger.Error("Unable to parse admin token", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, errors.New("invalid admin token provided"))
			return nil
		}

		// Get the request params from context and store the claims in them
		reqParams := middlewares.GetRequestParams(r)
		reqParams.Claims = claims
	}

	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddy.Provisioner = (*AuthHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*AuthHandler)(nil)
