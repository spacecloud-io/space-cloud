package auth

import (
	"crypto/aes"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// PostProcessMethod to do processing on result
func (m *Module) PostProcessMethod(postProcess *PostProcess, result interface{}) error {
	// Gracefully exist if the result is nil
	if result == nil {
		return nil
	}

	// convert to array of interfaces
	var resultArr []interface{}
	switch val := result.(type) {
	case map[string]interface{}:
		resultArr = []interface{}{val} //make an array of interface with val element
	case []interface{}:
		resultArr = val
	default:
		return errors.New("result is of invalid type")
	}

	for _, doc := range resultArr {
		for _, field := range postProcess.postProcessAction {
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
					logrus.Errorln("error loading value in matchEncrypt: ", err)
					return err
				}
				encrypted := make([]byte, len(loadedValue.(string)))
				err1 := encryptAESCFB(encrypted, []byte(loadedValue.(string)), m.aesKey, m.aesKey[:aes.BlockSize])
				if err1 != nil {
					logrus.Errorln("error encrypting value in matchEncrypt: ", err1)
					return err1
				}
				er := utils.StoreValue(field.Field, encrypted, map[string]interface{}{"res": doc})
				if er != nil {
					logrus.Errorln("error storing value in matchEncrypt: ", er)
					return er
				}

			case "decrypt":
				loadedValue, err := utils.LoadValue(field.Field, map[string]interface{}{"res": doc})
				if err != nil {
					logrus.Errorln("error loading value in matchEncrypt: ", err)
					return err
				}
				decrypted := make([]byte, len(loadedValue.(string)))
				err1 := decryptAESCFB(decrypted, []byte(loadedValue.(string)), m.aesKey, m.aesKey[:aes.BlockSize])
				if err1 != nil {
					logrus.Errorln("error decrypting value in matchEncrypt: ", err1)
					return err1
				}
				er := utils.StoreValue(field.Field, decrypted, map[string]interface{}{"res": doc})
				if er != nil {
					logrus.Errorln("error storing value in matchEncrypt: ", er)
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
