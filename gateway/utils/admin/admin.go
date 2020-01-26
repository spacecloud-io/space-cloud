package admin

import (
	"errors"
	"net/http"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Manager manages all admin transactions
type Manager struct {
	lock   sync.RWMutex
	admin  *config.Admin
	isProd bool
}

// New creates a new admin manager instance
func New() *Manager {
	return &Manager{}
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(admin *config.Admin) {
	m.lock.Lock()
	m.admin = admin
	m.lock.Unlock()
}

// SetEnv sets the env
func (m *Manager) SetEnv(isProd bool) {
	m.lock.Lock()
	m.isProd = isProd
	m.lock.Unlock()
}

// LoadEnv gets the env
func (m *Manager) LoadEnv() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.isProd
}

// Login handles the admin login operation
func (m *Manager) Login(user, pass string) (int, string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, u := range m.admin.Users {
		if u.User == user && u.Pass == pass {
			token, err := m.createToken(map[string]interface{}{"id": user, "role": user})
			if err != nil {
				return http.StatusInternalServerError, "", err
			}
			return http.StatusOK, token, nil
		}
	}

	return http.StatusUnauthorized, "", errors.New("invalid credentials provided")
}
// RefreshToken is used to create a new token!
func (m *Manager) RefreshToken(token string) (int,string,error){
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Parse the token to get userID and userRole
	tokenClaims ,err := m.parseToken(token)
	if err != nil {
		return http.StatusUnauthorized,"", err
	}
	// Create a new token
		newToken, err := m.createToken(map[string]interface{}{"id": tokenClaims["id"], "nodeID": tokenClaims["nodeID"]})
		if err != nil {
			return http.StatusInternalServerError, "", errors.New("Failed to create a JWT token")
		}
	return http.StatusOK,newToken,nil

}