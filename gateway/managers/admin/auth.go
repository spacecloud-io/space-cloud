package admin

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spaceuptech/helpers"
)

func (m *Manager) createToken(tokenClaims map[string]interface{}) (string, error) {

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}

	// Add expiry of one week
	claims["exp"] = time.Now().Add(24 * 7 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.user.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) parseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid signing method type", nil, nil)
		}

		return []byte(m.user.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		if err := claims.Valid(); err != nil {
			return nil, err
		}

		obj := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			obj[key] = val
		}
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Claim from request token", map[string]interface{}{"claims": claims, "type": "admin"})
		return obj, nil
	}

	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Admin: JWT token could not be verified", nil, nil)
}
