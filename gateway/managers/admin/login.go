package admin

import (
	"context"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Login handles the admin login operation
func (m *Manager) Login(ctx context.Context, user, pass string) (int, string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.user.User == user && m.user.Pass == pass {
		token, err := m.createToken(map[string]interface{}{"id": user, "role": user})
		if err != nil {
			return http.StatusInternalServerError, "", err
		}
		return http.StatusOK, token, nil
	}

	return http.StatusUnauthorized, "", utils.LogError("Invalid credentials provided", "admin", "login", nil)
}
