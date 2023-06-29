package auth

import (
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"github.com/spacecloud-io/space-cloud/utils"
)

type AuthenticateUser struct{}

func (AuthenticateUser) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_plugin_authenticate_user_handler",
		New: func() caddy.Module { return new(AuthenticateUser) },
	}
}

// ServeHTTP handles the http request
func (h *AuthenticateUser) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	authResult, p := utils.GetAuthenticationResult(r.Context())
	if !p {
		return utils.SendErrorResponse(w, http.StatusUnauthorized, errors.New("unable to authenticate request"))
	}
	if !authResult.IsAuthenticated {
		return utils.SendErrorResponse(w, http.StatusUnauthorized, errors.New("authentication failed"))
	}
	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*AuthenticateUser)(nil)
