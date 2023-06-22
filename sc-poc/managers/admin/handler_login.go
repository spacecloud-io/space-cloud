package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spacecloud-io/space-cloud/utils"
)

type Login struct {
	Secret   string `json:"secret"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Login) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_admin_login_handler",
		New: func() caddy.Module { return new(Login) },
	}
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	cred := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		_ = utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid request payload received"))
		return nil
	}

	if cred.Username != h.Username || cred.Password != h.Password {
		_ = utils.SendErrorResponse(w, http.StatusUnauthorized, errors.New("invalid credentials provided"))
		return nil
	}

	jwtClaims := make(jwt.MapClaims)
	jwtClaims["sub"] = cred.Username
	jwtClaims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, errors.New(err.Error()))
		return nil
	}

	return utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": tokenString})
}

var _ caddyhttp.MiddlewareHandler = (*Login)(nil)
