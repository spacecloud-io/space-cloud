package admin

import (
	"context"
	"sync"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Manager manages all admin transactions
type Manager struct {
	lock         sync.RWMutex
	user         *config.AdminUser
	integrations config.Integrations

	services model.ScServices
	isProd   bool

	syncMan        model.SyncManAdminInterface
	integrationMan IntegrationInterface

	nodeID, clusterID string
}

// New creates a new admin manager instance
func New(nodeID, clusterID string, isDev bool, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.nodeID = nodeID
	m.isProd = !isDev // set inverted
	m.clusterID = clusterID

	m.integrations = make(config.Integrations)
	m.user = adminUserInfo
	return m
}

// SetSyncMan sets syncman manager
func (m *Manager) SetSyncMan(s model.SyncManAdminInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.syncMan = s
}

// SetIntegrationMan sets integration manager
func (m *Manager) SetIntegrationMan(i IntegrationInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.integrationMan = i
}

// SetIntegrationConfig sets integration config
func (m *Manager) SetIntegrationConfig(integrations config.Integrations) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.integrations = integrations
}

// LoadEnv gets the env
func (m *Manager) LoadEnv() (bool, string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loginURL := "/mission-control/login"

	// Invoke integration hooks
	hookResponse := m.integrationMan.InvokeHook(ctx, model.RequestParams{
		Resource: "load-env",
		Op:       "read",
	})
	if hookResponse.CheckResponse() {
		if err := hookResponse.Error(); err == nil {
			loginURL = hookResponse.Result().(map[string]interface{})["loginUrl"].(string)
		}
	}
	return m.isProd, loginURL, nil
}
