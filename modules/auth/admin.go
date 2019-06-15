package auth

import (
	"errors"
	"net/http"
)

// AdminLogin handles the admin login operation
func (m *Module) AdminLogin(user, pass string) (int, string, error) {
	m.RLock()
	u, p, r := m.admin.User, m.admin.Pass, m.admin.Role
	m.RUnlock()

	if u != user || p != pass {
		return http.StatusUnauthorized, "", errors.New("invalid credentials provided")
	}

	token, err := m.CreateToken(TokenClaims{"id": r, "role": r})
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}
