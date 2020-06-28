package utils

import (
	"testing"
)

func TestGetServiceDomain(t *testing.T) {
	type args struct {
		projectID string
		serviceID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				projectID: "test",
				serviceID: "service",
			},
			want: "service.test.svc.cluster.local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceDomain(tt.args.projectID, tt.args.serviceID); got != tt.want {
				t.Errorf("GetServiceDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInternalServiceDomain(t *testing.T) {
	type args struct {
		projectID string
		serviceID string
		version   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				projectID: "test",
				serviceID: "service",
				version:   "v1",
			},
			want: "service.test-v1.svc.cluster.local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInternalServiceDomain(tt.args.projectID, tt.args.serviceID, tt.args.version); got != tt.want {
				t.Errorf("GetInternalServiceDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
