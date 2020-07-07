package modules

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/modules/global"
)

// Modules is an object that sets up the modules
type Modules struct {
	lock   sync.RWMutex
	blocks map[string]*Module

	nodeID string

	// Global Modules
	GlobalMods *global.Global

	// Managers
	Managers *managers.Managers
}

// New creates a new modules instance
func New(nodeID string, managers *managers.Managers, globalMods *global.Global) (*Modules, error) {
	return &Modules{
		blocks:     map[string]*Module{},
		nodeID:     nodeID,
		GlobalMods: globalMods,
		Managers:   managers,
	}, nil
}

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(config *config.Project) error {
	module, err := m.loadModule(config.ID)
	if err != nil {
		module, err = m.newModule(config.ID)
		if err != nil {
			return err
		}
	}
	_ = module.SetProjectConfig(config)
	return nil
}

// SetGlobalConfig sets the auth secret and AESKey
func (m *Modules) SetGlobalConfig(projectID, secretSource string, secret []*config.Secret, aesKey string) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetGlobalConfig(projectID, secretSource, secret, aesKey)
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
	_ = module.SetUsermanConfig(projectID, auth)
	return nil
}

func (m *Modules) projects() *config.Config {
	m.lock.RLock()
	defer m.lock.RUnlock()

	c := &config.Config{Projects: []*config.Project{}}
	for id := range m.blocks {
		c.Projects = append(c.Projects, &config.Project{ID: id})
	}
	return c
}

func (m *Modules) Delete(projectID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.blocks, projectID)

	// Remove config from global modules
	_ = m.LetsEncrypt().DeleteProjectDomains(projectID)
	m.Routing().DeleteProjectRoutes(projectID)
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
	projects := m.projects()
	m.lock.Lock()
	defer m.lock.Unlock()

	if ok := m.Managers.Admin().ValidateProjectSyncOperation(projects, &config.Project{ID: projectID}); !ok {
		logrus.Println("Cannot create new project. Upgrade your plan")
		return nil, errors.New("upgrade your plan to create new project")
	}

	module := newModule(m.nodeID, m.Managers, m.GlobalMods)
	m.blocks[projectID] = module
	return module, nil
}
