package bolt

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (b *Bolt) Read(ctx context.Context, project, col string, req *model.ReadRequest) (int64, interface{}, error) {

	switch req.Operation {
	case utils.All:

		var count int64
		results := []interface{}{}
		if err := b.client.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			c := tx.Bucket([]byte(project)).Cursor()

			prefix := []byte(col + "/")
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				count++
				result := map[string]interface{}{}
				if err := json.Unmarshal(v, result); err != nil {
					logrus.Errorf("error un marshalling while reading from bboltdb - %v", err)
					return err
				}
				if utils.Validate(req.Find, result) {
					results = append(results, result)
				}
			}
			return nil
		}); err != nil {
			return 0, nil, nil
		}

		return count, results, nil

	case utils.One:
		finalResult := map[string]interface{}{}
		if err := b.client.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			c := tx.Bucket([]byte(project)).Cursor()
			prefix := []byte(col + "/")
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				result := map[string]interface{}{}
				if err := json.Unmarshal(v, result); err != nil {
					logrus.Errorf("error un marshalling while reading from bboltdb - %v", err)
					return err
				}
				if utils.Validate(req.Find, result) {
					finalResult = result
				}
			}
			return nil
		}); err != nil {
			return 0, nil, nil
		}
		return 1, finalResult, nil

	default:
		return 0, nil, utils.ErrInvalidParams
	}
}
