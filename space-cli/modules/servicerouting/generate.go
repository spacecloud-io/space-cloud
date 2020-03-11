package servicerouting

import (
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateServiceRouting() (*model.SpecObject, error) {

	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter url"}, &project); err != nil {
		return nil, err
	}

	hosts := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter hosts by comma seperated value: "}, &hosts); err != nil {
		return nil, err
	}
	host := strings.Split(hosts, ",")

	h := make(map[string]interface{})
	for k, v := range host {
		h[strconv.Itoa(k)] = v
	}

	methods := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter methods by comma seperated value: "}, &methods); err != nil {
		return nil, err
	}
	method := strings.Split(methods, ",")

	m := make(map[string]interface{})
	for k, v := range method {
		h[strconv.Itoa(k)] = v
	}

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

	targets := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter targets by comma seperated value: "}, &targets); err != nil {
		return nil, err
	}
	target := strings.Split(targets, ",")

	t := make(map[string]interface{})
	for k, v := range target {
		t[strconv.Itoa(k)] = v
	}

	port := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter port", Default: "8080"}, &port); err != nil {
		return nil, err
	}

	scheme := ""
	if err := survey.AskOne(&survey.Select{Message: "Enter scheme", Options: []string{"HTTP", "HTTPS"}}, &scheme); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/routing",
		Type: "service-routing",
		Meta: map[string]string{
			"project": project,
		},
		Spec: map[string]interface{}{
			"host":       h,
			"method":     m,
			"type":       routingType,
			"url":        url,
			"rewriteURL": rewriteURL,
			"target":     t,
			"port":       port,
			"scheme":     scheme,
			"weight":     "100",
		},
	}

	return v, nil
}
