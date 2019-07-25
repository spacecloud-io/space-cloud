package filestore

import (
	"io"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/filestore/amazons3"
	"github.com/spaceuptech/space-cloud/modules/filestore/local"
	"github.com/spaceuptech/space-cloud/modules/filestore/gcpstorage"
)

// Module is responsible for managing the file storage module
type Module struct {
	sync.RWMutex
	store   FileStore
	enabled bool
	auth    *auth.Module
}

// Init creates a new instance of the file store object
func Init(auth *auth.Module) *Module {
	return &Module{enabled: false, store: nil, auth: auth}
}

// FileStore abstracts the implementation file storage operations
type FileStore interface {
	CreateFile(req *model.CreateFileRequest, file io.Reader) error
	CreateDir(req *model.CreateFileRequest) error

	ListDir(req *model.ListFilesRequest) ([]*model.ListFilesResponse, error)
	ReadFile(path string) (*model.File, error)

	DeleteDir(path string) error
	DeleteFile(path string) error

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

		// Close the store if present
		if m.store != nil {
			m.store.Close()
		}

		// Clear th store object
		m.store = nil
		return nil
	}

	// Create a new crud blocks
	s, err := initBlock(utils.FileStoreType(conf.StoreType), conf.Conn, conf.Endpoint, conf.Bucket)
	if err != nil {
		return err
	}
	m.store = s
	m.enabled = true
	return nil
}

// IsEnabled checks if the file store module is enabled
func (m *Module) IsEnabled() bool {
	m.RLock()
	defer m.RUnlock()
	return m.enabled
}

func initBlock(fileStoreType utils.FileStoreType, connection, endpoint, bucket string) (FileStore, error) {
	switch fileStoreType {
	case utils.Local:
		return local.Init(connection)
	case utils.AmazonS3:
		return amazons3.Init(connection, endpoint, bucket) // connection is the aws region code
	case utils.GCPStorage:
		return gcpstorage.Init(bucket)
	default:
		return nil, utils.ErrInvalidParams
	}
}
