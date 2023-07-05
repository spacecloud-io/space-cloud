package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

type RefreshHandler struct {
	logger   *zap.Logger
	adminMan *App
}

func (RefreshHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_admin_refresh_handler",
		New: func() caddy.Module { return new(RefreshHandler) },
	}
}

func (h *RefreshHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	adminManT, err := ctx.App("sc_admin")
	if err != nil {
		h.logger.Error("Unable to load the admin manager", zap.Error(err))
	}
	h.adminMan = adminManT.(*App)
	return nil
}

func (h *RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	tokenString, ok := getTokenFromHeader(r)
	if !ok {
		return utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("token not found in header"))
	}

	err := h.adminMan.VerifySCToken(tokenString)
	if err != nil {
		return utils.SendErrorResponse(w, http.StatusInternalServerError, err)
	}

	tokenString, err = h.adminMan.SignSCToken()
	if err != nil {
		return utils.SendErrorResponse(w, http.StatusUnauthorized, err)
	}

	utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": tokenString})
	return nil
}

func getTokenFromHeader(r *http.Request) (string, bool) {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	arr := strings.Split(tokens[0], " ")
	if strings.ToLower(arr[0]) == "bearer" && len(arr) >= 2 {
		return arr[1], true
	}

	return "", false
}

var _ caddyhttp.MiddlewareHandler = (*RefreshHandler)(nil)
