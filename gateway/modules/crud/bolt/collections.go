package bolt

import (
	"context"
	"fmt"
	"strings"

	"go.etcd.io/bbolt"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetCollections returns collection / tables name of specified database
func (b *Bolt) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {
	keys := make(map[string]bool)
	err := b.client.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(project))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for key, _ := c.First(); key != nil; key, _ = c.Next() {
			keys[strings.Split(string(key), "/")[0]] = true
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("could not get all collections for given project and database - %s", err.Error())
		return nil, err
	}
	dbCols := make([]utils.DatabaseCollections, len(keys))
	for col := range keys {
		dbCols = append(dbCols, utils.DatabaseCollections{TableName: col})
	}
	return dbCols, nil
}

// DeleteCollection deletes collection / tables name of specified database
func (b *Bolt) DeleteCollection(ctx context.Context, project, col string) error {
	return fmt.Errorf("error deleting collection operation not supported for selected database")
}
