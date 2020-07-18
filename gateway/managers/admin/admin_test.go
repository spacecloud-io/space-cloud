package admin

import (
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
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
	m := New("", "clusterID", false, &config.AdminUser{User: "admin", Pass: "123", Secret: "some-secret"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := m.Login(tt.args.user, tt.args.pass)
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

func TestManager_SetConfig(t *testing.T) {
	type fields struct {
		config    *config.Admin
		quotas    model.UsageQuotas
		user      *config.AdminUser
		isProd    bool
		clusterID string
	}
	type args struct {
		admin *config.Admin
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   *Manager
	}{
		{
			name:   "valid config provided",
			args:   args{admin: &config.Admin{ClusterID: "clusterID", ClusterKey: "clusterKey", License: "1"}},
			fields: fields{},
			want:   &Manager{config: &config.Admin{ClusterID: "clusterID", ClusterKey: "clusterKey", License: "1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				lock:      sync.RWMutex{},
				config:    tt.fields.config,
				quotas:    tt.fields.quotas,
				user:      tt.fields.user,
				isProd:    tt.fields.isProd,
				clusterID: tt.fields.clusterID,
			}
			_ = m.SetConfig(tt.args.admin)
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("SetConfig() = %v, want %v", m, tt.want)
			}
		})
	}
}
