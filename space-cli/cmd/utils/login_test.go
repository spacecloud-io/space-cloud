package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/file"
)

func Test_login(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		selectedAccount *model.Account
	}

	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		args           args
		want           *model.LoginResponse
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "succesfully login",
			schemaMockArgs: []mockArgs{
				{
					method:         "Post",
					args:           []interface{}{},
					paramsReturned: []interface{}{200, nil},
				},
			},
			args: args{
				selectedAccount: &model.Account{
					UserName: "username",
					Key:      "key",
				},
			},
			want:    &model.LoginResponse{},
			wantErr: false,
		},
		{
			name: "statusCode not ok",
			schemaMockArgs: []mockArgs{
				{
					method:         "Post",
					args:           []interface{}{},
					paramsReturned: []interface{}{201, nil},
				},
			},
			args: args{
				selectedAccount: &model.Account{
					UserName: "username",
					Key:      "key",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error making http post request",
			schemaMockArgs: []mockArgs{
				{
					method:         "Post",
					args:           []interface{}{},
					paramsReturned: []interface{}{200, fmt.Errorf("some-error")},
				},
			},
			args: args{
				selectedAccount: &model.Account{},
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

			got, err := login(tt.args.selectedAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("login() = %v, want %v", got, tt.want)
			}
			b, _ := json.MarshalIndent(got, "", " ")
			fmt.Println(string(b))
		})
	}
}

func TestLoginWithSelectedAccount(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		want           *model.Account
		want1          string
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
				{
					method:         "Post",
					args:           []interface{}{},
					paramsReturned: []interface{}{200, nil},
				},
			},
			want: &model.Account{
				ID:        "local-admin",
				UserName:  "local-admin",
				Key:       "81WZUGRTtHbG",
				ServerURL: "http://localhost:4122",
			},
			want1:   "",
			wantErr: false,
		},
		{
			name: "error getting selected account",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), fmt.Errorf("some-error")},
				},
			},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
		{
			name: "error login",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Post",
					args:           []interface{}{},
					paramsReturned: []interface{}{200, fmt.Errorf("some-error")},
				},
			},
			want:    nil,
			want1:   "",
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

			got, got1, err := LoginWithSelectedAccount()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginWithSelectedAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginWithSelectedAccount() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoginWithSelectedAccount() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
