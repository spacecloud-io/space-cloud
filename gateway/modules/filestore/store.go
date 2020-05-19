package filestore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"

	"github.com/spaceuptech/space-cloud/gateway/modules/filestore/amazons3"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore/gcpstorage"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore/local"
)

// Module is responsible for managing the file storage module
type Module struct {
	sync.RWMutex
	store       FileStore
	enabled     bool
	auth        model.AuthFilestoreInterface
	eventing    model.EventingModule
	metricsHook model.MetricFileHook

	// function to get secrets from runner
	getSecrets utils.GetSecrets
}

// Init creates a new instance of the file store object
func Init(auth model.AuthFilestoreInterface, hook model.MetricFileHook) *Module {
	return &Module{enabled: false, store: nil, auth: auth, metricsHook: hook}
}

// SetEventingModule sets the eventing module
func (m *Module) SetEventingModule(eventing model.EventingModule) {
	m.eventing = eventing
}

// FileStore abstracts the implementation file storage operations
type FileStore interface {
	CreateFile(req *model.CreateFileRequest, file io.Reader) error
	CreateDir(req *model.CreateFileRequest) error

	ListDir(req *model.ListFilesRequest) ([]*model.ListFilesResponse, error)
	ReadFile(path string) (*model.File, error)

	DeleteDir(path string) error
	DeleteFile(path string) error

	DoesExists(path string) error
	GetState(ctx context.Context) error

	GetStoreType() utils.FileStoreType
	Close() error
}

// SetConfig set the rules and secret key required by the filestore block
func (m *Module) SetConfig(project string, conf *config.FileStore) error {
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

		// Clear the store object
		m.store = nil
		return nil
	}

	// create aws and gcp file secret
	secretName, secretKey, isSecretExists := splitConnectionString(conf.Secret)
	if isSecretExists {
		value, err := m.getSecrets(project, secretName, secretKey)
		if err != nil {
			return utils.LogError("cannot get secrets from runner", "filestore", "setConfig", err)
		}
		if err := setFileSecret(utils.FileStoreType(conf.StoreType), secretKey, value); err != nil {
			return utils.LogError("cannot set fileStore secrets", "filestore", "setConfig", err)
		}
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

func setFileSecret(fileStoreType utils.FileStoreType, key, value string) error {
	switch fileStoreType {
	case utils.AmazonS3:
		path := fmt.Sprintf("%s/.aws/", os.ExpandEnv("$HOME"))
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(fmt.Sprintf("%s/credentials", path), []byte(value), 0755)
	case utils.GCPStorage:
		path := fmt.Sprintf("%s/.gcp/", os.ExpandEnv("$HOME"))
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(fmt.Sprintf("%s/credentials.json", path), []byte(value), 0755)
	default:
		return utils.ErrInvalidParams
	}
}

// splitConnectionString splits the connection string
func splitConnectionString(connection string) (string, string, bool) {
	s := strings.Split(connection, ".")
	if s[0] == "secrets" {
		return s[1], s[2], true
	}
	return "", "", false
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

// SetGetSecrets sets the GetSecrets function
func (m *Module) SetGetSecrets(function utils.GetSecrets) {
	m.Lock()
	defer m.Unlock()

	m.getSecrets = function
}
