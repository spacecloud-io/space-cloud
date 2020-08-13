package admin

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GenerateToken creates a token with the appropriate claims
func (m *Manager) GenerateToken(token string, tokenClaims map[string]interface{}) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	claims, err := m.parseToken(token)
	if err != nil {
		return "", err
	}

	res := m.integrationMan.HandleConfigAuth("admin-token", "create", claims, nil)
	if res.CheckResponse() {
		if err := res.Error(); err != nil {
			return "", err
		}

		return m.createToken(tokenClaims)
	}

	return "", errors.New("only integrations are allowed to generate token")
}

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

func (m *Manager) parseToken(token string) (map[string]interface{}, error) {
	// Parse the JWT token
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
		if err := claims.Valid(); err != nil {
			return nil, err
		}

		obj := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			obj[key] = val
		}

		return obj, nil
	}

	return nil, utils.LogError("Unable to verify JWT token", "admin", "parse-token", nil)
}
