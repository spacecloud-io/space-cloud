package admin

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

func (m *Manager) createToken(tokenClaims map[string]interface{}) (string, error) {

	// Throw error if in enterprise mode
	if m.isEnterpriseMode() {
		return "", errors.New("cannot sign tokens in enterprise m")
	}

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.user.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) parseToken(token string) (map[string]interface{}, error) {
	// Parse the JWT token
	if m.isEnterpriseMode() {
		// enterprise mode, validate the token
		_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, errors.New("invalid signing method type")
			}
			return m.publicKey, nil
		})
		return make(map[string]interface{}), err
	}
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method type")
		}

		return []byte(m.user.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		obj := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			obj[key] = val
		}

		return obj, nil
	}

	return nil, errors.New("Admin: JWT token could not be verified")
}
