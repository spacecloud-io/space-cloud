package remoteservices

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func generateService() (*model.SpecObject, error) {
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID: "}, &project); err != nil {
		return nil, err
	}
	service := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Service Name: "}, &service); err != nil {
		return nil, err
	}
	url := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Service URL: ", Help: "e.g -> http://localhost:8090"}, &url); err != nil {
		return nil, err
	}
	endpoints := []interface{}{}
	want := "y"
	for {
		endpointName := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Endpoint Name: "}, &endpointName); err != nil {
			return nil, err
		}
		methods := ""
		if err := input.Survey.AskOne(&survey.Select{Message: "Select Method: ", Options: []string{"POST", "PUT", "GET", "DELETE"}}, &methods); err != nil {
			return nil, err
		}

		path := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter URL Path:", Default: "/"}, &path); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, map[string]interface{}{endpointName: map[string]interface{}{"method": methods, "path": path}})
		if err := input.Survey.AskOne(&survey.Input{Message: "Add another host?(Y/n)", Default: "n"}, &want); err != nil {
			return nil, err
		}
		if strings.ToLower(want) == "n" {
			break
		}
	}
	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/remote-service/service/{id}",
		Type: "remote-services",
		Meta: map[string]string{
			"id":      service,
			"project": project,
		},
		Spec: map[string]interface{}{
			"url":       url,
			"endpoints": endpoints,
		},
	}

	return v, nil
}
