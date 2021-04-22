package managers

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/integration"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
)

// Managers holds all the managers
type Managers struct {
	adminMan       *admin.Manager
	syncMan        *syncman.Manager
	integrationMan *integration.Manager
}

// New creates a new managers instance
func New(nodeID, clusterID, storeType, runnerAddr string, isDev bool, adminUserInfo *config.AdminUser, ssl *config.SSL) (*Managers, error) {
	// Create the fundamental modules
	adminMan := admin.New(nodeID, clusterID, isDev, adminUserInfo)
	i := integration.New(adminMan)
	syncMan, err := syncman.New(nodeID, clusterID, storeType, runnerAddr, adminMan, i, ssl)
	if err != nil {
		return nil, err
	}
	adminMan.SetSyncMan(syncMan)
	adminMan.SetIntegrationMan(i)

	return &Managers{adminMan: adminMan, syncMan: syncMan, integrationMan: i}, nil
}

// Admin returns the admin manager
func (m *Managers) Admin() *admin.Manager {
	return m.adminMan
}

// Sync returns the sync manager
func (m *Managers) Sync() *syncman.Manager {
	return m.syncMan
}

// Integration returns the integration manager
func (m *Managers) Integration() *integration.Manager {
	return m.integrationMan
}
