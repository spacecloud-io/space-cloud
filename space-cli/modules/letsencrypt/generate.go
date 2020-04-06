package letsencrypt

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateLetsEncryptDomain() (*model.SpecObject, error) {
	whiteListedDomains := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter White Listed Domain by comma seperated value: "}, &whiteListedDomains); err != nil {
		return nil, err
	}

	whiteListedDomain := strings.Split(strings.TrimSuffix(whiteListedDomains, ","), ",")
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter project: "}, &project); err != nil {
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
		},
	}

	return v, nil
}
