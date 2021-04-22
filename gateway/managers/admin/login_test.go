package admin

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestManager_Login(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		user string
		pass string
	}
	tests := []struct {
		name            string
		args            args
		integrationArgs []mockArgs
		want            int
		wantErr         bool
	}{
		{
			name: "valid login credentials provided",
			args: args{
				user: "admin",
				pass: "123",
			},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{}},
				},
			},
			want:    http.StatusOK,
			wantErr: false,
		},
		{
			name: "Invalid login credentials provided",
			args: args{
				user: "ADMIN",
				pass: "123456",
			},
			integrationArgs: []mockArgs{
				{
					method:         "InvokeHook",
					args:           []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{}},
				},
			},
			want:    http.StatusUnauthorized,
			wantErr: true,
		},
		{
			name: "integration hijack - success",
			args: args{},
			integrationArgs: []mockArgs{
				{
					method: "InvokeHook",
					args:   []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{
						checkResponse: true,
						result: map[string]interface{}{
							"token": "abc",
						},
						status: 200,
					}},
				},
			},
			want: http.StatusOK,
		},
		{
			name: "integration hijack - hook failure",
			args: args{},
			integrationArgs: []mockArgs{
				{
					method: "InvokeHook",
					args:   []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{
						checkResponse: true,
						err:           "oops",
						status:        500,
					}},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "integration hijack - failure response",
			args: args{},
			integrationArgs: []mockArgs{
				{
					method: "InvokeHook",
					args:   []interface{}{mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{
						checkResponse: true,
						err:           "oops",
						status:        403,
					}},
				},
			},
			want:    http.StatusForbidden,
			wantErr: true,
		},
	}
	m := New("nodeID", "clusterID", true, &config.AdminUser{User: "admin", Pass: "123", Secret: "some-secret"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &mockIntegrationManager{}
			for _, arg := range tt.integrationArgs {
				i.On(arg.method, arg.args...).Return(arg.paramsReturned...)
			}

			m.integrationMan = i

			got, _, err := m.Login(context.Background(), tt.args.user, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}
