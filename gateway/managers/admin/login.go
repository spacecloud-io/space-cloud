package admin

import (
	"context"
	"errors"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Login handles the admin login operation
func (m *Manager) Login(ctx context.Context, user, pass string) (int, string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Prepare the request params object
	params := model.RequestParams{
		Payload:  map[string]string{"user": user, "key": pass},
		Method:   http.MethodPost,
		Resource: "admin-login",
		Op:       "access",
	}

	// Invoke integration hooks
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), "", err
		}

		// Check the status code first. Send error for non 200 status code
		res := hookResponse.Result().(map[string]interface{})

		// Return the token
		return http.StatusOK, res["token"].(string), nil
	}

	if m.user.User == user && m.user.Pass == pass {
		token, err := m.createToken(map[string]interface{}{"id": user, "role": user})
		if err != nil {
			return http.StatusInternalServerError, "", err
		}
		return http.StatusOK, token, nil
	}

	return http.StatusUnauthorized, "", errors.New("Invalid credentials provided")
}
