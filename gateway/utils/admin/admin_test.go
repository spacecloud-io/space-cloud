package admin

import (
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
		want1   string
		wantErr bool
	}{
		{
			name: "valid login credentials provided",
			args: args{
				user: "admin",
				pass: "123",
			},
			want:    http.StatusOK,
			want1:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4",
			wantErr: false,
		},
		{
			name: "Invalid login credentials provided",
			args: args{
				user: "ADMIN",
				pass: "123456",
			},
			want:    http.StatusUnauthorized,
			want1:   "",
			wantErr: true,
		},
	}
	m := New("nodeID", "clusterID", true, &config.AdminUser{User: "admin", Pass: "123", Secret: "some-secret"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := m.Login(tt.args.user, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Login() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
