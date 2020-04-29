package letsencrypt

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mholt/certmagic"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("error while storing in lets encrypt - %v", err)
		return err
	}
	if response.Status != http.StatusOK {
		logrus.Errorf("error while storing in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
		return fmt.Errorf("error while storing in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
	}
	return nil
}

// Load gets the value for specifed key
func (s *Storage) Load(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		logrus.Errorf("error while getting data in lets encrypting %v", err)
		return nil, err
	}
	if response.Status != http.StatusOK {
		logrus.Errorf("error while getting data in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
		return nil, fmt.Errorf("error while getting data in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
	}

	result, ok := response.Data["result"]
	if !ok {
		logrus.Errorf("error while getting data in lets encrypt unable to find field result in response body")
		return nil, fmt.Errorf("error while getting data in lets encrypt unable to find field result in response body")
	}
	data, ok := result.(map[string]interface{})
	if !ok {
		logrus.Errorf("error while getting data in lets encrypt unable to assert result to map")
		return nil, fmt.Errorf("error while getting data in lets encrypt unable to assert result to map")
	}
	value, ok := data["value"]
	if !ok {
		logrus.Errorf("error while getting data in lets encrypt unable to find field value in received data")
		return nil, fmt.Errorf("error while getting data in lets encrypt unable to find field value in received data")
	}

	return base64.StdEncoding.DecodeString(value.(string))
}

// Delete deletes specifed key from space cloud storage
func (s *Storage) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Delete(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		logrus.Errorf("error while deleting in lets encrypt - %v", err)
		return err
	}
	if response.Status != http.StatusOK {
		logrus.Errorf("error while deleting lets encrypt got http status %v with error message - %v", response.Status, response.Error)
		return fmt.Errorf("error while deleting lets encrypt got http status %v with error message - %v", response.Status, response.Error)
	}
	return nil
}

// Exists checks if key exists in space cloud storage
func (s *Storage) Exists(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Get(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil {
		logrus.Errorf("error while checking existence in lets encrypt - %v", err)
		return false
	}
	if response.Status != http.StatusOK {
		logrus.Errorf("error while checking existence of in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
		return false
	}
	result, ok := response.Data["result"]
	if !ok {
		logrus.Errorf("error while checking existence in lets encrypt unable to find field result in response body")
		return false
	}
	data, ok := result.([]interface{})
	if !ok {
		logrus.Errorf("error while checking existence in lets encrypt unable to assert result to array")
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
		logrus.Errorf("error while listing in lets encrypt - %v", err)
		return nil, err
	}
	if response.Status != http.StatusOK {
		logrus.Errorf("error while listing response of lets encrypt got http status %v with error message - %v", response.Status, response.Error)
		return nil, fmt.Errorf("error while listing response of lets encrypt got http status %v with error message - %v", response.Status, response.Error)
	}

	result, ok := response.Data["result"]
	if !ok {
		logrus.Errorf("error while listing in lets encrypt unable to find field result in response body")
		return nil, fmt.Errorf("error while listing in lets encrypt unable to find field result in response body")
	}
	data, ok := result.([]interface{})
	if !ok {
		logrus.Errorf("error while listing in lets encrypt unable to assert result to array")
		return nil, fmt.Errorf("error while listing in lets encrypt unable to assert result to array")
	}
	prefixArr := []string{}
	for _, v := range data {
		key, ok := v.(map[string]interface{})["_id"]
		if !ok {
			logrus.Errorf("error while listing in lets encrypt unable to find _id field in received data")
			return nil, fmt.Errorf("error while listing in lets encrypt unable to find _id field in received data")
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
		logrus.Errorf("error while getting stats in lets encrypt - %v", err)
		return certmagic.KeyInfo{}, err
	}
	if response.Status != http.StatusOK {
		logrus.Errorf("error while getting stats in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
		return certmagic.KeyInfo{}, fmt.Errorf("error while getting stats in lets encrypt got http status %v with error message - %v", response.Status, response.Error)
	}

	modifiedTime, err := time.Parse(time.RFC3339, response.Data["modified"].(string))
	if err != nil {
		return certmagic.KeyInfo{}, fmt.Errorf("error getting stats in lets encrypt unable to parse string to time - %v", err)
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
func (s *Storage) Lock(key string) error {
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
			logrus.Errorf("error creating lock file in lets encrypt - %v", err)
			return fmt.Errorf("error creating lock file in lets encrypt - %v", err)
		}

		// lock file already exists
		info, err := s.Stat(lockFile)
		switch {
		case err != nil:
			logrus.Errorf("error getting stats of lock file in lets encrypt - %v", err)
			return fmt.Errorf("error getting stats of lock file in lets encrypt - %v", err)

		case s.fileLockIsStale(info):
			logrus.Printf("error in lets encrypt lock file is in stale state removing and trying again")
			if err := s.deleteLockFile(lockFile); err != nil {
				return err
			}
			continue

		case time.Since(start) > staleLockDuration*2:
			// should never happen, hopefully
			logrus.Errorf("possible deadlock: %s passed trying to obtain lock for %s", time.Since(start), key)
			return fmt.Errorf("possible deadlock: %s passed trying to obtain lock for %s", time.Since(start), key)

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
		logrus.Errorf("error while creating lock file in lets encrypt - %v", err)
	}
	return err
}

func (s *Storage) deleteLockFile(keyPath string) error {
	err := s.Delete(keyPath)
	if err != nil {
		logrus.Errorf("error while deleting lock file in lets encrypt - %v", err)
		return fmt.Errorf("error while deleting lock file in lets encrypt - %v", err)
	}
	return nil
}
