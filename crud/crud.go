package crud

import (
	"context"
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/crud/mgo"
	"github.com/spaceuptech/space-cloud/crud/sql"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	blocks map[string]Crud
}

// Crud abstracts the implementation crud operations of databases
type Crud interface {
	Create(ctx context.Context, project, col string, req *model.CreateRequest) error
	Read(ctx context.Context, project, col string, req *model.ReadRequest) (interface{}, error)
	Update(ctx context.Context, project, col string, req *model.UpdateRequest) error
	Delete(ctx context.Context, project, col string, req *model.DeleteRequest) error
	Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error)
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
		c, err := initBlock(utils.DBType(k), v.Connection)
		if err != nil {
			return err
		}
		m.blocks[k] = c
	}

	return nil
}
