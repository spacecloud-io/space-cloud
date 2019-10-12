package admin

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

// Manager manages all admin transactions
type Manager struct {
	lock      sync.RWMutex
	nodeID    string
	admin     *config.Admin
	validator *validator
	isProd    bool
}

// New creates a new admin manager instance
func New(nodeID string) *Manager {
	m := new(Manager)
	m.nodeID = nodeID
	m.validator = newValidator(m.reduceOpMode)
	return m
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(admin *config.Admin) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if admin.Operation.Mode > 0 {
		// Start the validation process for higher op modes
		if err := m.validator.startValidation(m.nodeID, admin.Operation.UserID, admin.Operation.Key, admin.Operation.Mode); err != nil {
			log.Println("Could not start in enterprise mode:", err)
			log.Println("Switching to community edition")
			admin.Operation.Mode = 0
		}
	}

	m.admin = admin
}

// GetConfig returns the adming config
func (m *Manager) GetConfig() *config.Admin {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.admin
}

// SetOperationMode sets the operation mode
func (m *Manager) SetOperationMode(op *config.OperationConfig) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if op.Mode > 0 && (op.UserID == "" || op.Key == "") {
		return errors.New("Invalid operation setting provided")
	}

	if op.Mode > 0 {
		// Start the validation process for higher op modes
		if err := m.validator.startValidation(m.nodeID, op.UserID, op.Key, op.Mode); err != nil {
			return err
		}
	} else {
		// Stop validation for open source mode
		m.validator.stopValidation()
	}

	m.admin.Operation = *op
	return nil
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

	for _, u := range m.admin.Users {
		if u.User == user && u.Pass == pass {
			token, err := m.createToken(map[string]interface{}{"id": user, "role": user})
			if err != nil {
				return http.StatusInternalServerError, "", err
			}
			return http.StatusOK, token, nil
		}
	}

	return http.StatusUnauthorized, "", errors.New("invalid credentials provided")
}
