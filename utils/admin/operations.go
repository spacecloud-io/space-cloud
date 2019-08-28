package admin

import (
	"errors"
	"net/http"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// GetInternalAccessToken returns the token that can be used internally by Space Cloud
func (m *Manager) GetInternalAccessToken() (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.createToken(map[string]interface{}{"id": utils.InternalUserID})
}

// IsTokenValid checks if the token is valid
func (m *Manager) IsTokenValid(token string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return nil
	}

	_, err := m.parseToken(token)
	return err
}

// ValidateSyncOperation validates if an operation is permitted based on the mode
func (m *Manager) ValidateSyncOperation(c *config.Config, project *config.Project) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, p := range c.Projects {
		if p.ID == project.ID {
			return true
		}
	}

	maxProjects := 1
	if m.admin.Operation.Mode == 1 {
		maxProjects = 3
	} else if m.admin.Operation.Mode == 2 {
		maxProjects = 5
	}

	if len(c.Projects) == (maxProjects - 1) {
		return true
	}

	return false
}

// IsAdminOpAuthorised checks if the admin operation is authorised.
// TODO add scope level restrictions as well
func (m *Manager) IsAdminOpAuthorised(token, scope string) (int, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return http.StatusOK, nil
	}

	auth, err := m.parseToken(token)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	user, p := auth["id"]
	if !p {
		return http.StatusUnauthorized, errors.New("Invalid Token")
	}

	if user == utils.InternalUserID {
		return http.StatusOK, nil
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
