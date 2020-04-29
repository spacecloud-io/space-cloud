package crud

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/graph-gophers/dataloader"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/bolt"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"

	"github.com/spaceuptech/space-cloud/gateway/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/sql"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	project            string
	removeProjectScope bool
	schema             model.SchemaCrudInterface

	// batch operation
	batchMapTableToChan batchMap // every table gets mapped to group of channels

	dataLoader loader
	// Variables to store the hooks
	hooks      *model.CrudHooks
	metricHook model.MetricCrudHook

	// Extra variables for enterprise
	blocks map[string]Crud
	admin  *admin.Manager
}

type loader struct {
	loaderMap      map[string]*dataloader.Loader
	dataLoaderLock sync.Mutex
}

// Crud abstracts the implementation crud operations of databases
type Crud interface {
	Create(ctx context.Context, project, col string, req *model.CreateRequest) (int64, error)
	Read(ctx context.Context, project, col string, req *model.ReadRequest) (int64, interface{}, error)
	Update(ctx context.Context, project, col string, req *model.UpdateRequest) (int64, error)
	Delete(ctx context.Context, project, col string, req *model.DeleteRequest) (int64, error)
	Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error)
	Batch(ctx context.Context, project string, req *model.BatchRequest) ([]int64, error)
	DescribeTable(ctc context.Context, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error)
	RawExec(ctx context.Context, project string) error
	GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error)
	DeleteCollection(ctx context.Context, project, col string) error
	CreateDatabaseIfNotExist(ctx context.Context, project string) error
	RawBatch(ctx context.Context, batchedQueries []string) error
	GetDBType() utils.DBType
	IsClientSafe() error
	Close() error
	GetConnectionState(ctx context.Context) bool
}

// Init create a new instance of the Module object
func Init(removeProjectScope bool) *Module {
	return &Module{removeProjectScope: removeProjectScope, batchMapTableToChan: make(batchMap), dataLoader: loader{loaderMap: map[string]*dataloader.Loader{}}}
}

// SetSchema sets the schema module
func (m *Module) SetSchema(s model.SchemaCrudInterface) {
	m.schema = s
}

// SetAdminManager sets the admin manager
func (m *Module) SetAdminManager(a *admin.Manager) {
	m.admin = a
}

// SetHooks sets the internal hooks
func (m *Module) SetHooks(hooks *model.CrudHooks, metricHook model.MetricCrudHook) {
	m.hooks = hooks
	m.metricHook = metricHook
}

func (m *Module) initBlock(dbType utils.DBType, enabled bool, connection string) (Crud, error) {
	switch dbType {
	case utils.Mongo:
		return mgo.Init(enabled, connection)
	case utils.EmbeddedDB:
		return bolt.Init(enabled, connection)
	case utils.MySQL, utils.Postgres, utils.SQLServer:
		return sql.Init(dbType, enabled, m.removeProjectScope, connection)
	default:
		return nil, utils.ErrInvalidParams
	}
}

func (m *Module) getCrudBlock(dbAlias string) (Crud, error) {
	block, p := m.blocks[dbAlias]
	if !p {
		return nil, fmt.Errorf("crud module not initialized yet for %s", dbAlias)
	}
	return block, nil
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crud config.Crud) error {
	m.Lock()
	defer m.Unlock()

	if err := m.admin.IsDBConfigValid(crud); err != nil {
		return err
	}

	m.project = project

	// Close the previous database connection
	for _, block := range m.blocks {
		utils.CloseTheCloser(block)
	}

	// Reset the blocks
	m.blocks = map[string]Crud{}

	// clear previous data loader
	m.dataLoader = loader{loaderMap: map[string]*dataloader.Loader{}}

	// Create a new crud blocks
	for k, v := range crud {
		var c Crud
		var err error
		if v.Type == "" {
			v.Type = k
		}

		v.Type = strings.TrimPrefix(v.Type, "sql-")
		c, err = m.initBlock(utils.DBType(v.Type), v.Enabled, v.Conn)

		if v.Enabled {
			if err != nil {
				logrus.Errorf("Error connecting to " + k + " : " + err.Error())
				return err
			}
			logrus.Info("Successfully connected to " + k)
		}

		// Store the block
		m.blocks[strings.TrimPrefix(k, "sql-")] = c
	}

	m.closeBatchOperation()
	m.initBatchOperation(project, crud)
	return nil
}

// GetDBType returns the type of the db for the alias provided
func (m *Module) GetDBType(dbAlias string) (string, error) {
	dbAlias = strings.TrimPrefix(dbAlias, "sql-")
	block, p := m.blocks[dbAlias]
	if !p {
		return "", fmt.Errorf("crud module not initialized yet for %s", dbAlias)
	}

	return string(block.GetDBType()), nil
}
