package accounts

import (
	"errors"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/utils/file"
)

func Test_listAccounts(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		accountID string
		showKeys  bool
	}
	tests := []struct {
		name           string
		args           args
		schemaMockArgs []mockArgs
		wantErr        bool
	}{
		{
			name: "could not fetch stored credentials",
			args: args{accountID: "accountID", showKeys: false},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte{}, errors.New("could not read yaml file")},
				},
			},
			wantErr: true,
		},
		{
			name: "could fetch stored credentials but has no accounts stored",
			args: args{accountID: "accountID", showKeys: false},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte(""), nil},
				},
			},
		},
		{
			name: "account ID not specified and showKeys false",
			args: args{accountID: "", showKeys: false},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
		},
		{
			name: "account ID not specified and showKeys true",
			args: args{accountID: "", showKeys: true},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
		},
		{
			name: "account ID specified but does not match any existing account",
			args: args{accountID: "accountID", showKeys: false},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
		},
		{
			name: "account ID specified and matches an existing account with showKeys false",
			args: args{accountID: "local-admin", showKeys: false},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
		},
		{
			name: "account ID specified and matches an existing account with showKeys true",
			args: args{accountID: "local-admin", showKeys: true},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema

			if err := listAccounts(tt.args.accountID, tt.args.showKeys); (err != nil) != tt.wantErr {
				t.Errorf("listAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockSchema.AssertExpectations(t)
		})
	}
}
