package operations

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Destroy cleans the environment which has been setup. It removes the containers, secrets & host file
func Destroy() error {

	account, _, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return utils.LogError("Ensure cluster is up and running & space cloud is accessible outside the cluster", err)
	}

	isOk := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Space cloud cluster with id (%s) will be destoryed, Do you want to continue", account.ID),
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

	if err := utils.HelmUninstall(account.ID); err != nil {
		return err
	}

	if err := utils.RemoveAccount(account.ID); err != nil {
		return err
	}
	utils.LogInfo("Space cloud cluster has been destroyed successfully 😢")
	return nil
}
