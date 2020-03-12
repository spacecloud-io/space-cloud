package services

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateService() (*model.SpecObject, error) {
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project ID: "}, &project); err != nil {
		return nil, err
	}
	service := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter service: "}, &service); err != nil {
		return nil, err
	}
	url := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter url: "}, &url); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/services/{id}",
		Type: "remote-services",
		Meta: map[string]string{
			"id":      service,
			"project": project,
		},
		Spec: map[string]interface{}{
			"URL":       url,
			"Endpoints": "{}",
		},
	}

	return v, nil
}
