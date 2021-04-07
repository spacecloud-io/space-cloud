package admin

import (
	"context"
	"sync"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Manager manages all admin transactions
type Manager struct {
	lock   sync.RWMutex
	config *config.Admin
	user   *config.AdminUser
	isProd bool

	syncMan        model.SyncManAdminInterface
	integrationMan IntegrationInterface

	nodeID, clusterID      string
}

// New creates a new admin manager instance
func New(nodeID, clusterID string, isDev bool, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.nodeID = nodeID
	m.isProd = !isDev // set inverted
	m.clusterID = clusterID

	// Initialise all config
	m.config = new(config.Admin)
	m.user = adminUserInfo

	return m
}

func (m *Manager) SetSyncMan(s model.SyncManAdminInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.syncMan = s
}

func (m *Manager) SetIntegrationMan(i IntegrationInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.integrationMan = i
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(config *config.Admin, isInitialCall bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Set the admin config
	m.config = config

	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), "Gateway running in open source mode", nil)
	return nil
}

// GetConfig returns the admin config
func (m *Manager) GetConfig() *config.Admin {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.config
}

// LoadEnv gets the env
func (m *Manager) LoadEnv() (bool, string) {
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

	return m.isProd, loginURL
}
