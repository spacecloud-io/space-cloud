package admin

import (
	"errors"
	"net/http"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Manager manages all admin transactions
type Manager struct {
	lock   sync.RWMutex
	config *config.Admin
	quotas model.UsageQuotas
	user   *config.AdminUser
	isProd bool

	clusterID string
}

// New creates a new admin manager instance
func New(clusterID string, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.config = new(config.Admin)
	m.user = adminUserInfo
	m.quotas = model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1}
	m.clusterID = clusterID
	return m
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(admin *config.Admin) {
	m.lock.Lock()
	m.config = admin
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

	if m.user.User == user && m.user.Pass == pass {
		token, err := m.createToken(map[string]interface{}{"id": user, "role": user})
		if err != nil {
			return http.StatusInternalServerError, "", err
		}
		return http.StatusOK, token, nil
	}

	return http.StatusUnauthorized, "", errors.New("Invalid credentials provided")
}
