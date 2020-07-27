package syncman

import (
	"errors"
	"testing"
)

func TestManager_GetRunnerAddr(t *testing.T) {
	tests := []struct {
		name string
		s    *Manager
		want string
	}{
		{
			name: "got runner address",
			s:    &Manager{runnerAddr: "runnerAddress"},
			want: "runnerAddress",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetRunnerAddr(); got != tt.want {
				t.Errorf("Manager.GetRunnerAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetClusterType(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		admin AdminSyncmanInterface
	}
	tests := []struct {
		name          string
		s             *Manager
		args          args
		adminMockArgs []mockArgs
		want          string
		wantErr       bool
	}{
		{
			name: "runner address is empty",
			s:    &Manager{runnerAddr: ""},
			want: "none",
		},
		{
			name: "could not get internal access token",
			s:    &Manager{runnerAddr: "runnerAddr"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("could not get completed signed token")},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAdmin := mockAdminSyncmanInterface{}

			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.args.admin = &mockAdmin

			got, err := tt.s.GetClusterType(tt.args.admin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetClusterType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.GetClusterType() = %v, want %v", got, tt.want)
			}

			mockAdmin.AssertExpectations(t)
		})
	}
}
