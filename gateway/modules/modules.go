package modules

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Modules is an object that sets up the modules
type Modules struct {
	lock   sync.RWMutex
	blocks map[string]*Module

	nodeID             string
	removeProjectScope bool
	syncMan            *syncman.Manager
	adminMan           *admin.Manager
	metrics            *metrics.Module
	driver             *driver.Handler

	letsencrypt *letsencrypt.LetsEncrypt
	routing     *routing.Routing
}

// New creates a new modules instance
func New(nodeID string, removeProjectScope bool, syncMan *syncman.Manager, adminMan *admin.Manager, metrics *metrics.Module) (*Modules, error) {
	return &Modules{
		blocks:             map[string]*Module{},
		nodeID:             nodeID,
		removeProjectScope: removeProjectScope,
		syncMan:            syncMan,
		adminMan:           adminMan,
		metrics:            metrics,
		driver:             driver.New(removeProjectScope),
	}, nil
}

// SetGlobalModules sets the global modules
func (m *Modules) SetGlobalModules(letsencrypt *letsencrypt.LetsEncrypt, routing *routing.Routing) {
	m.letsencrypt = letsencrypt
	m.routing = routing
}

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(config *config.Project, le *letsencrypt.LetsEncrypt, ingressRouting *routing.Routing) error {
	module, err := m.loadModule(config.ID)
	if err != nil {
		module, err = m.newModule(config.ID)
		if err != nil {
			return err
		}
	}
	module.SetProjectConfig(config, le, ingressRouting)
	return nil
}

// SetGlobalConfig sets the auth secret and AESKey
func (m *Modules) SetGlobalConfig(projectID string, secret []*config.Secret, aesKey string) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetGlobalConfig(projectID, secret, aesKey)
}

// SetCrudConfig sets the config of db, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetCrudConfig(projectID, crudConfig)
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetServicesConfig(projectID, services)
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID string, fileStore *config.FileStore) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetFileStoreConfig(projectID, fileStore)
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetEventingConfig(projectID, eventingConfig)
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	module.SetUsermanConfig(projectID, auth)
	return nil
}

func (m *Modules) ProjectIDs() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ids := make([]string, 0)
	for id := range m.blocks {
		ids = append(ids, id)
	}
	return ids
}

func (m *Modules) Delete(projectID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.blocks, projectID)

	// Remove config from global modules
	_ = m.letsencrypt.DeleteProjectDomains(projectID)
	m.routing.DeleteProjectRoutes(projectID)
}

func (m *Modules) loadModule(projectID string) (*Module, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if module, p := m.blocks[projectID]; p {
		return module, nil
	}

	return nil, fmt.Errorf("project (%s) not found in server state", projectID)
}

func (m *Modules) newModule(projectID string) (*Module, error) {
	projectsIDs := m.ProjectIDs()
	m.lock.Lock()
	defer m.lock.Unlock()

	if ok := m.adminMan.ValidateProjectSyncOperation(projectsIDs, projectID); !ok {
		logrus.Println("Cannot create new project. Upgrade your plan")
		return nil, errors.New("upgrade your plan to create new project")
	}

	module := newModule(m.nodeID, m.removeProjectScope, m.syncMan, m.adminMan, m.metrics, m.driver)
	m.blocks[projectID] = module
	return module, nil
}
