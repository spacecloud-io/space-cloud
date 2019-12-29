package crud

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"

	"github.com/spaceuptech/space-cloud/modules/crud/driver"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	blocks  map[string]*stub
	project string

	// Variables to store the hooks
	hooks      *model.CrudHooks
	metricHook model.MetricCrudHook

	// Drivers handler
	h *driver.Handler

	// Admin manager
	adminMan *admin.Manager
}

// Init create a new instance of the Module object
func Init(h *driver.Handler, adminMan *admin.Manager) *Module {
	return &Module{blocks: make(map[string]*stub), h: h, adminMan: adminMan}
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crud config.Crud) error {
	m.Lock()
	defer m.Unlock()

	if err := m.adminMan.IsDBConfigValid(crud); err != nil {
		return err
	}

	m.project = project

	// Close the previous database connections
	for _, v := range m.blocks {
		m.h.RemoveBlock(v.dbType, v.conn)
	}

	m.blocks = make(map[string]*stub, len(crud))

	// Create a new crud blocks
	for dbAlias, v := range crud {
		// For backward compatibilty support
		if v.Type == "" {
			v.Type = dbAlias
		}
		v.Type = strings.TrimPrefix(v.Type, "sql-")

		// Initialise a new block
		c, err := m.h.InitBlock(utils.DBType(v.Type), v.Enabled, v.Conn)
		m.blocks[dbAlias] = &stub{c: c, conn: v.Conn, dbType: utils.DBType(v.Type)}

		if err != nil {
			log.Println("Error connecting to " + dbAlias + " : " + err.Error())
			return err
		}
		log.Println("Successfully connected to " + dbAlias)
	}
	return nil
}

// SetHooks sets the internal hooks
func (m *Module) SetHooks(hooks *model.CrudHooks, metricHook model.MetricCrudHook) {
	m.hooks = hooks
	m.metricHook = metricHook
}

type stub struct {
	conn   string
	c      driver.Crud
	dbType utils.DBType
}

func (m *Module) getCrudBlock(dbType string) (driver.Crud, error) {
	if crud, p := m.blocks[dbType]; p {
		return crud.c, nil
	}

	return nil, fmt.Errorf("database (%s) does not exist", dbType)
}

// GetDBType returns the type of the db for the alias provided
func (m *Module) GetDBType(dbAlias string) (string, error) {
	dbAlias = strings.TrimPrefix(dbAlias, "sql-")
	s, p := m.blocks[dbAlias]
	if !p {
		return "", fmt.Errorf("db (%s) not found", dbAlias)
	}
	return string(s.dbType), nil
}
