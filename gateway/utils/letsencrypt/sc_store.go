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

type storage struct {
	sync.RWMutex
	db         *db.DB
	path       string // required for lock
	collection string
}

func NewScStore() *storage {
	scProject := os.Getenv("LETSENCRYPT_SC_PROJECT")
	scAddr := os.Getenv("LETSENCRYPT_SC_ADDR")
	scDatabase := os.Getenv("LETSENCRYPT_SC_DATABASE")
	scCollection := os.Getenv("LETSENCRYPT_SC_COLLECTION")
	if scProject == "" {
		scProject = "space_cloud"
	}
	if scAddr == "" {
		scAddr = "store.space-cloud.svc.cluster.local:4122"
	}
	if scDatabase == "" {
		scDatabase = "bolt"
	}
	if scCollection == "" {
		scCollection = "certificates"
	}

	return &storage{db: apigo.New(scProject, scAddr, false).DB(scDatabase), collection: scCollection}
}

func (s *storage) Store(key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m := types.M{
		"_id":   key,
		"value": base64.StdEncoding.EncodeToString(value),
		"size":  len(value),
	}

	data, err := s.db.Upsert(s.collection).Set(m).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error while storing data of lets encrypt - %v %v", err, data.Error)
	}
	return err
}

func (s *storage) Load(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || response.Status != http.StatusOK {
		logrus.Errorf("error while retrieving response of lets encrypt - %v", err)
		return nil, err
	}

	result, ok := response.Data["result"]
	if !ok {
		logrus.Errorf("error while getting data of lets encrypt unable to find field result in response body")
		return nil, fmt.Errorf("error while getting data of lets encrypt unable to find field result in response body")
	}
	data, ok := result.(map[string]interface{})
	if !ok {
		logrus.Errorf("error while getting data of lets encrypt unable to assert result to map")
		return nil, fmt.Errorf("error while getting data of lets encrypt unable to assert result to map")
	}
	value, ok := data["value"]
	if !ok {
		logrus.Errorf("error while getting data of lets encrypt unable to find field value in received data")
		return nil, fmt.Errorf("error while getting data of lets encrypt unable to find field value in received data")
	}

	return base64.StdEncoding.DecodeString(value.(string))
}

func (s *storage) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.db.Delete(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error while deleting data of lets encrypt - %v %v", err, data.Error)
	}
	return err
}

func (s *storage) Exists(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Get(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || response.Status != http.StatusOK {
		logrus.Errorf("error while checking existence of in lets encrypt unable to specified key - %v %v", err, response.Error)
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
	logrus.Errorf("error while checking existence in lets encrypt length less than zero")
	return false
}

func (s *storage) List(prefix string, recursive bool) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.db.Get(s.collection).Where(types.Cond("_id", "regex", prefix)).Apply(ctx)
	if err != nil || response.Status != http.StatusOK {
		logrus.Errorf("error while listing response of lets encrypt - %v", err)
		return nil, err
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

func (s *storage) Stat(key string) (certmagic.KeyInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error getting stats in lets encrypt - %v", err)
		return certmagic.KeyInfo{}, err
	}

	modifiedTime, err := time.Parse(time.RFC3339, data.Data["modified"].(string))
	if err != nil {
		return certmagic.KeyInfo{}, fmt.Errorf("error getting stats in lets encrypt unable to parse string to time - %v", err)
	}

	return certmagic.KeyInfo{
		Key:        key,
		Modified:   modifiedTime,
		Size:       data.Data["size"].(int64),
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

var StorageKeys certmagic.KeyBuilder

// Lock obtains a lock named by the given key. It blocks
// until the lock can be obtained or an error is returned.
func (s *storage) Lock(key string) error {
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
			s.deleteLockFile(lockFile)
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
func (s *storage) Unlock(key string) error {
	return s.deleteLockFile(s.lockFileName(key))
}

func (s *storage) String() string {
	return "storage:" + s.path
}

func (s *storage) lockFileName(key string) string {
	return filepath.Join(s.lockDir(), fmt.Sprintf("%s.lock", StorageKeys.Safe(key)))
}

func (s *storage) lockDir() string {
	return filepath.Join(s.path, "locks")
}

func (s *storage) fileLockIsStale(info certmagic.KeyInfo) bool {
	return time.Since(info.Modified) > staleLockDuration
}

func (s *storage) createLockFile(filename string) error {
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

func (s *storage) deleteLockFile(keyPath string) error {
	err := s.Delete(keyPath)
	if err != nil {
		logrus.Errorf("error while deleting lock file in lets encrypt - %v", err)
	}
	return err
}
