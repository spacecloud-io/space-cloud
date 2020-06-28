package utils

import (
	"crypto/aes"
	"encoding/base64"
	"reflect"
	"testing"
)

func Test_encryptAESCFB(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "invalid key",
			args:    args{dst: make([]byte, len("username1")), src: []byte("username1"), key: []byte("invalidKey"), iv: []byte("invalidKey123456")[:aes.BlockSize]},
			wantErr: true,
		},
		{
			name: "encryption takes place",
			args: args{dst: make([]byte, len("username1")), src: []byte("username1"), key: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g="), iv: []byte(base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g="))[:aes.BlockSize]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := encryptAESCFB(tt.args.dst, tt.args.src, tt.args.key, tt.args.iv); (err != nil) != tt.wantErr {
				t.Errorf("encryptAESCFB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && reflect.DeepEqual(tt.args.dst, tt.args.src) {
				t.Errorf("encryptAESCFB() encryption did not take place")
			}
		})
	}
}

func base64DecodeString(key string) []byte {
	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	return decodedKey
}
