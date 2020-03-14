package letsencrypt

import (
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateLetsEncryptDomain() (*model.SpecObject, error) {
	whiteListedDomains := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter White Listed Domain by comma seperated value: "}, &whiteListedDomains); err != nil {
		return nil, err
	}

	id := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter  id"}, &id); err != nil {
		return nil, err
	}

	whiteListedDomain := strings.Split(whiteListedDomains, ",")
	wld := make(map[string]interface{})
	for k, v := range whiteListedDomain {
		wld[strconv.Itoa(k)] = v
	}
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter project: "}, &project); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/letsencrypt/config/{id}",
		Type: "eventing-rule",
		Meta: map[string]string{
			"project": project,
			"id":      id,
		},
		Spec: map[string]interface{}{
			"white listed domain": wld,
		},
	}

	return v, nil
}
