package bolt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spaceuptech/helpers"
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
				return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to insert data into bboltdb cannot assert document to map", nil, nil)
			}
			objs = append(objs, doc)
		} else {
			docs, ok := req.Document.([]interface{})
			if !ok {
				return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to insert data into bboltdb cannot assert document to slice of interface", nil, nil)
			}
			objs = docs
		}

		if err := b.client.Update(func(tx *bbolt.Tx) error {

			for _, objToSet := range objs {
				// get _id from create request
				id, ok := objToSet.(map[string]interface{})["_id"]
				if !ok {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to insert data _id not found in create request", nil, nil)
				}
				// check if specified already exists in database
				count, _, _, err := b.Read(ctx, col, &model.ReadRequest{
					Find: map[string]interface{}{
						"_id": id,
					},
					Operation: utils.Count,
				})
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to read existing data", err, nil)
				}
				if count > 0 {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to insert data already exists", nil, nil)
				}

				b, err := tx.CreateBucketIfNotExists([]byte(b.bucketName))
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error creating bucket in bboltdb while inserting- %v", err), nil, nil)
				}

				// store value as json string
				value, err := json.Marshal(&objToSet)
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error marshalling while inserting in bboltdb - %v", err), nil, nil)
				}

				// insert document in bucket
				if err = b.Put([]byte(fmt.Sprintf("%s/%s", col, id)), value); err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error inserting in bbolt db - %v", err), nil, nil)
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
