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

// GetSpaceCloudClusterDirectory gets the root space cloud directory
func GetSpaceCloudClusterDirectory(clusterID string) string {
	return fmt.Sprintf("%s/.space-cloud/%s", getHomeDirectory(), clusterID)
}

// GetSpaceCloudHostsFilePath returns the path of the hosts files used in space cloud
func GetSpaceCloudHostsFilePath(id string) string {
	if id == "default" {
		return fmt.Sprintf("%s/hosts", GetSpaceCloudDirectory())
	}
	return fmt.Sprintf("%s/%s/hosts", GetSpaceCloudDirectory(), id)
}

// GetSpaceCloudRoutingConfigPath returns the path of the file storing the service routing config
func GetSpaceCloudRoutingConfigPath(id string) string {
	if id == "default" {
		return fmt.Sprintf("%s/routing-config.json", GetSpaceCloudDirectory())
	}
	return fmt.Sprintf("%s/%s/routing-config.json", GetSpaceCloudDirectory(), id)
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

// GetMountHostsFilePath returns the path of the hosts files to be mounted in in space cloud
func GetMountHostsFilePath(id string) string {
	if runtime.GOOS == "windows" {
		if id == "default" {
			return fmt.Sprintf("%s/.space-cloud/hosts", getWindowsUserDirectory())
		}
		return fmt.Sprintf("%s/.space-cloud/%s/hosts", getWindowsUserDirectory(), id)
	}
	return GetSpaceCloudHostsFilePath(id)
}

// GetMountConfigFilePath returns the path of the config files to be mounted in space cloud
func GetMountConfigFilePath(id string) string {
	if runtime.GOOS == "windows" {
		if id == "default" {
			return fmt.Sprintf("%s/.space-cloud/config.yaml", getWindowsUserDirectory())
		}
		return fmt.Sprintf("%s/.space-cloud/%s/config.yaml", getWindowsUserDirectory(), id)
	}
	return GetSpaceCloudConfigFilePath(id)
}

// GetMountSecretsDir returns the path of the secret dir to be mounted in space cloud
func GetMountSecretsDir(id string) string {
	if runtime.GOOS == "windows" {
		if id == "default" {
			return fmt.Sprintf("%s/.space-cloud/secrets", getWindowsUserDirectory())
		}
		return fmt.Sprintf("%s/.space-cloud/%s/secrets", getWindowsUserDirectory(), id)
	}
	return GetSecretsDir(id)
}

// GetMountTempSecretsDir returns the path of the temp secret dir to be mounted in space cloud
func GetMountTempSecretsDir(id string) string {
	if runtime.GOOS == "windows" {
		if id == "default" {
			return fmt.Sprintf("%s/.space-cloud/secrets/temp-secrets", getWindowsUserDirectory())
		}
		return fmt.Sprintf("%s/.space-cloud/%s/secrets/temp-secrets", getWindowsUserDirectory(), id)
	}
	return GetTempSecretsDir(id)
}

// GetMountRoutingConfigPath returns the path of the routing config to be mounted in space cloud
func GetMountRoutingConfigPath(id string) string {
	if runtime.GOOS == "windows" {
		if id == "default" {
			return fmt.Sprintf("%s/.space-cloud/routing-config.json", getWindowsUserDirectory())
		}
		return fmt.Sprintf("%s/.space-cloud/%s/routing-config.json", getWindowsUserDirectory(), id)
	}
	return GetSpaceCloudRoutingConfigPath(id)
}

// getWindowsUserDirectory gets home directory to support setup on window
func getWindowsUserDirectory() string {
	// eg. HOMEDRIVE = "C:" and HOMEPATH = "\User\username	"
	homeDrive := strings.ToLower(strings.Split(os.Getenv("HOMEDRIVE"), ":")[0])
	homePath := strings.ReplaceAll(os.Getenv("HOMEPATH"), "\\", "/")

	return "/" + homeDrive + homePath
}
