package filestore

import (
	"context"
	"io"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/filestore/local"
)

// Module is responsible for managing the file storage module
type Module struct {
	sync.RWMutex
	store   FileStore
	enabled bool
}

// Init creates a new instance of the file store object
func Init() *Module {
	return &Module{enabled: false, store: nil}
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
	if m.store != nil {
		err := m.store.Close()
		if err != nil {
			return err
		}
	}

	// Disable the module if file store is not enabled
	if conf == nil || !conf.Enabled {
		m.enabled = false
		m.store = nil
		return nil
	}

	// Create a new crud blocks
	s, err := initBlock(utils.FileStoreType(conf.StoreType), conf.Conn)
	if err != nil {
		return err
	}
	m.store = s
	m.enabled = true
	return nil
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()
	return m.enabled
}

func initBlock(fileStoreType utils.FileStoreType, connection string) (FileStore, error) {
	switch fileStoreType {
	case utils.Local:
		return local.Init(connection)

	default:
		return nil, utils.ErrInvalidParams
	}
}
