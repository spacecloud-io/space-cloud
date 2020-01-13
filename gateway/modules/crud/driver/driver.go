package driver

import (
	"context"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/sql"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

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
	CreateProjectIfNotExist(ctx context.Context, project string) error
	RawBatch(ctx context.Context, batchedQueries []string) error
	GetDBType() utils.DBType
	IsClientSafe() error
	Close() error
	GetConnectionState(ctx context.Context) bool
}

// Handler is the object managing the database connections
type Handler struct {
	lock               sync.Mutex
	drivers            map[string]*stub
	RemoveProjectScope bool
}

// New creates a new driver handler
func New(removeProjectScope bool) *Handler {
	return &Handler{drivers: map[string]*stub{}, RemoveProjectScope: removeProjectScope}
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
		c, err = mgo.Init(enabled, connection)

	case utils.MySQL, utils.Postgres, utils.SqlServer:
		c, err = sql.Init(dbType, enabled, h.RemoveProjectScope, connection)

	default:
		c, err = nil, utils.ErrInvalidParams
	}

	h.drivers[generateKey(dbType, connection)] = newStub(c)
	return c, err
}

func generateKey(dbType utils.DBType, connection string) string {
	return string(dbType) + ":" + connection
}
