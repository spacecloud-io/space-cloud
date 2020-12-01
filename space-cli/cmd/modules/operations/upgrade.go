package operations

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Upgrade upgrades existing space cloud cluster
func Upgrade(setValuesFlag, valuesYamlFile, chartLocation string) error {
	_ = utils.CreateDirIfNotExist(utils.GetSpaceCloudDirectory())

	selectedAccount, _, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return err
	}

	isOk := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Space cloud cluster with id (%s) will be upgraded, Do you want to continue", selectedAccount.ID),
	}
	if err := survey.AskOne(prompt, &isOk); err != nil {
		return err
	}
	if !isOk {
		return nil
	}

	valuesFileObj, err := utils.ExtractValuesObj(setValuesFlag, valuesYamlFile)
	if err != nil {
		return err
	}

	// override the clusterId
	value, ok := valuesFileObj["clusterId"]
	if ok {
		return fmt.Errorf("you cannot set a new cluster id (%s) while upgrading an existing cluster, revmove the value & try again", value)
	}

	_, err = utils.HelmUpgrade(selectedAccount.ID, chartLocation, model.HelmSpaceCloudChartDownloadURL, "", valuesFileObj)
	if err != nil {
		return err
	}

	fmt.Println()
	utils.LogInfo(fmt.Sprintf("Space Cloud (cluster id: \"%s\") has been successfully upgraded! 👍", selectedAccount.ID))
	utils.LogInfo(fmt.Sprintf("You can visit mission control at %s/mission-control 💻", selectedAccount.ServerURL))
	utils.LogInfo("Note: The url is only valid if you have done port forwarding using the commnad as per the docs at https://docs.spaceuptech.com/install/kubernetes/minikube/")
	utils.LogInfo("Command => kubectl port-forward -n istio-system deployments/istio-ingressgateway 4122:8080")
	utils.LogInfo("If you have done forwarding on other port, use the login command to change the url")
	utils.LogInfo(fmt.Sprintf("Your login credentials: [username: \"%s\"; key: \"%s\"] 🤫", selectedAccount.UserName, selectedAccount.Key))
	return nil
}
