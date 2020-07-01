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

	nodeID, clusterID string
}

// New creates a new admin manager instance
func New(nodeID, clusterID string, isDev bool, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.config = new(config.Admin)
	m.user = adminUserInfo
	m.quotas = model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1}
	m.nodeID = nodeID
	m.clusterID = clusterID
	m.isProd = !isDev
	return m
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(admin *config.Admin) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.config = admin
	return nil
}

// GetConfig returns the admin config
func (m *Manager) GetConfig() *config.Admin {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.config
}

// LoadEnv gets the env
func (m *Manager) LoadEnv() (bool, string, model.UsageQuotas) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.isProd, "space-cloud-open", m.quotas
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
