package helpers

import (
	"context"
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// PostProcessMethod to do processing on result
func PostProcessMethod(ctx context.Context, aesKey []byte, postProcess *model.PostProcess, result interface{}) error {
	// Gracefully exits if the result is nil
	if result == nil || postProcess == nil {
		return nil
	}

	// convert to array of interfaces
	var resultArr []interface{}
	switch val := result.(type) {
	case map[string]interface{}:
		resultArr = []interface{}{val} // make an array of interface with val element
	case []interface{}:
		resultArr = val
	default:
		return errors.New("result is of invalid type")
	}

	for _, doc := range resultArr {
		for _, field := range postProcess.PostProcessAction {
			// apply Action on all elements
			switch field.Action {
			case "force":
				if err := utils.StoreValue(ctx, field.Field, field.Value, map[string]interface{}{"res": doc}); err != nil {
					return err
				}

			case "remove":
				if err := utils.DeleteValue(ctx, field.Field, map[string]interface{}{"res": doc}); err != nil {
					return err
				}

			case "encrypt":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to load value in post process", err, nil)
				}
				stringValue, ok := loadedValue.(string)
				if !ok {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid data type found", fmt.Errorf("value should be of type string got (%T)", loadedValue), nil)
				}
				encryptedValue, err := utils.Encrypt(aesKey, stringValue)
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to encrypt string in post process", err, map[string]interface{}{"valueToEncrypt": stringValue})
				}
				er := utils.StoreValue(ctx, field.Field, encryptedValue, map[string]interface{}{"res": doc})
				if er != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to store value in post process", er, nil)
				}

			case "decrypt":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to load value in post process", err, map[string]interface{}{"decrypt": true})
				}
				stringValue, ok := loadedValue.(string)
				if !ok {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid data type found", fmt.Errorf("value should be of type string got (%T)", loadedValue), map[string]interface{}{"decrypt": true})
				}
				decodedValue, err := base64.StdEncoding.DecodeString(stringValue)
				if err != nil {
					return err
				}
				decrypted := make([]byte, len(decodedValue))
				err1 := DecryptAESCFB(decrypted, decodedValue, aesKey, aesKey[:aes.BlockSize])
				if err1 != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to decrypt string in post process", err1, map[string]interface{}{"valueToDecrypt": decodedValue})
				}
				er := utils.StoreValue(ctx, field.Field, string(decrypted), map[string]interface{}{"res": doc})
				if er != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to store value in post process", er, map[string]interface{}{"decrypt": true})
				}

			case "hash":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to load value in post process", err, map[string]interface{}{"hash": true})
				}
				stringValue, ok := loadedValue.(string)
				if !ok {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid data type found", fmt.Errorf("value should be of type string got (%T)", loadedValue), map[string]interface{}{"hash": true})
				}
				h := sha256.New()
				_, _ = h.Write([]byte(stringValue))
				hashed := hex.EncodeToString(h.Sum(nil))
				er := utils.StoreValue(ctx, field.Field, hashed, map[string]interface{}{"res": doc})
				if er != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to store value in post process", er, map[string]interface{}{"hash": true})

				}

			default:
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid action (%s) received in post processing read op", field.Action), nil, nil)
			}
		}
	}
	return nil
}
