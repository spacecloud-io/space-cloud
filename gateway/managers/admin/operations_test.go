package admin

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

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
		wantErr bool
	}{
		{
			name:    "valid case",
			wantErr: false,
		},
	}
	m := New("", "", true, &config.AdminUser{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.GetInternalAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInternalAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
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
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
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
		mockI   []mockArgs
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
			name:   "valid token and no integration",
			fields: fields{isProd: true, config: &config.Admin{}, user: &config.AdminUser{Secret: "some-secret"}},
			args:   args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4"},
			mockI: []mockArgs{
				{
					method:         "HandleConfigAuth",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{checkResponse: false}},
				},
			},
			wantErr: false,
		},
		{
			name:   "valid token with integration",
			fields: fields{isProd: true, config: &config.Admin{}, user: &config.AdminUser{Secret: "some-secret"}},
			args:   args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4"},
			mockI: []mockArgs{
				{
					method:         "HandleConfigAuth",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{checkResponse: true}},
				},
			},
			wantErr: false,
		},
		{
			name:   "valid token with integration error",
			fields: fields{isProd: true, config: &config.Admin{}, user: &config.AdminUser{Secret: "some-secret"}},
			args:   args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkbWluIiwicm9sZSI6ImFkbWluIn0.N4aa9nBNQHsvnWPUfzmKjMG3YD474ChIyOM5FEUuVm4"},
			mockI: []mockArgs{
				{
					method:         "HandleConfigAuth",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{mockIntegrationResponse{checkResponse: true, err: "some-eror"}},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &mockIntegrationManager{}

			for _, v := range tt.mockI {
				i.On(v.method, v.args...).Return(v.paramsReturned...)
			}

			m := &Manager{
				config:         tt.fields.config,
				quotas:         tt.fields.quotas,
				user:           tt.fields.user,
				isProd:         tt.fields.isProd,
				clusterID:      tt.fields.clusterID,
				integrationMan: i,
			}
			if _, err := m.IsTokenValid(context.Background(), tt.args.token, "", "", nil); (err != nil) != tt.wantErr {
				t.Errorf("IsTokenValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !i.AssertExpectations(t) {
				t.Error("Integration expections failed")
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
	m := New("", "", false, &config.AdminUser{Secret: "some-secret"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.RefreshToken(context.Background(), tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				c:       &config.Config{Projects: []*config.Project{{ID: "abc"}}},
				project: &config.Project{ID: "abc"},
			},
			want: true,
		},
		{
			name: "project max projects creation limit not reached",
			args: args{
				c:       &config.Config{Projects: []*config.Project{}},
				project: &config.Project{ID: "abc"},
			},
			want: true,
		},
		{
			name: "project max projects creation limit reached",
			args: args{
				c:       &config.Config{Projects: []*config.Project{{ID: "abc1"}}},
				project: &config.Project{ID: "abc2"},
			},
			want: false,
		},
	}

	m := New("nodeID", "clusterID", true, &config.AdminUser{})
	m.quotas = model.UsageQuotas{MaxProjects: 1, MaxDatabases: 1}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.ValidateProjectSyncOperation(tt.args.c, tt.args.project); got != tt.want {
				t.Errorf("ValidateProjectSyncOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
