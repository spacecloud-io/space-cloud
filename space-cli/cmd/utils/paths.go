package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

// GetAccountConfigPath get the path to account config yaml file
func getAccountConfigPath() string {
	return fmt.Sprintf("%s/accounts.yaml", GetSpaceCloudDirectory())
}

// GetMountHostsFilePath returns the path of the hosts files to be mounted in in space cloud
func GetMountHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s/.space-cloud/hosts", getWindowsUserDirectory())
	}
	return GetSpaceCloudHostsFilePath()
}

// GetMountConfigFilePath returns the path of the config files to be mounted in space cloud
func GetMountConfigFilePath() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s/.space-cloud/config.yaml", getWindowsUserDirectory())
	}
	return GetSpaceCloudConfigFilePath()
}

// GetMountSecretsDir returns the path of the secret dir to be mounted in space cloud
func GetMountSecretsDir() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s/.space-cloud/secrets", getWindowsUserDirectory())
	}
	return GetSecretsDir()
}

// GetMountTempSecretsDir returns the path of the temp secret dir to be mounted in space cloud
func GetMountTempSecretsDir() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s/.space-cloud/secrets/temp-secrets", getWindowsUserDirectory())
	}
	return GetTempSecretsDir()
}

// GetMountRoutingConfigPath returns the path of the routing config to be mounted in space cloud
func GetMountRoutingConfigPath() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s/.space-cloud/routing-config.json", getWindowsUserDirectory())
	}
	return GetSpaceCloudRoutingConfigPath()
}

// getWindowsUserDirectory gets home directory to support setup on window
func getWindowsUserDirectory() string {
	// eg. HOMEDRIVE = "C:" and HOMEPATH = "\User\username	"
	homeDrive := strings.ToLower(strings.Split(os.Getenv("HOMEDRIVE"), ":")[0])
	homePath := strings.ReplaceAll(os.Getenv("HOMEPATH"), "\\", "/")

	return "/" + homeDrive + homePath
}
