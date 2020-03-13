package crud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/graph-gophers/dataloader"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/bolt"
	"github.com/spaceuptech/space-cloud/gateway/utils"

	"github.com/spaceuptech/space-cloud/gateway/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/sql"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	block              Crud
	dbType             string
	alias              string
	project            string
	removeProjectScope bool
	schema             model.SchemaCrudInterface
	// batch operation
	batchMapTableToChan batchMap // every table gets mapped to group of channels

	dataLoader loader
	// Variables to store the hooks
	hooks      *model.CrudHooks
	metricHook model.MetricCrudHook
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

func (m *Module) getCrudBlock(dbType string) (Crud, error) {
	if m.block != nil {
		return m.block, nil
	}
	return nil, fmt.Errorf("crud module not initialized yet for %q", dbType)
}

// SetServiceRoutes set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crud config.Crud) error {
	m.Lock()
	defer m.Unlock()
	m.closeBatchOperation()

	if len(crud) > 1 {
		return errors.New("crud module cannot have more than 1 db")
	}

	m.project = project

	// Close the previous database connection
	if m.block != nil {
		utils.CloseTheCloser(m.block)
	}

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

		m.dbType = v.Type
		m.block = c
		m.alias = strings.TrimPrefix(k, "sql-")
	}
	m.initBatchOperation(project, crud)
	return nil
}

// GetDBType returns the type of the db for the alias provided
func (m *Module) GetDBType(dbAlias string) (string, error) {
	dbAlias = strings.TrimPrefix(dbAlias, "sql-")
	if dbAlias != m.alias {
		return "", fmt.Errorf("db (%s) not found", dbAlias)
	}
	return m.dbType, nil
}
