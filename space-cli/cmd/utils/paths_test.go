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
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: "/home/.space-cloud/hosts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudHostsFilePath(); got != tt.want {
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
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: "/home/.space-cloud/routing-config.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudRoutingConfigPath(); got != tt.want {
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
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: "/home/.space-cloud/secrets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSecretsDir(); got != tt.want {
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
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: "/home/.space-cloud/secrets/temp-secrets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTempSecretsDir(); got != tt.want {
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
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: "/home/.space-cloud/config.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpaceCloudConfigFilePath(); got != tt.want {
				t.Errorf("GetSpaceCloudConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
