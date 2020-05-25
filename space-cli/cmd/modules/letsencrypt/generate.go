package letsencrypt

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func generateLetsEncryptDomain() (*model.SpecObject, error) {
	whiteListedDomains := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter White Listed Domain by comma seperated value: "}, &whiteListedDomains); err != nil {
		return nil, err
	}

	email := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Email ID: "}, &email); err != nil {
		return nil, err
	}

	whiteListedDomain := strings.Split(strings.TrimSuffix(whiteListedDomains, ","), ",")
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter project: "}, &project); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/letsencrypt/config/{id}",
		Type: "letsencrypt",
		Meta: map[string]string{
			"project": project,
			"id":      "letsencrypt-config",
		},
		Spec: map[string]interface{}{
			"domains": whiteListedDomain,
			"email":   email,
		},
	}

	return v, nil
}
