package bolt

import (
	"context"
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Create inserts a document (or multiple when op is "all") into the database
func (b *Bolt) Create(ctx context.Context, col string, req *model.CreateRequest) (int64, error) {
	objs := []interface{}{}
	switch req.Operation {
	case utils.All, utils.One:
		if req.Operation == utils.One {
			doc, ok := req.Document.(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("error inserting into bboltdb cannot assert document to map")
			}
			objs = append(objs, doc)
		} else {
			docs, ok := req.Document.([]interface{})
			if !ok {
				return 0, fmt.Errorf("error inserting into bboltdb cannot assert document to slice of interface")
			}
			objs = docs
		}

		if err := b.client.Update(func(tx *bbolt.Tx) error {

			for _, objToSet := range objs {
				// get _id from create request
				id, ok := objToSet.(map[string]interface{})["_id"]
				if !ok {
					return fmt.Errorf("error creating _id not found in create request")
				}
				// check if specified already exists in database
				count, _, err := b.Read(ctx, col, &model.ReadRequest{
					Find: map[string]interface{}{
						"_id": id,
					},
					Operation: utils.Count,
				})
				if err != nil {
					return fmt.Errorf("error reading existing data - %s", err.Error())
				}
				if count > 0 {
					return fmt.Errorf("error inserting into bboltdb data already exists - %v", count)
				}

				b, err := tx.CreateBucketIfNotExists([]byte(b.bucketName))
				if err != nil {
					return fmt.Errorf("error creating bucket in bboltdb while inserting- %v", err)
				}

				// store value as json string
				value, err := json.Marshal(&objToSet)
				if err != nil {
					return fmt.Errorf("error marshalling while inserting in bboltdb - %v", err)
				}

				// insert document in bucket
				if err = b.Put([]byte(fmt.Sprintf("%s/%s", col, id)), value); err != nil {
					return fmt.Errorf("error inserting in bbolt db - %v", err)
				}
			}
			return nil
		}); err != nil {
			return 0, err
		}
		return int64(len(objs)), nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
