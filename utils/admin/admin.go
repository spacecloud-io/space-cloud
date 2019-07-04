package admin

import (
	"errors"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/spaceuptech/space-cloud/config"
)

// Manager manages all admin transactions
type Manager struct {
	admin *config.Admin
}

// New creates a new admin manager instance
func New() *Manager {
	return &Manager{}
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(admin *config.Admin) {
	m.admin = admin
}

// Login handles the admin login operation
func (m *Manager) Login(user, pass string) (int, string, error) {
	log.Println("Admin", m.admin, m)
	u, p, r := m.admin.User, m.admin.Pass, m.admin.Role

	if u != user || p != pass {
		return http.StatusUnauthorized, "", errors.New("invalid credentials provided")
	}

	token, err := m.createToken(map[string]interface{}{"id": r, "role": r})
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}

// IsAdminOpAuthorised checks if the admin operation is authorised
func (m *Manager) IsAdminOpAuthorised(token string) (int, error) {

	auth, err := m.parseToken(token)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	role, p := auth["role"]
	if !p {
		return http.StatusUnauthorized, errors.New("Invalid Token")
	}

	if role != m.admin.Role {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	return http.StatusOK, nil
}

// CreateToken generates a new JWT Token with the token claims
func (m *Manager) createToken(tokenClaims map[string]interface{}) (string, error) {

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.admin.Secret))
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

		return []byte(m.admin.Secret), nil
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
