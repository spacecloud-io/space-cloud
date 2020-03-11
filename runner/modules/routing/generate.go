package routing

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cloud/runner/model"
)

func generateServiceRouting() (*model.SpecObject, error) {

	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter target"}, &project); err != nil {
		return nil, err
	}

	source := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter source"}, &source); err != nil {
		return nil, err
	}

	target := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter target"}, &target); err != nil {
		return nil, err
	}

	sourceType := ""
	if err := survey.AskOne(&survey.Select{Message: "Select source type ", Options: []string{"version", "External"}}, &sourceType); err != nil {
		return nil, err
	}
	host := ""
	port := ""
	version := ""
	switch sourceType {
	case "External":

		if err := survey.AskOne(&survey.Input{Message: "Enter host", Default: "serviceID.projectID.svc.cluster.local"}, &host); err != nil {
			return nil, err
		}

	case "version":
		if err := survey.AskOne(&survey.Input{Message: "Enter version", Default: "v1"}, &version); err != nil {
			return nil, err
		}
	}

	if err := survey.AskOne(&survey.Input{Message: "Enter port", Default: "8080"}, &port); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/routing",
		Type: "service-routing",
		Meta: map[string]string{
			"project": project,
		},
		Spec: map[string]interface{}{
			"source":  source,
			"target":  target,
			"type":    sourceType,
			"host":    host,
			"version": version,
			"port":    port,
		},
	}

	return v, nil
}
