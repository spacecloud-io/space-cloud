package admin

import (
	"context"
	"net/http"

	"github.com/spaceuptech/helpers"
)

// Login handles the admin login operation
func (m *Manager) Login(ctx context.Context, user, pass string) (int, string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.user.User == user && m.user.Pass == pass {
		token, err := m.createToken(map[string]interface{}{"id": user, "role": "admin"})
		if err != nil {
			return http.StatusInternalServerError, "", err
		}
		return http.StatusOK, token, nil
	}

	return http.StatusUnauthorized, "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid username or password provided", nil, map[string]interface{}{"user": user, "pass": pass})
}
