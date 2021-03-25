package letsencrypt

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/spaceuptech/helpers"
	apigo "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"
	"github.com/spaceuptech/space-api-go/types"
)

// Storage is the object for storing space cloud storage information
type Storage struct {
	sync.RWMutex
	db         *db.DB
	path       string // required for lock
	collection string
}

// NewScStore returns a new instance of space cloud storage
func NewScStore() *Storage {
	scProject := os.Getenv("LETSENCRYPT_SC_PROJECT")
	scAddr := os.Getenv("LETSENCRYPT_SC_ADDR")
	scDatabase := os.Getenv("LETSENCRYPT_SC_DATABASE")
	scCollection := os.Getenv("LETSENCRYPT_SC_COLLECTION")
	if scProject == "" {
		scProject = "space_cloud"
	}
	if scAddr == "" {
		scAddr = "store.space_cloud.svc.cluster.local:4122"
	}
	if scDatabase == "" {
		scDatabase = "bolt"
	}
	if scCollection == "" {
		scCollection = "certificates"
	}

	return &Storage{db: apigo.New(scProject, scAddr, false).DB(scDatabase), collection: scCollection, path: "certmagic"}
}

// Store sets the key value in space cloud storage
func (s *Storage) Store(key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m := types.M{
		"_id":   key,
		"value": base64.StdEncoding.EncodeToString(value),
		"size":  len(value),
	}

	response, err := s.db.Upsert(s.collection).Set(m).Apply(ctx)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "error while storing in lets encrypt", err, nil)
	}
	if response.Status != http.StatusOK {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to store lets encrypt config database service responded with status code (%v)", response.Status), fmt.Errorf(response.Error), nil)
	}
	return nil
}

// Load gets the value for specifed key
func (s *Storage) Load(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to fetch lets encrypt key (%s)", key), err, map[string]interface{}{"collection": s.collection})
	}
	if response.Status != http.StatusOK {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to fetch lets encrypt key (%s) database service responded with status code (%v)", key, response.Status), fmt.Errorf(response.Error), nil)
	}

	result, ok := response.Data["result"]
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found while fetching lets encypt key (%s)", key), fmt.Errorf("field (result) not found in data object"), nil)
	}
	data, ok := result.(map[string]interface{})
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data type found for lets encypt key (%s)", key), fmt.Errorf("field (result) is of type (%v) but wanted object", reflect.TypeOf(result)), nil)
	}
	value, ok := data["value"]
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found for lets encypt key (%s)", key), fmt.Errorf("field (value) not found in result object"), nil)
	}

	return base64.StdEncoding.DecodeString(value.(string))
}

// Delete deletes specified key from space cloud storage
func (s *Storage) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Delete(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to delete lets encrypt key (%s)", key), err, nil)
	}
	if response.Status != http.StatusOK {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to delete lets encrypt key (%s) database service responded with status (%v)", key, response.Status), err, nil)
	}
	return nil
}

// Exists checks if key exists in space cloud storage
func (s *Storage) Exists(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Get(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable check if lets encrypt key (%s) exists", key), err, nil)
		return false
	}
	if response.Status != http.StatusOK {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable check if lets encrypt key (%s) exists database service responded with status code (%v)", key, response.Status), err, nil)
		return false
	}
	result, ok := response.Data["result"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found while checking for existence of lets encypt key (%s)", key), fmt.Errorf("field (result) not found in data object"), nil)
		return false
	}
	data, ok := result.([]interface{})
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found while checking for existence of lets encypt key (%s)", key), fmt.Errorf("field (result) is not an array got (%v)", reflect.TypeOf(result)), nil)
		return false
	}
	if len(data) > 0 {
		return true
	}
	return false
}

