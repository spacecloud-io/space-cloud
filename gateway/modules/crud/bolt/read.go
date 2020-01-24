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

func (b *Bolt) Read(ctx context.Context, project, col string, req *model.ReadRequest) (int64, interface{}, error) {

	switch req.Operation {
	case utils.All, utils.One:
		var count int64
		results := []interface{}{}
		if err := b.client.View(func(tx *bbolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(project))
			if bucket == nil {
				return nil
			}

			cursor := bucket.Cursor()
			prefix := []byte(col + "/")
			for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
				result := map[string]interface{}{}
				if err := json.Unmarshal(v, &result); err != nil {
					logrus.Errorf("error un marshalling while reading from bboltdb - %v", err)
					return err
				}
				if utils.Validate(req.Find, result) {
					results = append(results, result)
					count++
					if req.Operation == utils.One {
						break
					}
				}
			}
			return nil
		}); err != nil {
			return 0, nil, err
		}
		if req.Operation == utils.One {
			if count == 0 {
				return 0, nil, fmt.Errorf("error reading from bbolt db no match found for specifed find clause")
			}
			if count == 1 {
				return count, results[0], nil
			}
		}

		return count, results, nil
	case utils.Count:
		var count int64
		err := b.client.View(func(tx *bbolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(project))
			if bucket == nil {
				return nil
			}
			// not nil means value exists
			if bucket.Get([]byte(fmt.Sprintf("%s/%s", col, req.Find["_id"]))) != nil {
				count = 1
			}

			return nil
		})
		return count, nil, err

	default:
		return 0, nil, utils.ErrInvalidParams
	}
}
