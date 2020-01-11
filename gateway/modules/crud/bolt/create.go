package bolt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Create inserts a document (or multiple when op is "all") into the database
func (b *Bolt) Create(ctx context.Context, project, col string, req *model.CreateRequest) (int64, error) {
	// Create a collection object
	// todo check if exists
	switch req.Operation {
	case utils.One:
		// Insert single document
		if err := b.client.Update(func(tx *bolt.Tx) error {
			count, _, err := b.Read(ctx, project, col, &model.ReadRequest{
				Find:      req.Document.(map[string]interface{}),
				Operation: utils.All,
			})
			if count != 0 || err != nil {
				logrus.Error("error inserting into bboltdb data already exists - %v", err)
			}

			b, err := tx.CreateBucketIfNotExists([]byte(project))
			if err != nil {
				logrus.Errorf("error creating bucket in bboltdb while inserting- %v", err)
				return err
			}
			id, ok := req.Document.(map[string]interface{})["_id"]
			if ok {
				value, err := json.Marshal(req.Document)
				if err != nil {
					logrus.Errorf("error marshalling while inserting in bboltdb - %v", err)
					return err
				}

				if err = b.Put([]byte(fmt.Sprintf("%s/%s", col, id)), value); err != nil {
					return err
				}
			}
			logrus.Errorf("error bboltdb unable to find _id while inserting - %v", err)
			return err
		}); err != nil {
			return 0, nil
		}
		return 1, nil

	case utils.All:
		// Insert multiple documents
		objs, ok := req.Document.([]interface{})
		if !ok {
			return 0, utils.ErrInvalidParams
		}

		if err := b.client.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(project))
			if err != nil {
				logrus.Errorf("error creating bucket in bboltdb - %v", err)
				return err
			}
			for _, obj := range objs {
				count, _, err := b.Read(ctx, project, col, &model.ReadRequest{
					Find:      req.Document.(map[string]interface{}),
					Operation: utils.All,
				})
				if count != 0 || err != nil {
					logrus.Error("error inserting into bboltdb data already exists - %v", err)
				}
				id, ok := obj.(map[string]interface{})["_id"]
				if !ok {
					logrus.Errorf("error bboltdb unable to find _id while multiple inserts - %v", err)
					return err
				}
				value, err := json.Marshal(obj)
				if err != nil {
					logrus.Errorf("error marshalling while inserting in bboltdb - %v", err)
					return err
				}

				if err = bucket.Put([]byte(fmt.Sprintf("%s/%s", col, id)), value); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return 0, nil
		}
		return int64(len(objs)), nil

	default:
		return 0, utils.ErrInvalidParams
	}
}
