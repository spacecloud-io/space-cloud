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
func GetSpaceCloudHostsFilePath() string {
	return fmt.Sprintf("%s/hosts", GetSpaceCloudDirectory())
}

// GetSpaceCloudRoutingConfigPath returns the path of the file storing the service routing config
func GetSpaceCloudRoutingConfigPath() string {
	return fmt.Sprintf("%s/routing-config.json", GetSpaceCloudDirectory())
}

// GetSpaceCloudConfigFilePath returns the path of the file storing the config
func GetSpaceCloudConfigFilePath() string {
	return fmt.Sprintf("%s/config.yaml", GetSpaceCloudDirectory())
}

// GetSecretsDir returns the path of the directory storing all the secrets
func GetSecretsDir() string {
	return fmt.Sprintf("%s/secrets", GetSpaceCloudDirectory())
}

// GetTempSecretsDir gets the path of the directory storing all the temp secrets
func GetTempSecretsDir() string {
	return fmt.Sprintf("%s/secrets/temp-secrets", GetSpaceCloudDirectory())
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

func getAccountConfigPath() string {
	return fmt.Sprintf("%s/accounts.yaml", GetSpaceCloudDirectory())
}
