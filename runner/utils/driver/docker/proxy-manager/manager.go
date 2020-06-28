package manager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// Manager is used for proxy management
type Manager struct {
	lock sync.RWMutex

	serviceRoutes model.Config
	servers       map[int32]*http.Server
	path          string
}

// New creates a new manager
func New(path string) (*Manager, error) {
	manager := &Manager{
		servers:       map[int32]*http.Server{},
		serviceRoutes: model.Config{},
		path:          path,
	}

	// Load the config from the file
	if err := manager.loadConfigFromFile(); err != nil {
		return nil, err
	}

	// Start all servers
	manager.adjustProxyServers()

	return manager, nil
}

// loadConfigFromFile loads the route config from file
func (m *Manager) loadConfigFromFile() error {
	content, err := ioutil.ReadFile(m.path)
	if err != nil {
		logrus.Info("Could not load service routing file. Resetting to default")
		m.serviceRoutes = model.Config{}
		return nil
	}

	return json.Unmarshal(content, &m.serviceRoutes)
}

// writeConfigFromFile writes the route config to file
func (m *Manager) writeConfigToFile() error {
	routeConfig, err := json.MarshalIndent(m.serviceRoutes, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.path, routeConfig, 0644)
}
