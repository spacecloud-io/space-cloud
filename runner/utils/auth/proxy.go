package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// VerifyProxyToken is used for authenticating websocket requests from metrics proxy
func (m *Module) VerifyProxyToken(token string) (map[string]interface{}, error) {
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}

		// Return the key
		return []byte(m.config.ProxySecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		tokenClaims := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			tokenClaims[key] = val
		}

		return tokenClaims, nil
	}

	return nil, errors.New("token could not be verified")
}
