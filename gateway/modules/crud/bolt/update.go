package bolt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Update updates the document(s) which match the condition provided.
func (b *Bolt) Update(ctx context.Context, col string, req *model.UpdateRequest) (int64, error) {
	var count int64
	switch req.Operation {
	case utils.One, utils.All, utils.Upsert:
		if err := b.client.Update(func(tx *bbolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(b.bucketName))
			c := bucket.Cursor()

			// get all keys matching the prefix
			prefix := []byte(col + "/")
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				currentObj := map[string]interface{}{}
				if err := json.Unmarshal(v, &currentObj); err != nil {
					logrus.Errorf("error un marshalling while reading from bboltdb - %v", err)
					return err
				}
				// if valid then update
				if utils.Validate(req.Find, currentObj) {
					objToSet, ok := req.Update["$set"].(map[string]interface{})
					if !ok {
						return fmt.Errorf("error unable to update in bbolt - $set db operator not found or the operator value is not map")
					}

					for objToSetKey, objToSetValue := range objToSet {
						currentObj[objToSetKey] = objToSetValue
					}
					value, err := json.Marshal(&currentObj)
					if err != nil {
						logrus.Errorf("error marshalling while updating in bboltdb - %v", err)
						return err
					}

					// over ride the data
					if err = bucket.Put(k, value); err != nil {
						return err
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

		if req.Operation == utils.Upsert && count == 0 {
			objToSet, ok := req.Update["$set"].(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("error unable to update in bbolt - $set db operator not found or the operator value is not map")
			}

			for findName, findValue := range req.Find {
				_, ok := objToSet[findName]
				if !ok {
					objToSet[findName] = findValue
				}
			}
			rowsAffected, err := b.Create(ctx, col, &model.CreateRequest{Operation: utils.One, Document: objToSet})
			if err != nil || rowsAffected == 0 {
				return 0, fmt.Errorf("error while upserting in bbolt db - %v rows affected %v", err, rowsAffected)
			}
			count = rowsAffected
		}
		return count, nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
