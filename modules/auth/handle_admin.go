package auth

import (
	"errors"
	"net/http"
)

// AdminLogin handles the admin login operation
func (m *Module) AdminLogin(project, user, pass string) (int, string, error) {
	m.RLock()
	u, p, r := m.admin.User, m.admin.Pass, m.admin.Role
	proj := m.project
	m.RUnlock()

	if u != user || p != pass || proj != project {
		return http.StatusUnauthorized, "", errors.New("invalid credentials provided")
	}

	token, err := m.CreateToken(TokenClaims{"id": r, "role": r})
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}

// IsAdminOpAuthorised checks if the admin operation is authorised
func (m *Module) IsAdminOpAuthorised(project, token string) (int, error) {
	m.RLock()
	defer m.RUnlock()

	auth, err := m.parseToken(token)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	role, err := auth.GetRole()
	if err != nil {
		return http.StatusUnauthorized, err
	}

	if project != m.project {
		return http.StatusForbidden, errors.New("Invalid project")

	}

	if role != m.admin.Role {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	return http.StatusOK, nil
}
