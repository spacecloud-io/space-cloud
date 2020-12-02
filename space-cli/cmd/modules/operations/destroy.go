package operations

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Destroy cleans the environment which has been setup. It removes the containers, secrets & host file
func Destroy() error {

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
		Message: fmt.Sprintf("Space cloud cluster with id (%s) will be destoryed, Do you want to continue", charList[0].Name),
	}
	if err := survey.AskOne(prompt, &isOk); err != nil {
		return err
	}
	if !isOk {
		return nil
	}

	// Delete all projects
	objects, err := project.GetProjectConfig("*", "projects", nil)
	if err != nil {
		return err
	}
	for _, object := range objects {
		projectID, ok := object.Meta["project"]
		if !ok {
			continue
		}

		if err := project.DeleteProject(projectID); err != nil {
			return err
		}
	}

	if err := utils.HelmUninstall(charList[0].Name); err != nil {
		return err
	}

	if err := utils.RemoveAccount(charList[0].Name); err != nil {
		return err
	}
	utils.LogInfo("Space cloud cluster has been destroyed successfully 😢")
	return nil
}
