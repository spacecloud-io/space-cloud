package admin

import (
	"errors"
	"net/http"
)

// IsTokenValid checks if the token is valid
func (m *Manager) IsTokenValid(token string) error {
	_, err := m.parseToken(token)
	return err
}

// IsAdminOpAuthorised checks if the admin operation is authorised.
// TODO add scope level restrictions as well
func (m *Manager) IsAdminOpAuthorised(token, scope string) (int, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	auth, err := m.parseToken(token)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	user, p := auth["id"]
	if !p {
		return http.StatusUnauthorized, errors.New("Invalid Token")
	}

	for _, u := range m.admin.Users {
		if u.User == user {

			// Allow full access for scope name `all`
			if _, p := u.Scopes["all"]; p {
				return http.StatusOK, nil
			}

			// Check if scope is present
			if _, p := u.Scopes[scope]; p {
				return http.StatusOK, nil
			}

			break
		}
	}

	return http.StatusForbidden, errors.New("You are not authorized to make this request")
}
