package admin

import (
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

type AuthenticateSCUserPlugin struct {
	Name     string `json:"name"`
	logger   *zap.Logger
	adminMan *App
}

func (AuthenticateSCUserPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_plugin_authenticate_sc_user_handler",
		New: func() caddy.Module { return new(AuthenticateSCUserPlugin) },
	}
}

func (h *AuthenticateSCUserPlugin) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	adminManT, err := ctx.App("sc_admin")
	if err != nil {
		h.logger.Error("Unable to load the admin manager", zap.Error(err))
	}
	h.adminMan = adminManT.(*App)
	return nil
}

func (h *AuthenticateSCUserPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	tokenString, ok := getTokenFromHeader(r)
	if !ok {
		return utils.SendErrorResponse(w, http.StatusForbidden, errors.New("token not found in header"))
	}

	err := h.adminMan.VerifySCToken(tokenString)
	if err != nil {
		return utils.SendErrorResponse(w, http.StatusInternalServerError, err)
	}

	return next.ServeHTTP(w, r)
}

var _ caddyhttp.MiddlewareHandler = (*AuthenticateSCUserPlugin)(nil)
