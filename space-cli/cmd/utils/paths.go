package utils

import (
	"fmt"
	"os"
	"runtime"
)

// GetSpaceCloudDirectory gets the root space cloud directory
func GetSpaceCloudDirectory() string {
	return fmt.Sprintf("%s/.space-cloud", getHomeDirectory())
}

// GetSpaceCloudHostsFilePath returns the path of the hosts files used in space cloud
func GetSpaceCloudHostsFilePath(id string) string {
	if id == "default" {
		return fmt.Sprintf("%s/hosts", GetSpaceCloudDirectory())
	}
	return fmt.Sprintf("%s/%s/hosts", GetSpaceCloudDirectory(), id)
}

// GetSpaceCloudRoutingConfigPath returns the path of the file storing the service routing config
func GetSpaceCloudRoutingConfigPath() string {
	return fmt.Sprintf("%s/routing-config.json", GetSpaceCloudDirectory())
}

// GetSpaceCloudConfigFilePath returns the path of the file storing the config
func GetSpaceCloudConfigFilePath(id string) string {
	if id == "default" {
		return fmt.Sprintf("%s/config.yaml", GetSpaceCloudDirectory())
	}
	return fmt.Sprintf("%s/%s/config.yaml", GetSpaceCloudDirectory(), id)
}

// GetSecretsDir returns the path of the directory storing all the secrets
func GetSecretsDir(id string) string {
	if id == "default" {
		return fmt.Sprintf("%s/secrets", GetSpaceCloudDirectory())
	}
	return fmt.Sprintf("%s/%s/secrets", GetSpaceCloudDirectory(), id)
}

// GetTempSecretsDir gets the path of the directory storing all the temp secrets
func GetTempSecretsDir(id string) string {
	if id == "default" {
		return fmt.Sprintf("%s/secrets/temp-secrets", GetSpaceCloudDirectory())
	}
	return fmt.Sprintf("%s/%s/secrets/temp-secrets", GetSpaceCloudDirectory(), id)
}

func getHomeDirectory() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// GetAccountConfigPath get the path to account config yaml file
func getAccountConfigPath() string {
	return fmt.Sprintf("%s/accounts.yaml", GetSpaceCloudDirectory())
}
