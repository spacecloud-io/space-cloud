package utils

import (
	"os"
	"testing"
)

func Test_getHomeDirectory(t *testing.T) {

	if err := os.Setenv("HOME", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEPATH", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEDRIVE", ""); err != nil {
		return
	}
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: "/home",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHomeDirectory(); got != tt.want {
				t.Errorf("getHomeDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpaceCloudHostsFilePath(t *testing.T) {

	if err := os.Setenv("HOME", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEPATH", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEDRIVE", ""); err != nil {
		return
	}
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: "default",
			want: "/home/.space-cloud/hosts",
		},
		{
			name: "test1",
			args: "name",
			want: "/home/.space-cloud/name/hosts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudHostsFilePath(tt.args); got != tt.want {
				t.Errorf("GetSpaceCloudHostsFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpaceCloudRoutingConfigPath(t *testing.T) {

	if err := os.Setenv("HOME", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEPATH", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEDRIVE", ""); err != nil {
		return
	}
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: "default",
			want: "/home/.space-cloud/routing-config.json",
		},
		{
			name: "test1",
			args: "name",
			want: "/home/.space-cloud/name/routing-config.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudRoutingConfigPath(tt.args); got != tt.want {
				t.Errorf("GetSpaceCloudRoutingConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSecretsDir(t *testing.T) {

	if err := os.Setenv("HOME", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEPATH", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEDRIVE", ""); err != nil {
		return
	}
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: "default",
			want: "/home/.space-cloud/secrets",
		},
		{
			name: "test1",
			args: "name",
			want: "/home/.space-cloud/name/secrets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSecretsDir(tt.args); got != tt.want {
				t.Errorf("GetSecretsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTempSecretsDir(t *testing.T) {

	if err := os.Setenv("HOME", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEPATH", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEDRIVE", ""); err != nil {
		return
	}
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: "name",
			want: "/home/.space-cloud/name/secrets/temp-secrets",
		},
		{
			name: "test1",
			args: "default",
			want: "/home/.space-cloud/secrets/temp-secrets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTempSecretsDir(tt.args); got != tt.want {
				t.Errorf("GetTempSecretsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpaceCloudConfigFilePath(t *testing.T) {

	if err := os.Setenv("HOME", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEPATH", "/home"); err != nil {
		return
	}
	if err := os.Setenv("HOMEDRIVE", ""); err != nil {
		return
	}
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: "default",
			want: "/home/.space-cloud/config.yaml",
		},
		{
			name: "test1",
			args: "name",
			want: "/home/.space-cloud/name/config.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudConfigFilePath(tt.args); got != tt.want {
				t.Errorf("GetSpaceCloudConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpaceCloudClusterDirectory(t *testing.T) {
	type args struct {
		clusterID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{clusterID: "name"},
			want: "/home/.space-cloud/name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudClusterDirectory(tt.args.clusterID); got != tt.want {
				t.Errorf("GetSpaceCloudClusterDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
