package bolt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/spaceuptech/helpers"
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
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to unmarshal data read from bbbolt db", err, nil)
				}
				// if valid then update
				if utils.Validate(string(model.EmbeddedDB), req.Find, currentObj) {
					objToSet, ok := req.Update["$set"].(map[string]interface{})
					if !ok {
						return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to update in bbolt - $set db operator not found or the operator value is not map", nil, nil)
					}

					for objToSetKey, objToSetValue := range objToSet {
						currentObj[objToSetKey] = objToSetValue
					}
					value, err := json.Marshal(&currentObj)
					if err != nil {
						return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to unmarshal data updated from bbbolt db", err, nil)
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
				return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to update in bbolt - $set db operator not found or the operator value is not map", nil, nil)
			}

			for findName, findValue := range req.Find {
				_, ok := objToSet[findName]
				if !ok {
					objToSet[findName] = findValue
				}
			}
			rowsAffected, err := b.Create(ctx, col, &model.CreateRequest{Operation: utils.One, Document: objToSet})
			if err != nil || rowsAffected == 0 {
				return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to upsert in bbolt db - %v rows affected %v", err, rowsAffected), nil, nil)
			}
			count = rowsAffected
		}
		return count, nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
