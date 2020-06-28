package auth

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// PostProcessMethod to do processing on result
func (m *Module) PostProcessMethod(postProcess *model.PostProcess, result interface{}) error {
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
				if err := utils.StoreValue(field.Field, field.Value, map[string]interface{}{"res": doc}); err != nil {
					return err
				}

			case "remove":
				if err := utils.DeleteValue(field.Field, map[string]interface{}{"res": doc}); err != nil {
					return err
				}

			case "encrypt":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					logrus.Errorln("error loading value in postProcessMethod: ", err)
					return err
				}
				stringValue, ok := loadedValue.(string)
				if !ok {
					return fmt.Errorf("Value should be of type string and not %T", loadedValue)
				}
				encryptedValue, err := utils.Encrypt(m.aesKey, stringValue)
				if err != nil {
					return utils.LogError("Unable to encrypt string", "auth", "post-process", err)
				}
				er := utils.StoreValue(field.Field, encryptedValue, map[string]interface{}{"res": doc})
				if er != nil {
					logrus.Errorln("error storing value in postProcessMethod: ", er)
					return er
				}

			case "decrypt":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					logrus.Errorln("error loading value in postProcessMethod: ", err)
					return err
				}
				stringValue, ok := loadedValue.(string)
				if !ok {
					return fmt.Errorf("Value should be of type string and not %T", loadedValue)
				}
				decodedValue, err := base64.StdEncoding.DecodeString(stringValue)
				if err != nil {
					return err
				}
				decrypted := make([]byte, len(decodedValue))
				err1 := decryptAESCFB(decrypted, decodedValue, m.aesKey, m.aesKey[:aes.BlockSize])
				if err1 != nil {
					logrus.Errorln("error decrypting value in postProcessMethod: ", err1)
					return err1
				}
				er := utils.StoreValue(field.Field, string(decrypted), map[string]interface{}{"res": doc})
				if er != nil {
					logrus.Errorln("error storing value in postProcessMethod: ", er)
					return er
				}

			case "hash":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					logrus.Errorln("error loading value in postProcessMethod: ", err)
					return err
				}
				stringValue, ok := loadedValue.(string)
				if !ok {
					return fmt.Errorf("Value should be of type string and not %T", loadedValue)
				}
				h := sha256.New()
				_, _ = h.Write([]byte(stringValue))
				hashed := hex.EncodeToString(h.Sum(nil))
				er := utils.StoreValue(field.Field, hashed, map[string]interface{}{"res": doc})
				if er != nil {
					logrus.Errorln("error storing value in matchHash: ", er)
					return er
				}

			default:
				err := fmt.Errorf("invalid action (%s) received in post processing read op", field.Action)
				return err
			}
		}
	}
	return nil
}
