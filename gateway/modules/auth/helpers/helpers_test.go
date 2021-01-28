package helpers

import (
	"crypto/aes"
	"encoding/base64"
	"reflect"
	"testing"
)

func base64DecodeString(key string) []byte {
	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	return decodedKey
}

func Test_decryptAESCFB(t *testing.T) {

	type args struct {
		dst []byte
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid key",
			args:    args{dst: make([]byte, len("username1")), src: []byte("username1"), key: []byte("invalidKey"), iv: []byte("invalidKey123456")[:aes.BlockSize]},
			wantErr: true,
		},
		{
			name: "decryption takes place",
			args: args{dst: make([]byte, len("username1")), src: []byte{5, 120, 168, 68, 222, 6, 202, 246, 108}, key: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g="), iv: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")[:aes.BlockSize]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecryptAESCFB(tt.args.dst, tt.args.src, tt.args.key, tt.args.iv); (err != nil) != tt.wantErr {
				t.Errorf("decryptAESCFB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && reflect.DeepEqual(tt.args.dst, tt.args.src) {
				t.Errorf("decryptAESCFB() decryption did not take place")
			}
		})
	}
}
