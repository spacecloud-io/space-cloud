package bolt

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/spaceuptech/helpers"
	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetCollections returns collection / tables name of specified database
func (b *Bolt) GetCollections(ctx context.Context) ([]utils.DatabaseCollections, error) {
	keys := make(map[string]bool)
	err := b.client.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
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
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to query database to get tables in database (%s)", b.bucketName), err, nil)
	}
	dbCols := make([]utils.DatabaseCollections, 0)
	for col := range keys {
		dbCols = append(dbCols, utils.DatabaseCollections{TableName: col})
	}
	return dbCols, nil
}

// DeleteCollection deletes collection / tables name of specified database
func (b *Bolt) DeleteCollection(ctx context.Context, col string) error {
	err := b.client.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))

		if b == nil {
			return nil
		}

		c := b.Cursor()

		prefix := []byte(col)
		for key, _ := c.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, _ = c.Next() {
			err := b.Delete(key)
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error deleting collection from embedded db", err, nil)
			}
		}
		return nil
	})
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error deleting collection from embedded db", err, nil)
	}
	return nil
}
