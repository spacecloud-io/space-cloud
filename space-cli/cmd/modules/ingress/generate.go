package ingress

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func generateIngressRouting() (*model.SpecObject, error) {

	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter project"}, &project); err != nil {
		return nil, err
	}

	hosts := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter hosts by comma separated value: "}, &hosts); err != nil {
		return nil, err
	}
	host := strings.Split(hosts, ",")

	methods := []string{}
	if err := input.Survey.AskOne(&survey.MultiSelect{Message: "Select Methods: ", Options: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE"}}, &methods); err != nil {
		return nil, err
	}

	url := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter url", Default: "/"}, &url); err != nil {
		return nil, err
	}

	rewriteURL := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter rewriteURL"}, &rewriteURL); err != nil {
		return nil, err
	}

	routingType := ""
	if err := input.Survey.AskOne(&survey.Select{Message: "Select routing type", Options: []string{"prefix", "exact"}}, &routingType); err != nil {
		return nil, err
	}
	var target []interface{}
	totalWeight := 0
	want := "y"
	for {

		host1 := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter host"}, &host1); err != nil {
			return nil, err
		}

		port := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter port", Default: "8080"}, &port); err != nil {
			return nil, err
		}

		scheme := ""
		if err := input.Survey.AskOne(&survey.Select{Message: "Enter scheme", Options: []string{"HTTP", "HTTPS"}}, &scheme); err != nil {
			return nil, err
		}

		weight := 0
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter weight"}, &weight); err != nil {
			return nil, err
		}
		t := map[string]interface{}{"host": host1, "port": port, "schema": scheme, "type": "", "weight": weight, "version": ""}
		target = append(target, t)
		totalWeight += weight
		if err := input.Survey.AskOne(&survey.Input{Message: "Add another host?(Y/n)", Default: "n"}, &want); err != nil {
			return nil, err
		}
		if strings.ToLower(want) == "n" {
			break
		}

	}
	if totalWeight != 100 {
		_ = utils.LogError("sum of weights of all targets should be 100", nil)
		return nil, nil
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/routing/ingress/{id}",
		Type: "ingress-routes",
		Meta: map[string]string{
			"project": project,
			"id":      ksuid.New().String(),
		},
		Spec: map[string]interface{}{
			"source": map[string]interface{}{
				"hosts":      host,
				"method":     methods,
				"port":       0,
				"type":       routingType,
				"url":        url,
				"rewriteURL": rewriteURL,
			},
			"targets": target,
		},
	}

	return v, nil
}
