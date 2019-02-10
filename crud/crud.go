package crud

import (
	"context"
	"errors"
	"sync"

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
}

// Init create a new instance of the Module object
func Init() *Module {
	return &Module{blocks: make(map[string]Crud)}
}

func (mod *Module) getCrudBlock(dbType string) (Crud, error) {
	if crud, p := mod.blocks[dbType]; p {
		return crud, nil
	}

	return nil, errors.New("CRUD: No crud block present for db")
}
