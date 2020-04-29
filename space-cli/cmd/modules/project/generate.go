package project

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func generateProject() (*model.SpecObject, error) {
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project ID: "}, &project); err != nil {
		return nil, err
	}
	if project == "" {
		_ = utils.LogError("project id cannot be empty", nil)
		return nil, nil
	}
	projectName := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project Name: ", Default: project}, &projectName); err != nil {
		return nil, err
	}

	key := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter AES Key: "}, &key); err != nil {
		return nil, err
	}

	contextTime := 0
	if err := survey.AskOne(&survey.Input{Message: "Enter Graphql Query Timeout: ", Default: "10"}, &contextTime); err != nil {
		return nil, err
	}
	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}",
		Type: "project",
		Meta: map[string]string{
			"project": project,
		},
		Spec: map[string]interface{}{
			"id":                 project,
			"aesKey":             key,
			"name":               projectName,
			"secrets":            []map[string]interface{}{{"isPrimary": true, "secret": ksuid.New().String()}},
			"contextTimeGraphQL": contextTime,
		},
	}

	return v, nil
}
