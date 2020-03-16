package userman

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateUserManagement() (*model.SpecObject, error) {
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project"}, &project); err != nil {
		return nil, err
	}
	provider := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Collection Name"}, &provider); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/user-management/provider{provider}",
		Type: "auth-providers",
		Meta: map[string]string{
			"project":  project,
			"provider": provider,
		},
		Spec: map[string]interface{}{
			"enabled": true,
			"id":      "",
			"secret":  "",
		},
	}

	return v, nil
}
