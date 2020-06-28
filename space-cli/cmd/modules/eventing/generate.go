package eventing

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func generateEventingRule() (*model.SpecObject, error) {
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &project); err != nil {
		return nil, err
	}
	ruleType := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter rule type"}, &ruleType); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/rules/{id}",
		Type: "eventing-rule",
		Meta: map[string]string{
			"id":      ruleType,
			"project": project,
		},
		Spec: map[string]interface{}{
			"rule": "allow",
		},
	}

	return v, nil
}

func generateEventingSchema() (*model.SpecObject, error) {
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &project); err != nil {
		return nil, err
	}
	ruleType := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter rule type"}, &ruleType); err != nil {
		return nil, err
	}
	schema := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Schema"}, &schema); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/schema/{id}",
		Type: "eventing-schema",
		Meta: map[string]string{
			"id":      ruleType,
			"project": project,
		},
		Spec: map[string]interface{}{
			"schema": schema,
		},
	}

	return v, nil
}

func generateEventingConfig() (*model.SpecObject, error) {
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter project"}, &project); err != nil {
		return nil, err
	}

	dbAlias := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter DB Alias"}, &dbAlias); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/config/{id}",
		Type: "eventing-config",
		Meta: map[string]string{
			"project": project,
			"id":      "eventing-config",
		},
		Spec: map[string]interface{}{
			"dbAlias": dbAlias,
			"enabled": true,
		},
	}

	return v, nil
}

func generateEventingTrigger() (*model.SpecObject, error) {
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter project"}, &project); err != nil {
		return nil, err
	}
	triggerName := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "trigger name"}, &triggerName); err != nil {
		return nil, err
	}

	source := ""
	if err := input.Survey.AskOne(&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &source); err != nil {
		return nil, err
	}
	operationType := ""
	var dbAlias string
	col := ""
	options := map[string]interface{}{}
	switch source {
	case "Database":

		if err := input.Survey.AskOne(&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, &operationType); err != nil {
			return nil, err
		}

		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Database Alias "}, &dbAlias); err != nil {
			return nil, err
		}

		if err := input.Survey.AskOne(&survey.Input{Message: "Enter collection/table name"}, &col); err != nil {
			return nil, err
		}
		options = map[string]interface{}{"db": dbAlias, "col": col}
	case "File Storage":
		if err := input.Survey.AskOne(&survey.Select{Message: "Select trigger operation", Options: []string{"FILE_CREATE", "FILE_DELETE"}}, &operationType); err != nil {
			return nil, err
		}
	case "Custom":
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter trigger type"}, &operationType); err != nil {
			return nil, err
		}
	}
	url := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "webhook url"}, &url); err != nil {
		return nil, err
	}
	wantAdvancedSettings := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Do you want advanced settings? (Y / n) ?", Default: "n"}, &wantAdvancedSettings); err != nil {
		return nil, err
	}
	retries := 3
	timeout := 5000

	if strings.ToLower(wantAdvancedSettings) == "y" {

		if err := input.Survey.AskOne(&survey.Input{Message: "Retries count", Default: "3"}, &retries); err != nil {
			return nil, err
		}

		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Timeout", Default: "5000"}, &timeout); err != nil {
			return nil, err
		}

	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/triggers/{id}",
		Type: "eventing-triggers",
		Meta: map[string]string{
			"project": project,
			"id":      triggerName,
		},
		Spec: map[string]interface{}{
			"type":    operationType,
			"url":     url,
			"retries": retries,
			"timeout": timeout,
			"options": options,
		},
	}

	return v, nil
}
