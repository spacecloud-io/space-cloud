package accounts

import (
	"errors"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/utils/file"
)

func Test_deleteAccount(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		accountID string
	}
	tests := []struct {
		name           string
		args           args
		schemaMockArgs []mockArgs
		wantErr        bool
	}{
		{
			name: "could not fetch stored credentials",
			args: args{accountID: "accountID"},
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
			name: "accountID provided doesn't match any existing account",
			args: args{accountID: "accountID"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
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
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "accountID provided matches an existing account but can not update yaml file",
			args: args{accountID: "local-admin"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, errors.New("could not get file info")},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{true},
				},
				{
					method:         "MkdirAll",
					paramsReturned: []interface{}{errors.New("could not make directory")},
				},
			},
			wantErr: true,
		},
		{
			name: "accountID provided matches an existing account and yaml file updated",
			args: args{accountID: "local-admin"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					paramsReturned: []interface{}{nil},
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

			if err := deleteAccount(tt.args.accountID); (err != nil) != tt.wantErr {
				t.Errorf("deleteAccount() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockSchema.AssertExpectations(t)
		})
	}
}
