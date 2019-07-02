package crud

import (
	"context"
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/modules/crud/sql"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	blocks    map[string]Crud
	primaryDB string
}

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

// Init create a new instance of the Module object
func Init() *Module {
	return &Module{blocks: make(map[string]Crud)}
}

func initBlock(dbType utils.DBType, connection string) (Crud, error) {
	switch dbType {
	case utils.Mongo:
		return mgo.Init(connection)

	case utils.MySQL, utils.Postgres:
		return sql.Init(dbType, connection)

	default:
		return nil, utils.ErrInvalidParams
	}
}

// GetPrimaryDB get the database configured as primary
func (m *Module) GetPrimaryDB() (Crud, error) {
	m.RLock()
	defer m.RUnlock()

	c, p := m.blocks[m.primaryDB]
	if !p {
		return nil, errors.New("CRUD: Primary DB not configured")
	}

	return c, nil
}

func (m *Module) getCrudBlock(dbType string) (Crud, error) {
	if crud, p := m.blocks[dbType]; p {
		return crud, nil
	}

	return nil, errors.New("CRUD: No crud block present for db")
}

// SetConfig set the rules adn secret key required by the crud block
func (m *Module) SetConfig(crud config.Crud) error {
	m.Lock()
	defer m.Unlock()

	// Close the previous database connections
	for _, v := range m.blocks {
		v.Close()
	}
	m.blocks = make(map[string]Crud, len(crud))

	// Create a new crud blocks
	for k, v := range crud {
		// Skip this block if it is not enabled
		if !v.Enabled {
			continue
		}

		c, err := initBlock(utils.DBType(k), v.Conn)
		if err != nil {
			return errors.New("CRUD: Error - " + k + " could not be initialised")
		}
		if v.IsPrimary {
			m.primaryDB = k
		}
		m.blocks[k] = c
	}

	return nil
}
