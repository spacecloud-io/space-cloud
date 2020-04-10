package admin

import (
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_LoadEnv(t *testing.T) {
	type fields struct {
		config    *config.Admin
		quotas    model.UsageQuotas
		user      *config.AdminUser
		isProd    bool
		clusterID string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "valid config",
			fields: fields{isProd: true},
			want:   true,
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
			if got := m.LoadEnv(); got != tt.want {
				t.Errorf("LoadEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	m := New("clusterID", &config.AdminUser{User: "admin", Pass: "123", Secret: "some-secret"})
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
			args:   args{admin: &config.Admin{ClusterID: "clusterID", ClusterKey: "clusterKey", Version: 1}},
			fields: fields{},
			want:   &Manager{config: &config.Admin{ClusterID: "clusterID", ClusterKey: "clusterKey", Version: 1}},
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
			m.SetConfig(tt.args.admin)
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("SetConfig() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestManager_SetEnv(t *testing.T) {
	type fields struct {
		config    *config.Admin
		quotas    model.UsageQuotas
		user      *config.AdminUser
		isProd    bool
		clusterID string
	}
	type args struct {
		isProd bool
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   *Manager
	}{
		{
			name:   "valid config provided",
			args:   args{isProd: true},
			fields: fields{},
			want:   &Manager{isProd: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				config:    tt.fields.config,
				quotas:    tt.fields.quotas,
				user:      tt.fields.user,
				isProd:    tt.fields.isProd,
				clusterID: tt.fields.clusterID,
			}
			m.SetEnv(true)
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("SetEnv() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		clusterID     string
		adminUserInfo *config.AdminUser
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "valid config provided",
			args: args{
				clusterID: "clusterID",
				adminUserInfo: &config.AdminUser{
					User:   "admin",
					Pass:   "123",
					Secret: "some-secret",
				},
			},
			want: &Manager{clusterID: "clusterID", user: &config.AdminUser{User: "admin", Pass: "123", Secret: "some-secret"}, quotas: model.UsageQuotas{MaxProjects: 1, MaxDatabases: 1}, config: &config.Admin{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.clusterID, tt.args.adminUserInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