// List returns all keys matching prefix
func (s *Storage) List(prefix string, recursive bool) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Get(s.collection).Where(types.Cond("_id", "regex", fmt.Sprintf("^%s", prefix))).Apply(ctx)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to lists lets encypt keys with prefix (%s)", prefix), err, map[string]interface{}{"isRecursive": recursive})
	}
	if response.Status != http.StatusOK {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to lists lets encypt keys with prefix (%s) database service responded with status code (%v)", prefix, response.Status), err, map[string]interface{}{"isRecursive": recursive})
	}

	result, ok := response.Data["result"]
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found while listing lets encypt keys with prefix (%s)", prefix), fmt.Errorf("field (result) not found in data object"), nil)
	}
	data, ok := result.([]interface{})
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found while listing lets encypt keys with prefix (%s)", prefix), fmt.Errorf("field (result) is not an array got (%v)", reflect.TypeOf(result)), nil)
	}
	prefixArr := []string{}
	for _, v := range data {
		key, ok := v.(map[string]interface{})["_id"]
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid data found while listing lets encypt keys with prefix (%s)", prefix), fmt.Errorf("field (_id) not found in data object"), nil)
		}
		prefixArr = append(prefixArr, key.(string))
	}
	return prefixArr, nil
}

// Stat get stats for a particular key
func (s *Storage) Stat(key string) (certmagic.KeyInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		return certmagic.KeyInfo{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to get stats of lets encrypt key (%s)", key), err, nil)
	}
	if response.Status != http.StatusOK {
		return certmagic.KeyInfo{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to get stats of lets encrypt key (%s) database service responded with status code (%v)", key, response.Status), fmt.Errorf(response.Error), nil)
	}

	modifiedTime, err := time.Parse(time.RFC3339Nano, response.Data["modified"].(string))
	if err != nil {
		return certmagic.KeyInfo{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to parse (modified) field of lets encrypt key (%s) time to string ", key), err, nil)
	}

	return certmagic.KeyInfo{
		Key:        key,
		Modified:   modifiedTime,
		Size:       response.Data["size"].(int64),
		IsTerminal: true,
	}, nil
}

const lockFileExists = "Lock file already exists"

// staleLockDuration is the length of time
// before considering a lock to be stale.
const staleLockDuration = 2 * time.Hour

// fileLockPollInterval is how frequently
// to check the existence of a lock file
const fileLockPollInterval = 1 * time.Second

// StorageKeys is used to store certmagic keys
var StorageKeys certmagic.KeyBuilder

// Lock obtains a lock named by the given key. It blocks
// until the lock can be obtained or an error is returned.
func (s *Storage) Lock(ctx context.Context, key string) error {
	start := time.Now()
	lockFile := s.lockFileName(key)

	for {
		err := s.createLockFile(lockFile)
		if err == nil {
			// got the lock
			return nil
		}

		if err.Error() != lockFileExists {
			// unexpected error
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create lock file in lets encrypt", err, nil)
		}

		// lock file already exists
		info, err := s.Stat(lockFile)
		switch {
		case err != nil:
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to get stats of lock file in lets encrypt", err, nil)

		case s.fileLockIsStale(info):
			helpers.Logger.LogWarn(helpers.GetRequestID(context.TODO()), "lets encrypt lock file is in stale state removing and trying again", nil)
			if err := s.deleteLockFile(lockFile); err != nil {
				return err
			}
			continue

		case time.Since(start) > staleLockDuration*2:
			// should never happen, hopefully
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Lets encrypt dead lock by passing", fmt.Errorf("possible deadlock: %s passed trying to obtain lock for %s", time.Since(start), key), nil)

		default:
			// lockfile exists and is not stale;
			// just wait a moment and try again
			time.Sleep(fileLockPollInterval)

		}
	}
}

// Unlock releases the lock for name.
func (s *Storage) Unlock(key string) error {
	return s.deleteLockFile(s.lockFileName(key))
}

func (s *Storage) String() string {
	return "storage:" + s.path
}

func (s *Storage) lockFileName(key string) string {
	return filepath.Join(s.lockDir(), fmt.Sprintf("%s.lock", StorageKeys.Safe(key)))
}

func (s *Storage) lockDir() string {
	return filepath.Join(s.path, "locks")
}

func (s *Storage) fileLockIsStale(info certmagic.KeyInfo) bool {
	return time.Since(info.Modified) > staleLockDuration
}

func (s *Storage) createLockFile(filename string) error {
	exists := s.Exists(filename)
	if exists {
		return fmt.Errorf(lockFileExists)
	}

	err := s.Store(filename, []byte("lock"))
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create lock file in lets encrypt", err, nil)
	}
	return err
}

func (s *Storage) deleteLockFile(keyPath string) error {
	err := s.Delete(keyPath)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete lock file of lets encrypt", err, nil)
	}
	return nil
}
