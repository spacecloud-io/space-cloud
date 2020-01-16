package letsencrypt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
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
	collection string
}

func NewScStore(project, url, db string) *storage {
	return &storage{db: apigo.New(project, url, false).DB(db), collection: "certificates"}
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
		logrus.Errorf("error while storing data of lets encrypt - %v", err)
		return err
	}
	return nil
}

func (s *storage) Load(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error while retrieving data of lets encrypt - %v", err)
		return nil, err
	}
	return json.Marshal(data.Data)
}

func (s *storage) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.db.DeleteOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error while deleting data of lets encrypt - %v", err)
		return err
	}
	return nil
}

func (s *storage) Exists(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error while checking existence of in lets encrypt unable to specified key - %v", err)
		return false
	}
	// 4 is the number of columns in table so the map should contain 4 fields
	if len(data.Data) != 4 {
		logrus.Errorf("error while checking existence of specified key of lets encrypt")
		return false
	}
	return true
}

func (s *storage) List(prefix string, recursive bool) ([]string, error) {
	// todo verify this
	return []string{prefix}, nil
}

func (s *storage) Stat(key string) (certmagic.KeyInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.db.GetOne(s.collection).Where(types.M{"_id": key}).Apply(ctx)
	if err != nil || data.Status != http.StatusOK {
		logrus.Errorf("error while retrieving data of lets encrypt - %v", err)
		return certmagic.KeyInfo{}, err
	}

	modifiedTime, err := time.Parse(time.RFC3339, data.Data["modified"].(string))
	if err != nil {
		return certmagic.KeyInfo{}, fmt.Errorf("error while generating stats in lets encrypt unable to parse string to time - %v", err)
	}

	return certmagic.KeyInfo{
		Key:        key,
		Modified:   modifiedTime,
		Size:       data.Data["size"].(int64),
		IsTerminal: false, // todo check this
	}, nil
}

func (s *storage) Lock(key string) error {
	return nil
}

func (s *storage) Unlock(key string) error {
	return nil
}
