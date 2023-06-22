package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spacecloud-io/space-cloud/utils"
)

type Refresh struct {
	Secret   string `json:"secret"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Refresh) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_admin_refresh_handler",
		New: func() caddy.Module { return new(Refresh) },
	}
}

func (h *Refresh) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	tokenString, ok := getTokenFromHeader(r)
	if !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("token not found in header"))
		return nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.Secret), nil
	})
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err)
		return nil
	}

	// Check if the token is valid
	if !token.Valid {
		utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid token"))
		return nil
	}

	// Get the claims from the existing token
	claims := token.Claims.(jwt.MapClaims)

	// Update the expiration time
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	// Create a new token with the updated claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the new token with the secret key
	tokenString, err = newToken.SignedString([]byte(h.Secret))
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err)
		return nil
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

var _ caddyhttp.MiddlewareHandler = (*Refresh)(nil)
