package driver

import (
	"context"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/bolt"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/sql"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

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
	Close() error
}

// Handler is the object managing the database connections
type Handler struct {
	lock    sync.Mutex
	drivers map[string]*stub
}

// New creates a new driver handler
func New() *Handler {
	return &Handler{drivers: map[string]*stub{}}
}

// InitBlock creates and returns a new crud object. If the driver already exists, it returns that instead
func (h *Handler) InitBlock(dbType utils.DBType, enabled bool, connection string) (Crud, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	// See if driver is present in cache
	s, p := h.getBlock(dbType, connection)
	if p {
		s.addCount()
		return s.getCrud(), nil
	}

	return h.addBlock(dbType, enabled, connection)
}

// RemoveBlock removes a crud object if other module is referencing it
func (h *Handler) RemoveBlock(dbType utils.DBType, connection string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	// See if driver is present in cache
	s, p := h.getBlock(dbType, connection)
	if !p {
		return
	}

	// Subtract the count of references
	s.subtractCount()

	// Close then delete the driver if there are no references
	if s.getCount() == 0 {
		_ = s.getCrud().Close()
		delete(h.drivers, generateKey(dbType, connection))
	}
}

func (h *Handler) getBlock(dbType utils.DBType, connection string) (s *stub, p bool) {
	s, p = h.drivers[generateKey(dbType, connection)]
	return
}

func (h *Handler) addBlock(dbType utils.DBType, enabled bool, connection string) (Crud, error) {
	var c Crud
	var err error

	switch dbType {
	case utils.Mongo:
		c, err = mgo.Init(enabled, connection, "")

	case utils.EmbeddedDB:
		c, err = bolt.Init(enabled, connection, "")

	case utils.MySQL, utils.Postgres, utils.SQLServer:
		c, err = sql.Init(dbType, enabled, connection, "")

	default:
		c, err = nil, utils.ErrInvalidParams
	}

	h.drivers[generateKey(dbType, connection)] = newStub(c)
	return c, err
}

func generateKey(dbType utils.DBType, connection string) string {
	return string(dbType) + ":" + connection
}
