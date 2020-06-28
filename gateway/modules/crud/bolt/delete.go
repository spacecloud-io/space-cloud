package bolt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Delete deletes a document (or multiple when op is "all") from the database
func (b *Bolt) Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error) {
	var count int64
	switch req.Operation {
	case utils.One, utils.All:
		if err := b.client.Update(func(tx *bbolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(b.bucketName))
			c := bucket.Cursor()

			// get all keys matching the prefix
			prefix := []byte(col + "/")
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				result := map[string]interface{}{}
				if err := json.Unmarshal(v, &result); err != nil {
					return fmt.Errorf("error un marshalling while reading from bboltdb - %v", err)
				}
				// if valid then delete
				if utils.Validate(req.Find, result) {
					// delete data
					if err := bucket.Delete(k); err != nil {
						return fmt.Errorf("error deleting from bboltdb - %v", err)
					}
					count++
					if req.Operation == utils.One {
						// exit the loop
						break
					}
				}
			}
			return nil
		}); err != nil {
			return 0, err
		}
		return count, nil

	default:
		return 0, errors.New("Invalid operation")
	}
}
