package auth

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

// VerifyToken checks if the token is valid and returns the token claims
func (m *Module) VerifyToken(token string) (map[string]interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.config.IsDev {
		return nil, nil
	}
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect it to be
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}

		// Return the key
		return []byte(m.config.Secret), nil
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
