package crud

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/graph-gophers/dataloader"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"

	"github.com/spaceuptech/space-cloud/gateway/modules/crud/driver"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	blocks  map[string]*stub
	project string

	// batch operation
	batchMapTableToChan batchMap // every table gets mapped to group of channels

	dataLoader loader
	// Variables to store the hooks
	hooks      *model.CrudHooks
	metricHook model.MetricCrudHook

	schema model.SchemaCrudInterface

	// Drivers handler
	h *driver.Handler

	// Admin manager
	adminMan *admin.Manager
}

type loader struct {
	loaderMap      map[string]*dataloader.Loader
	dataLoaderLock sync.Mutex
}

// Init create a new instance of the Module object
func Init(h *driver.Handler, adminMan *admin.Manager) *Module {
	return &Module{blocks: make(map[string]*stub), h: h, adminMan: adminMan}
}

// SetSchema sets the schema module
func (m *Module) SetSchema(s model.SchemaCrudInterface) {
	m.schema = s
}

// SetHooks sets the internal hooks
func (m *Module) SetHooks(hooks *model.CrudHooks, metricHook model.MetricCrudHook) {
	m.hooks = hooks
	m.metricHook = metricHook
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crud config.Crud) error {
	m.Lock()
	defer m.Unlock()
	m.closeBatchOperation()

	if err := m.adminMan.IsDBConfigValid(crud); err != nil {
		return err
	}

	m.project = project

	// Close the previous database connections
	for _, v := range m.blocks {
		m.h.RemoveBlock(v.dbType, v.conn)
	}
	m.blocks = make(map[string]*stub, len(crud))

	// clear previous data loader
	m.dataLoader = loader{loaderMap: map[string]*dataloader.Loader{}}

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
	m.initBatchOperation(project, crud)
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
