package ingress

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/model"
)

func generateIngressRouting() (*model.SpecObject, error) {

	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter project"}, &project); err != nil {
		return nil, err
	}

	routeID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter route ID"}, &routeID); err != nil {
		return nil, err
	}

	hosts := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter hosts by comma seperated value: "}, &hosts); err != nil {
		return nil, err
	}
	host := strings.Split(hosts, ",")

	methods := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter methods by comma seperated value: "}, &methods); err != nil {
		return nil, err
	}
	method := strings.Split(methods, ",")

	url := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter url", Default: "/"}, &url); err != nil {
		return nil, err
	}

	rewriteURL := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter rewriteURL"}, &rewriteURL); err != nil {
		return nil, err
	}

	routingType := ""
	if err := survey.AskOne(&survey.Select{Message: "Select routing type", Options: []string{"prefix", "exact"}}, &routingType); err != nil {
		return nil, err
	}
	var target []interface{}
	var t []string

	want := "y"
	for {

		host1 := ""
		if err := survey.AskOne(&survey.Input{Message: "Enter host"}, &host1); err != nil {
			return nil, err
		}

		port := ""
		if err := survey.AskOne(&survey.Input{Message: "Enter port", Default: "8080"}, &port); err != nil {
			return nil, err
		}

		scheme := ""
		if err := survey.AskOne(&survey.Select{Message: "Enter scheme", Options: []string{"HTTP", "HTTPS"}}, &scheme); err != nil {
			return nil, err
		}

		weight := ""
		if err := survey.AskOne(&survey.Input{Message: "Enter weight"}, &weight); err != nil {
			return nil, err
		}

		t = []string{"host:" + host1, "port:" + port, "scheme:" + scheme, "weight:" + weight}
		target = append(target, t)

		if err := survey.AskOne(&survey.Input{Message: "Add another host?(Y/n)", Default: "n"}, &want); err != nil {
			return nil, err
		}
		if strings.ToLower(want) == "n" {
			break
		}

	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/routing/ingress/{id}",
		Type: "ingress-routes",
		Meta: map[string]string{
			"project": project,
			"id":      routeID,
		},
		Spec: map[string]interface{}{
			"source": map[string]interface{}{
				"hosts":      host,
				"method":     method,
				"type":       routingType,
				"url":        url,
				"rewriteURL": rewriteURL,
			},
			"targets": target,
		},
	}

	return v, nil
}
