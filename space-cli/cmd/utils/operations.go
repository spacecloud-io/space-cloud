package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
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
func RemoveAccount(clusterName string) error {
	credential, err := GetCredentials()
	if err != nil {
		return err
	}

	index := 0
	isFound := false
	for i, v := range credential.Accounts {
		accountName := strings.Split(v.ID, "--")[0]
		if accountName == clusterName {
			isFound = true
			index = i
			credential.SelectedAccount = ""
		}
	}
	if !isFound {
		return nil
	}
	copy(credential.Accounts[index:], credential.Accounts[index+1:])
	credential.Accounts[len(credential.Accounts)-1] = nil
	credential.Accounts = credential.Accounts[:len(credential.Accounts)-1]

	if err := GenerateAccountsFile(credential); err != nil {
		return err
	}

	return nil
}

// ChangeSelectedAccount change selected account according to cluster name provided
func ChangeSelectedAccount(clusterName string) error {
	credential, err := GetCredentials()
	if err != nil {
		return err
	}

	credential.SelectedAccount = ""
	for _, v := range credential.Accounts {
		// With version (v0.19.0) account id has a clusterName prefix separated by -- (default--someId)
		clusterNameCumAccountID := strings.Split(v.ID, "--")[0]
		if clusterNameCumAccountID == "" {
			// this condition occurs when space cli is logged in with a remote server
			continue
		}
		if clusterNameCumAccountID == clusterName {
			credential.SelectedAccount = v.ID
			break
		}
		// This is for compatibility with version (v0.18.0)
		if clusterNameCumAccountID == v.ID {
			credential.SelectedAccount = v.ID
			break
		}
	}

	if credential.SelectedAccount == "" {
		return fmt.Errorf("no account found in account.yaml")
	}

	if err := GenerateAccountsFile(credential); err != nil {
		return err
	}

	return nil
}

//GetSCImageName get the sc image name and add the image prefix when required
func GetSCImageName(imagePrefix, version string, t model.ImageType) string {
	if imagePrefix != "" && !strings.HasSuffix(imagePrefix,"/")  {
		imagePrefix += "/"
	}
	switch t {
	case model.ImageRunner:
		return fmt.Sprintf("%s%s:%s", imagePrefix, "spaceuptech/runner", version)
	case model.ImageGateway:
		return fmt.Sprintf("%s%s:%s", imagePrefix, "spaceuptech/gateway", version)
	default:
		LogInfo("Invalid image type provided for getting sc image name with image prefix")
		return ""
	}
}

// GetNetworkName provides network name of particular cluster
func GetNetworkName(clusterName string) string {
	if clusterName == "default" {
		return "space-cloud"
	}
	return fmt.Sprintf("space-cloud-%s", clusterName)
}

// GetScContainers provides name for space-cloud containers
func GetScContainers(clusterName, name string) string {
	if clusterName == "default" {
		return fmt.Sprintf("space-cloud-%s", name)
	}
	return fmt.Sprintf("space-cloud-%s-%s", clusterName, name)
}

// GetDatabaseContainerName provides name for database container
func GetDatabaseContainerName(id, alias string) string {
	if id == "default" {
		return fmt.Sprintf("space-cloud--addon--db--%s", alias)
	}
	return fmt.Sprintf("space-cloud-%s--addon--db--%s", id, alias)
}

// GetRegistryContainerName provides name for registry container
func GetRegistryContainerName(clusterName string) string {
	if clusterName == "default" {
		return "space-cloud--addon--registry"
	}
	return fmt.Sprintf("space-cloud-%s--addon--registry", clusterName)
}
