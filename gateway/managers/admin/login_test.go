package admin

import (
	"context"
	"net/http"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestManager_Login(t *testing.T) {
	type args struct {
		user string
		pass string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "valid login credentials provided",
			args: args{
				user: "admin",
				pass: "123",
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
			want:    http.StatusUnauthorized,
			wantErr: true,
		},
	}
	m := New("nodeID", "clusterID", true, &config.AdminUser{User: "admin", Pass: "123", Secret: "some-secret"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
