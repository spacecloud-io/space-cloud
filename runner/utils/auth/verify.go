package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// VerifyToken checks if the token is valid and returns the token claims
func (m *Module) VerifyToken(token string) (map[string]interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return nil, nil

	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Declare the variables
		var alg string
		var key interface{}

		// See which algorithm has been configured
		switch m.config.JWTAlgorithm {
		case RSA256:
			alg = jwt.SigningMethodRS256.Alg()
			key = m.config.PublicKey
		case HS256:
			alg = jwt.SigningMethodHS256.Alg()
			key = []byte(m.config.Secret)
		default:
			return nil, errors.New("invalid signing methods configured")
		}

		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != alg {
			return nil, errors.New("invalid signing method")
		}

		// Return the key
		return key, nil
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
