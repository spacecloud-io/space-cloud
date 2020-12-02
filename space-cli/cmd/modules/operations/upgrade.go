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

	charList, err := utils.HelmList(model.HelmSpaceCloudNamespace)
	if err != nil {
		return err
	}
	if len(charList) < 1 {
		utils.LogInfo("space cloud cluster not found, setup a new cluster using the setup command")
		return nil
	}

	isOk := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Space cloud cluster with id (%s) will be upgraded, Do you want to continue", charList[0].Name),
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

	_, err = utils.HelmUpgrade(charList[0].Name, chartLocation, model.HelmSpaceCloudChartDownloadURL, "", valuesFileObj)
	if err != nil {
		return err
	}

	fmt.Println()
	utils.LogInfo(fmt.Sprintf("Space Cloud (cluster id: \"%s\") has been successfully upgraded! 👍", charList[0].Name))
	return nil
}
