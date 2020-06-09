package utils

import (
	"fmt"
	"net"

	"github.com/AlecAivazis/survey/v2"
)

// CheckPortAvailability checks if specified port is available on local machine
func CheckPortAvailability(port, s string) (string, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		LogInfo(fmt.Sprintf("The port %s is current busy", port))
		if err := survey.AskOne(&survey.Input{Message: fmt.Sprintf("Enter %s port", s)}, &port); err != nil {
			return "", LogError("error getting port", err)
		}
		if port == "" {
			return "", LogError("Invalid port", err)
		}
		return CheckPortAvailability(port, s)
	}
	_ = ln.Close()
	return port, nil
}

// RemoveAccount removes account from accounts file
func RemoveAccount(id string) error {
	credential, err := GetCredentials()
	if err != nil {
		return err
	}

	index := 0
	for i, v := range credential.Accounts {
		if v.ID == id {
			index = i
			credential.SelectedAccount = ""
		}
	}
	copy(credential.Accounts[index:], credential.Accounts[index+1:])
	credential.Accounts[len(credential.Accounts)-1] = nil
	credential.Accounts = credential.Accounts[:len(credential.Accounts)-1]

	if err := GenerateAccountsFile(credential); err != nil {
		return err
	}

	return err
}

// GetNetworkName provides network name of particular cluster
func GetNetworkName(id string) string {
	if id == "default" {
		return "space-cloud"
	}
	return fmt.Sprintf("space-cloud-%s", id)
}

// GetScContainers provides name for space-cloud containers
func GetScContainers(clusterID, name string) string {
	if clusterID == "default" {
		return fmt.Sprintf("space-cloud-%s", name)
	}
	return fmt.Sprintf("space-cloud-%s-%s", clusterID, name)
}

// GetDatabaseContainerName provides name for database container
func GetDatabaseContainerName(id, alias string) string {
	if id == "default" {
		return fmt.Sprintf("space-cloud--addon--db--%s", alias)
	}
	return fmt.Sprintf("space-cloud-%s--addon--db--%s", id, alias)
}

// GetRegistryContainerName provides name for registry container
func GetRegistryContainerName(id string) string {
	if id == "default" {
		return "space-cloud--addon--registry"
	}
	return fmt.Sprintf("space-cloud-%s--addon--registry", id)
}
