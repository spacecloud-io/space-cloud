package bolt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (b *Bolt) Delete(ctx context.Context, project, col string, req *model.DeleteRequest) (int64, error) {
	switch req.Operation {
	case utils.One:
		// Update single document
		if err := b.client.Update(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(project))
			c := bucket.Cursor()

			// get all keys matching the prefix
			prefix := []byte(col + "/")
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				result := map[string]interface{}{}
				if err := json.Unmarshal(v, result); err != nil {
					logrus.Errorf("error un marshalling while reading from bboltdb - %v", err)
					return err
				}
				// if valid then update
				if utils.Validate(req.Find, result) {
					// delete data
					if err := bucket.Delete(k); err != nil {
						return err
					}
					// exit the loop
					break
				}
			}
			return nil
		}); err != nil {
			return 0, nil
		}
		return 1, nil

	case utils.All:
		count := int64(0)
		// Update all matching document
		if err := b.client.Update(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(project))
			c := bucket.Cursor()

			// get all keys matching the prefix
			prefix := []byte(col + "/")
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				result := map[string]interface{}{}
				if err := json.Unmarshal(v, result); err != nil {
					logrus.Errorf("error un marshalling while reading from bboltdb - %v", err)
					return err
				}
				// if valid then delete
				if utils.Validate(req.Find, result) {
					count++
					// delete data
					if err := bucket.Delete(k); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return 0, nil
		}
		return count, nil

	default:
		return 0, errors.New("Invalid operation")
	}
}
