package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

type Login struct {
	logger   *zap.Logger
	adminMan *App
}

func (Login) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_admin_login_handler",
		New: func() caddy.Module { return new(Login) },
	}
}

func (h *Login) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	adminManT, err := ctx.App("sc_admin")
	if err != nil {
		h.logger.Error("Unable to load the admin manager", zap.Error(err))
	}
	h.adminMan = adminManT.(*App)
	return nil
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	creds := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		_ = utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid request payload received"))
		return nil
	}

	if err := h.adminMan.VerifyCredentials(creds); err != nil {
		return utils.SendErrorResponse(w, http.StatusUnauthorized, err)
	}

	tokenString, err := h.adminMan.SignSCToken()
	if err != nil {
		return utils.SendErrorResponse(w, http.StatusUnauthorized, err)
	}

	return utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": tokenString})
}

var (
	_ caddy.Provisioner           = (*Login)(nil)
	_ caddyhttp.MiddlewareHandler = (*Login)(nil)
)
