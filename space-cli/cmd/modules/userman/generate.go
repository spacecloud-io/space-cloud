package userman

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cli/cmd/model"
)

func generateUserManagement() (*model.SpecObject, error) {
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project"}, &project); err != nil {
		return nil, err
	}
	provider := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Provider Name"}, &provider); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/user-management/provider/{id}",
		Type: "auth-providers",
		Meta: map[string]string{
			"project": project,
			"id":      provider,
		},
		Spec: map[string]interface{}{
			"enabled": true,
			"secret":  "",
		},
	}

	return v, nil
}
