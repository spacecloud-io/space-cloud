package eventing

import (
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// Module is responsible for managing the eventing system
type Module struct {
	lock sync.RWMutex

	// Configurable variables
	project string
	config  *config.Eventing

	// Atomic maps to handle events being processed
	processingEvents sync.Map

	// Variables defined during initialisation
	auth      *auth.Module
	crud      *crud.Module
	functions *functions.Module
	adminMan  *admin.Manager
	syncMan   *syncman.Manager
}

// New creates a new instance of the eventing module
func New(auth *auth.Module, crud *crud.Module, functions *functions.Module, adminMan *admin.Manager, syncMan *syncman.Manager) *Module {

	m := &Module{
		auth:      auth,
		crud:      crud,
		functions: functions,
		adminMan:  adminMan,
		syncMan:   syncMan,
		config:    &config.Eventing{Enabled: false, InternalRules: map[string]config.EventingRule{}},
	}

	// Start the internal processes
	go m.routineProcessIntents()
	go m.routineProcessStaged()

	return m
}

// SetConfig sets the module config
func (m *Module) SetConfig(project string, eventing *config.Eventing) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if eventing == nil || !eventing.Enabled {
		m.config.Enabled = false
		return nil
	}

	if eventing.DBType == "" || eventing.Col == "" {
		return errors.New("invalid eventing config provided")
	}

	m.project = project
	m.config = eventing

	if m.config.Rules == nil {
		m.config.Rules = map[string]config.EventingRule{}
	}

	// Reset the internal rules
	m.config.InternalRules = map[string]config.EventingRule{}

	return nil
}
