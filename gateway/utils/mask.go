package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"

	"github.com/spaceuptech/helpers"
)

// HashString hashes a string value and base64 encodes the result
func HashString(stringValue string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(stringValue))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Encrypt encrypts a string and base64 encodes the result
func Encrypt(aesKey []byte, value string) (string, error) {
	encrypted := make([]byte, len(value))
	if err := encryptAESCFB(encrypted, []byte(value), aesKey, aesKey[:aes.BlockSize]); err != nil {
		return "", helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to encrypt value in match encrypt", err, nil)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func encryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}
