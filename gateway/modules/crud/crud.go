package crud

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/spaceuptech/helpers"

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
	block   Crud
	dbType  string
	alias   string
	project string
	schema  model.SchemaCrudInterface
	queries config.DatabasePreparedQueries
	// batch operation
	batchMapTableToChan batchMap // every table gets mapped to group of channels

	config *config.DatabaseConfig

	dataLoader loader
	// Variables to store the hooks
	hooks      *model.CrudHooks
	metricHook model.MetricCrudHook

	// function to get secrets from runner
	getSecrets utils.GetSecrets
}

type loader struct {
	loaderMap      map[string]*dataloader.Loader
	dataLoaderLock sync.RWMutex
}

// Crud abstracts the implementation crud operations of databases
type Crud interface {
	Create(ctx context.Context, col string, req *model.CreateRequest) (int64, error)
	Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, error)
	Update(ctx context.Context, col string, req *model.UpdateRequest) (int64, error)
	Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error)
	Aggregate(ctx context.Context, col string, req *model.AggregateRequest) (interface{}, error)
	Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error)
	DescribeTable(ctc context.Context, col string) ([]model.InspectorFieldType, []model.ForeignKeysType, []model.IndexType, error)
	RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error)
	GetCollections(ctx context.Context) ([]utils.DatabaseCollections, error)
	DeleteCollection(ctx context.Context, col string) error
	CreateDatabaseIfNotExist(ctx context.Context, name string) error
	RawBatch(ctx context.Context, batchedQueries []string) error
	GetDBType() model.DBType
	IsClientSafe(ctx context.Context) error
	IsSame(conn, dbName string, driverConf config.DriverConfig) bool
	Close() error
	GetConnectionState(ctx context.Context) bool
	SetQueryFetchLimit(limit int64)
}

// Init create a new instance of the Module object
func Init() *Module {
	return &Module{batchMapTableToChan: make(batchMap), dataLoader: loader{loaderMap: map[string]*dataloader.Loader{}}}
}

func (m *Module) initBlock(dbType model.DBType, enabled bool, connection, dbName string, driverConf config.DriverConfig) (Crud, error) {
	switch dbType {
	case model.Mongo:
		return mgo.Init(enabled, connection, dbName, driverConf)
	case model.EmbeddedDB:
		return bolt.Init(enabled, connection, dbName)
	case model.MySQL, model.Postgres, model.SQLServer:
		c, err := sql.Init(dbType, enabled, connection, dbName, driverConf)
		if err == nil && enabled {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := c.CreateDatabaseIfNotExist(ctx, dbName); err != nil {
				return nil, err
			}
		}
		if dbType == model.MySQL {
			return sql.Init(dbType, enabled, fmt.Sprintf("%s%s", connection, dbName), dbName, driverConf)
		}
		return c, err
	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unsupported database (%s) provided", dbType), nil, map[string]interface{}{})
	}
}

// GetDBType returns the type of the db for the alias provided
func (m *Module) GetDBType(dbAlias string) (string, error) {
	dbAlias = strings.TrimPrefix(dbAlias, "sql-")
	if dbAlias != m.alias {
		return "", fmt.Errorf("cannot get db type as invalid db alias (%s) provided", dbAlias)
	}
	return m.dbType, nil
}

// CloseConfig close the rules and secret key required by the crud block
func (m *Module) CloseConfig() error {
	// Acquire a lock
	m.Lock()
	defer m.Unlock()

	for k := range m.queries {
		delete(m.queries, k)
	}
	for k := range m.dataLoader.loaderMap {
		delete(m.dataLoader.loaderMap, k)
	}

	if m.block != nil {
		err := m.block.Close()
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to close database connection", err, map[string]interface{}{})
		}
	}

	m.closeBatchOperation()

	return nil
}
