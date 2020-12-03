package managers

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
)

// Managers holds all the managers
type Managers struct {
	adminMan *admin.Manager
	syncMan  *syncman.Manager
}

// New creates a new managers instance
func New(nodeID, clusterID, storeType, runnerAddr string, isDev bool, adminUserInfo *config.AdminUser, ssl *config.SSL) (*Managers, error) {
	// Create the fundamental modules
	adminMan := admin.New(nodeID, clusterID, isDev, adminUserInfo)
	syncMan, err := syncman.New(nodeID, clusterID, storeType, runnerAddr, adminMan, ssl)
	if err != nil {
		return nil, err
	}

	return &Managers{adminMan: adminMan, syncMan: syncMan}, nil
}

// Admin returns the admin manager
func (m *Managers) Admin() *admin.Manager {
	return m.adminMan
}

// Sync returns the sync manager
func (m *Managers) Sync() *syncman.Manager {
	return m.syncMan
}
