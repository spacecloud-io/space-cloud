package crud

import (
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/crud/driver"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	blocks    map[string]*stub
	primaryDB string
	h         *driver.Handler
}

// Init create a new instance of the Module object
func Init(h *driver.Handler) *Module {
	return &Module{blocks: make(map[string]*stub), h: h}
}

// GetPrimaryDB get the database configured as primary
func (m *Module) GetPrimaryDB() (driver.Crud, error) {
	m.RLock()
	defer m.RUnlock()

	c, p := m.blocks[m.primaryDB]
	if !p {
		return nil, errors.New("CRUD: Primary DB not configured")
	}

	return c.c, nil
}

// SetConfig set the rules adn secret key required by the crud block
func (m *Module) SetConfig(crud config.Crud) error {
	m.Lock()
	defer m.Unlock()

	// Close the previous database connections
	for _, v := range m.blocks {
		m.h.RemoveBlock(v.dbType, v.conn)
	}

	m.blocks = make(map[string]*stub, len(crud))

	// Create a new crud blocks
	for dbType, v := range crud {
		// Skip this block if it is not enabled
		if !v.Enabled {
			continue
		}

		// Initialise a new block
		c, err := m.h.InitBlock(utils.DBType(dbType), v.Conn)
		if err != nil {
			return errors.New("CURD: Error - " + dbType + " could not be initialised - " + err.Error())
		}

		if v.IsPrimary {
			m.primaryDB = dbType
		}
		m.blocks[dbType] = &stub{c: c, conn: v.Conn, dbType: utils.DBType(dbType)}
	}

	return nil
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

	return nil, errors.New("CRUD: No crud block present for db")
}
