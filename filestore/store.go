package filestore

import (
	"context"
	"io"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/crud"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/filestore/local"
)

// Module is responsible for managing the file storage module
type Module struct {
	sync.RWMutex
	crud  *crud.Module
	store FileStore
}

// Init creates a new instance of the file store object
func Init(crud *crud.Module) *Module {
	return &Module{crud: crud}
}

// FileStore abstracts the implementation file storage operations
type FileStore interface {
	CreateFile(ctx context.Context, project string, req *model.CreateFileRequest, file io.Reader) error
	CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error

	ListDir(ctx context.Context, project string, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error)
	ReadFile(ctx context.Context, project, path string) (*model.File, error)

	DeleteDir(ctx context.Context, project, path string) error
	DeleteFile(ctx context.Context, project, path string) error

	GetStoreType() utils.FileStoreType
	Close() error
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(conf *config.FileStore) error {
	m.Lock()
	defer m.Unlock()

	// Close the previous store connections
	err := m.store.Close()
	if err != nil {
		return err
	}

	// Create a new crud blocks
	s, err := initBlock(utils.FileStoreType(conf.StoreType), conf.Connection)
	if err != nil {
		return err
	}
	m.store = s
	return nil
}

func initBlock(fileStoreType utils.FileStoreType, connection string) (FileStore, error) {
	switch fileStoreType {
	case utils.Local:
		return local.Init(connection)

	default:
		return nil, utils.ErrInvalidParams
	}
}
