package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/file"
)

func Test_getSelectedAccount(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		want           *model.Account
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "proper",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
			want: &model.Account{
				ID:        "local-admin",
				UserName:  "local-admin",
				Key:       "81WZUGRTtHbG",
				ServerURL: "http://localhost:4122",
			},
			wantErr: false,
		},
		{
			name: "improper file content",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("wrong content"), nil},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error reading file",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("ok content"), fmt.Errorf("error reading file")},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema
			got, err := getSelectedAccount()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSelectedAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSelectedAccount() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestStoreCredentials(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		account *model.Account
	}

	tests := []struct {
		name           string
		args           args
		schemaMockArgs []mockArgs
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "added new account to existing file",
			args: args{&model.Account{
				ID:        "local-admin1",
				UserName:  "local-admin1",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
		},
		{
			name: "error adding new account to existing file",
			args: args{&model.Account{
				ID:        "local-admin1",
				UserName:  "local-admin1",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, fmt.Errorf("some_%s", "error")},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{fmt.Errorf("error writing to existing file")},
				},
			},
			wantErr: true,
		},
		{
			name: "updating account to existing file",
			args: args{&model.Account{
				ID:        "local-admin",
				UserName:  "local-admin1",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
		},

		{
			name: "error updating account to existing file",
			args: args{&model.Account{
				ID:        "local-admin",
				UserName:  "local-admin1",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, fmt.Errorf("some_%s", "error")},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{fmt.Errorf("error writing to existing file")},
				},
			},
			wantErr: true,
		},
		{
			name: "error unmarshalling",
			args: args{&model.Account{
				ID:        "local-admin",
				UserName:  "local-admin1",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("wrong con"), nil},
				},
			},
			wantErr: true,
		},
		{
			name: "added new account to non-existing file",
			args: args{&model.Account{
				ID:        "local-admin",
				UserName:  "local-admin",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte(""), fmt.Errorf("%s file does not exists", "accounts.yaml")},
				},
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, fmt.Errorf("some_%s", "error")},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: false,
		},
		{
			name: "added new account to non-existing file but directory could not be created",
			args: args{&model.Account{
				ID:        "local-admin",
				UserName:  "local-admin",
				Key:       "gibberish",
				ServerURL: "http://localhost:4122",
			}},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte(""), fmt.Errorf("%s file does not exists", "accounts.yaml")},
				},
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, fmt.Errorf("some_%s", "error")},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "MkdirAll",
					args:           []interface{}{},
					paramsReturned: []interface{}{fmt.Errorf("cannot create directory")},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema

			if err := StoreCredentials(tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("StoreCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
