package crud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

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
	block   Crud
	dbType  string
	alias   string
	project string
	schema  model.SchemaCrudInterface
	queries map[string]*config.PreparedQuery
	// batch operation
	batchMapTableToChan batchMap // every table gets mapped to group of channels

	dataLoader loader
	// Variables to store the hooks
	hooks      *model.CrudHooks
	metricHook model.MetricCrudHook

	// function to get secrets from runner
	getSecrets utils.GetSecrets
}

type loader struct {
	loaderMap      map[string]*dataloader.Loader
	dataLoaderLock sync.Mutex
}

// Crud abstracts the implementation crud operations of databases
type Crud interface {
	Create(ctx context.Context, col string, req *model.CreateRequest) (int64, error)
	Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, error)
	Update(ctx context.Context, col string, req *model.UpdateRequest) (int64, error)
	Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error)
	Aggregate(ctx context.Context, col string, req *model.AggregateRequest) (interface{}, error)
	Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error)
	DescribeTable(ctc context.Context, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error)
	RawExec(ctx context.Context, query string) error
	RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error)
	GetCollections(ctx context.Context) ([]utils.DatabaseCollections, error)
	DeleteCollection(ctx context.Context, col string) error
	CreateDatabaseIfNotExist(ctx context.Context, name string) error
	RawBatch(ctx context.Context, batchedQueries []string) error
	GetDBType() utils.DBType
	IsClientSafe() error
	IsSame(conn, dbName string) bool
	Close() error
	GetConnectionState(ctx context.Context) bool
}

// Init create a new instance of the Module object
func Init() *Module {
	return &Module{batchMapTableToChan: make(batchMap), dataLoader: loader{loaderMap: map[string]*dataloader.Loader{}}}
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

func (m *Module) initBlock(dbType utils.DBType, enabled bool, connection, dbName string) (Crud, error) {
	switch dbType {
	case utils.Mongo:
		return mgo.Init(enabled, connection, dbName)
	case utils.EmbeddedDB:
		return bolt.Init(enabled, connection, dbName)
	case utils.MySQL, utils.Postgres, utils.SQLServer:
		c, err := sql.Init(dbType, enabled, connection, dbName)
		if err == nil && enabled {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := c.CreateDatabaseIfNotExist(ctx, dbName); err != nil {
				return nil, err
			}
		}
		if dbType == utils.MySQL {
			return sql.Init(dbType, enabled, fmt.Sprintf("%s%s", connection, dbName), dbName)
		}
		return c, err
	default:
		return nil, utils.ErrInvalidParams
	}
}

func (m *Module) getCrudBlock(dbAlias string) (Crud, error) {
	if m.block != nil && m.alias == dbAlias {
		return m.block, nil
	}
	return nil, fmt.Errorf("crud module not initialized yet for %q", dbAlias)
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crud config.Crud) error {
	m.Lock()
	defer m.Unlock()

	if len(crud) > 1 {
		return errors.New("crud module cannot have more than 1 db")
	}

	m.closeBatchOperation()
	m.project = project

	// Reset all existing prepared query
	m.queries = map[string]*config.PreparedQuery{}

	// clear previous data loader
	m.dataLoader = loader{loaderMap: map[string]*dataloader.Loader{}}

	// Create a new crud blocks
	for k, v := range crud {
		if v.Type == "" {
			v.Type = k
		}

		// set default database name to project id
		if v.DBName == "" {
			v.DBName = project
		}

		if m.block != nil {
			// Skip if the connection string is the same
			if m.block.IsSame(v.Conn, v.DBName) {
				break
			}

			// Close the previous database connection
			if err := m.block.Close(); err != nil {
				_ = utils.LogError("Unable to close database connections", "crud", "set-config", err)
			}
		}

		var c Crud
		var err error

		// check if connection string starts with secrets
		secretName, secretKey, isSecretExists := splitConnectionString(v.Conn)
		if isSecretExists {
			v.Conn, err = m.getSecrets(project, secretName, secretKey)
			if err != nil {
				return utils.LogError("cannot get secrets from runner", "crud", "setConfig", err)
			}
		}

		v.Type = strings.TrimPrefix(v.Type, "sql-")
		c, err = m.initBlock(utils.DBType(v.Type), v.Enabled, v.Conn, v.DBName)

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

		// Add the prepared queries in this db
		for id, query := range v.PreparedQueries {
			m.queries[getPreparedQueryKey(strings.TrimPrefix(k, "sql-"), id)] = query
		}
	}
	m.initBatchOperation(project, crud)
	return nil
}

// splitConnectionString splits the connection string
func splitConnectionString(connection string) (string, string, bool) {
	s := strings.Split(connection, ".")
	if s[0] == "secrets" {
		return s[1], s[2], true
	}
	return "", "", false
}

// GetDBType returns the type of the db for the alias provided
func (m *Module) GetDBType(dbAlias string) (string, error) {
	dbAlias = strings.TrimPrefix(dbAlias, "sql-")
	if dbAlias != m.alias {
		return "", fmt.Errorf("db (%s) not found", dbAlias)
	}
	return m.dbType, nil
}

// SetGetSecrets sets the GetSecrets function
func (m *Module) SetGetSecrets(function utils.GetSecrets) {
	m.Lock()
	defer m.Unlock()

	m.getSecrets = function
}
