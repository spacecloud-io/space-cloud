package secrets

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cloud/runner/model"
)

func generateSecrets() (*model.SpecObject, error) {

	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter project name"}, &project); err != nil {
		return nil, err
	}

	name := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter name"}, &name); err != nil {
		return nil, err
	}

	secretID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter secret ID"}, &secretID); err != nil {
		return nil, err
	}

	secretType := ""
	if err := survey.AskOne(&survey.Select{Message: "Enter secret type", Options: []string{"docker", "file", "env"}}, &secretType); err != nil {
		return nil, err
	}

	rootpath := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter rootpath", Default: "/secret/" + secretID}, &rootpath); err != nil {
		return nil, err
	}

	username := ""
	password := ""
	url := ""
	filename := ""
	filedata := ""
	envname := ""
	envdata := ""
	data := make(map[string]interface{})
	switch secretType {
	case "docker":

		if err := survey.AskOne(&survey.Input{Message: "Enter username"}, &username); err != nil {
			return nil, err
		}
		if err := survey.AskOne(&survey.Password{Message: "Enter password"}, &password); err != nil {
			return nil, err
		}
		if err := survey.AskOne(&survey.Input{Message: "Enter URL"}, &url); err != nil {
			return nil, err
		}
	case "file":
		wantToAddMore := "y"

		for {

			if err := survey.AskOne(&survey.Input{Message: "Enter file name"}, &filename); err != nil {
				return nil, err
			}

			if err := survey.AskOne(&survey.Password{Message: "Enter file data"}, &filedata); err != nil {
				return nil, err
			}

			data[filename] = filedata

			if err := survey.AskOne(&survey.Input{Message: "Add another field?(Y/n)", Default: "n"}, &wantToAddMore); err != nil {
				return nil, err
			}
			if strings.ToLower(wantToAddMore) == "n" {
				break
			}
		}
	case "env":
		wantToAddMore := "y"

		for {

			if err := survey.AskOne(&survey.Input{Message: "Enter env name"}, &envname); err != nil {
				return nil, err
			}

			if err := survey.AskOne(&survey.Password{Message: "Enter env data"}, &envdata); err != nil {
				return nil, err
			}

			data[envname] = envdata

			if err := survey.AskOne(&survey.Input{Message: "Add another field?(Y/n)", Default: "n"}, &wantToAddMore); err != nil {
				return nil, err
			}
			if strings.ToLower(wantToAddMore) == "n" {
				break
			}
		}

	}

	v := &model.SpecObject{
		API:  "/v1/runner/{project}/secrets/{name}",
		Type: "apply-service",
		Meta: map[string]string{
			"project": project,
			"name":    name,
		},
		Spec: map[string]interface{}{
			"secretID": secretID,
			"type":     secretType,
			"rootpath": rootpath,
			"data":     data,
		},
	}

	return v, nil
}
