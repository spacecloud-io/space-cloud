package admin

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_GetClusterID(t *testing.T) {
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
		want   string
	}{
		{
			name:   "Valid case",
			fields: fields{clusterID: "clusterID"},
			want:   "clusterID",
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
			if got := m.GetClusterID(); got != tt.want {
				t.Errorf("GetClusterID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetCredentials(t *testing.T) {
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
		want   map[string]interface{}
	}{
		{
			name:   "Valid case",
			fields: fields{user: &config.AdminUser{User: "admin", Pass: "123"}},
			want:   map[string]interface{}{"user": "admin", "pass": "123"},
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
			if got := m.GetCredentials(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetInternalAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid case",
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.S99PtbIiuUlWtntVzjtSugibEPVwZc00jgCzpErgg6Y",
			wantErr: false,
		},
	}
	m := New("", &config.AdminUser{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetInternalAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInternalAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInternalAccessToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetQuotas(t *testing.T) {
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
		want   *model.UsageQuotas
	}{
		{
			name:   "Valid case",
			fields: fields{quotas: model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1}},
			want:   &model.UsageQuotas{MaxProjects: 1, MaxDatabases: 1},
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
			if got := m.GetQuotas(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuotas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_IsTokenValid(t *testing.T) {
	type fields struct {
		config    *config.Admin
		quotas    model.UsageQuotas
		user      *config.AdminUser
		isProd    bool
		clusterID string
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "valid not production mode",
			fields:  fields{isProd: false},
			args:    args{token: "some-token"},
			wantErr: false,
		},
		{
			name:    "valid token",
			fields:  fields{isProd: true, user: &config.AdminUser{Secret: "some-secret"}},
			args:    args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4"},
			wantErr: false,
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
			if err := m.IsTokenValid(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("IsTokenValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_RefreshToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid token",
			args:    args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4"},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4",
			wantErr: false,
		},
	}
	m := New("", &config.AdminUser{Secret: "some-secret"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.RefreshToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RefreshToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_ValidateSyncOperation(t *testing.T) {
	type args struct {
		c       *config.Config
		project *config.Project
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "project already exists",
			args: args{
				c:       &config.Config{Projects: []*config.Project{{ID: "projectID"}}},
				project: &config.Project{ID: "projectID"},
			},
			want: true,
		},
		{
			name: "project max projects creation limit not reached",
			args: args{
				c:       &config.Config{Projects: []*config.Project{}},
				project: &config.Project{ID: "projectID"},
			},
			want: true,
		},
		{
			name: "project max projects creation limit reached",
			args: args{
				c:       &config.Config{Projects: []*config.Project{{ID: "project1"}}},
				project: &config.Project{ID: "project2"},
			},
			want: false,
		},
	}
	m := New("clusterID", &config.AdminUser{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.ValidateSyncOperation(tt.args.c, tt.args.project); got != tt.want {
				t.Errorf("ValidateSyncOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
