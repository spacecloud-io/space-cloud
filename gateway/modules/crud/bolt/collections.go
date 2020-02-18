package bolt

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"go.etcd.io/bbolt"
)

// GetCollections returns collection / tables name of specified database
func (b *Bolt) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {
	return nil, fmt.Errorf("error getting collection operation not supported for selected database")
}

// DeleteCollection deletes collection / tables name of specified database
func (b *Bolt) DeleteCollection(ctx context.Context, project, col string) error {
	err := b.client.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(project))

		if b == nil {
			return nil
		}

		c := b.Cursor()

		prefix := []byte(col)
		for key, _ := c.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, _ = c.Next() {
			err := b.Delete(key)
			if err != nil {
				logrus.Errorf("error deleting collection from embedded db - %s", err.Error())
				return err
			}
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("error deleting collection from embedded db - %s", err.Error())
		return err
	}
	return nil
}
