package driver

import (
	"context"
	"sync"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/modules/crud/sql"
	"github.com/spaceuptech/space-cloud/utils"
)

// Crud abstracts the implementation crud operations of databases
type Crud interface {
	Create(ctx context.Context, project, col string, req *model.CreateRequest) error
	Read(ctx context.Context, project, col string, req *model.ReadRequest) (interface{}, error)
	Update(ctx context.Context, project, col string, req *model.UpdateRequest) error
	Delete(ctx context.Context, project, col string, req *model.DeleteRequest) error
	Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error)
	Batch(ctx context.Context, project string, req *model.BatchRequest) error
	GetDBType() utils.DBType
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
func (h *Handler) InitBlock(dbType utils.DBType, connection string) (Crud, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	// See if driver is present in cache
	s, p := h.getBlock(dbType, connection)
	if p {
		s.addCount()
		return s.getCrud(), nil
	}

	return h.addBlock(dbType, connection)
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
		s.getCrud().Close()
		delete(h.drivers, generateKey(dbType, connection))
	}
}

func (h *Handler) getBlock(dbType utils.DBType, connection string) (s *stub, p bool) {
	s, p = h.drivers[generateKey(dbType, connection)]
	return
}

func (h *Handler) addBlock(dbType utils.DBType, connection string) (Crud, error) {
	var c Crud
	var err error

	switch dbType {
	case utils.Mongo:
		c, err = mgo.Init(connection)

	case utils.MySQL, utils.Postgres:
		c, err = sql.Init(dbType, connection)

	default:
		c, err = nil, utils.ErrInvalidParams
	}

	// Return the error of exists
	if err != nil {
		return nil, err
	}

	h.drivers[generateKey(dbType, connection)] = newStub(c)
	return c, nil
}

func generateKey(dbType utils.DBType, connection string) string {
	return string(dbType) + ":" + connection
}
