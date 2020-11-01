package filestore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/spaceuptech/helpers"

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
	CreateFile(ctx context.Context, req *model.CreateFileRequest, file io.Reader) error
	CreateDir(ctx context.Context, req *model.CreateFileRequest) error

	ListDir(ctx context.Context, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error)
	ReadFile(ctx context.Context, path string) (*model.File, error)

	DeleteDir(ctx context.Context, path string) error
	DeleteFile(ctx context.Context, path string) error

	DoesExists(ctx context.Context, path string) error
	GetState(ctx context.Context) error

	GetStoreType() utils.FileStoreType
	Close() error
}

// SetConfig set the rules and secret key required by the filestore block
func (m *Module) SetConfig(project string, conf *config.FileStoreConfig) error {
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
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to fetch secret from runner", err, nil)
		}
		if err := setFileSecret(utils.FileStoreType(conf.StoreType), secretKey, value); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create credential file in gateway", err, nil)
		}
	}

	// Create a new crud blocks
	s, err := initBlock(conf)
	if err != nil {
		return err
	}
	m.store = s
	m.enabled = true
	return nil
}

// CloseConfig closes file store
func (m *Module) CloseConfig() error {
	m.Lock()
	defer m.Unlock()

	if m.store == nil {
		return nil
	}

	return m.store.Close()
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
		return s[1], strings.Join(s[2:], "."), true
	}
	return "", "", false
}

// IsEnabled checks if the file store module is enabled
func (m *Module) IsEnabled() bool {
	m.RLock()
	defer m.RUnlock()
	return m.enabled
}

func initBlock(conf *config.FileStoreConfig) (FileStore, error) {
	switch utils.FileStoreType(conf.StoreType) {
	case utils.Local:
		return local.Init(conf.Conn)
	case utils.AmazonS3:
		return amazons3.Init(conf.Conn, conf.Endpoint, conf.Bucket, conf.DisableSSL, conf.ForcePathStyle) // connection is the aws region code
	case utils.GCPStorage:
		return gcpstorage.Init(conf.Bucket)
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
