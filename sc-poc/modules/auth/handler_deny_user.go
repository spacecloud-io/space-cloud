package auth

import (
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"github.com/spacecloud-io/space-cloud/utils"
)

type DenyUser struct{}

func (DenyUser) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_plugin_deny_user_handler",
		New: func() caddy.Module { return new(DenyUser) },
	}
}

// ServeHTTP handles the http request
func (h *DenyUser) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return utils.SendErrorResponse(w, http.StatusForbidden, errors.New("access denied"))
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*DenyUser)(nil)
