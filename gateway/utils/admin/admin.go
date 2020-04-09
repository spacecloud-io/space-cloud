package admin

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Manager manages all admin transactions
type Manager struct {
	lock      sync.RWMutex
	nodeID    string
	publicKey *rsa.PublicKey
	isProd    bool

	quotas model.UsageQuotas
	config *config.Admin

	closeFetchPublicRSAKey chan struct{}

	// config
	user *config.AdminUser
}

// New creates a new admin manager instance
func New(nodeID string, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.nodeID = nodeID

	// Initialise all config
	m.config = new(config.Admin)
	m.user = adminUserInfo
	m.quotas = model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1, Version: 0}

	// Initialise channel for closing fetch public key
	m.closeFetchPublicRSAKey = make(chan struct{}, 1)

	// Return the admin manager
	return m
}

func (m *Manager) startOperation() error {
	logrus.Infoln("Starting gateway in enterprise mode")

	// Stop the previous fetch public key go routine
	m.closeFetchPublicRSAKey <- struct{}{} // close previous go routine

	// Get the public key
	if err := m.fetchPublicKeyWithoutLock(); err != nil {
		logrus.Errorln(err)
		return fmt.Errorf("unable to fetch public key (%s)", err)
	}

	go m.fetchPublicKeyRoutine()

	// We'll fetch the quotas only if the config version is greater than the one we have stored.
	// m.quotas.Version can be updated by the fetchQuotas method
	if m.quotas.Version < m.config.Version {
		if err := m.fetchQuotas(); err != nil {
			return fmt.Errorf("unable to fetch quotas (%s)", err.Error())
		}
	}

	return nil
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(config *config.Admin) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Set the admin
	m.config = config

	if m.isEnterpriseMode() {
		return m.startOperation()
	}

	// Reset quotas and version to defaults
	m.quotas.MaxProjects = 1
	m.quotas.MaxDatabases = 1
	m.quotas.Version = 0
	m.config.Version = 0
	return nil
}

func (m *Manager) SetVersion(version int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Send error if not in enterprise mode
	if !m.isEnterpriseMode() {
		return errors.New("cannot set version when not in enterprise mode")
	}

	m.config.Version = version

	return m.startOperation()
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

	// Dont allow logins when in enterprise mode
	if m.isEnterpriseMode() {
		return http.StatusUnauthorized, "", errors.New("cannot login when in enterprise mode")
	}

	if m.user.User == user && m.user.Pass == pass {
		token, err := m.createToken(map[string]interface{}{"id": user, "role": user})
		if err != nil {
			return http.StatusInternalServerError, "", err
		}
		return http.StatusOK, token, nil
	}

	return http.StatusUnauthorized, "", errors.New("Invalid credentials provided")
}
