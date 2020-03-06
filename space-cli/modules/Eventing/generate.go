package eventing

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateEventingRule() (*model.SpecObject, error) {
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &project); err != nil {
		return nil, err
	}
	ruleType := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter rule type"}, &ruleType); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/rules/{type}",
		Type: "eventing-rule",
		Meta: map[string]string{
			"type":    ruleType,
			"project": project,
		},
		Spec: map[string]interface{}{
			"rules": map[string]interface{}{
				"create": map[string]interface{}{
					"rule": "allow",
				},
				"delete": map[string]interface{}{
					"rule": "allow",
				},
				"read": map[string]interface{}{
					"rule": "allow",
				},
				"update": map[string]interface{}{
					"rule": "allow",
				},
			},
		},
	}

	return v, nil
}

func generateEventingSchema() (*model.SpecObject, error) {
	project := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &project); err != nil {
		return nil, err
	}
	ruleType := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter rule type"}, &ruleType); err != nil {
		return nil, err
	}
	schema := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Schema"}, &schema); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/schema/{type}",
		Type: "eventing-schema",
		Meta: map[string]string{
			"type":    ruleType,
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
	if err := survey.AskOne(&survey.Input{Message: "Enter project"}, &project); err != nil {
		return nil, err
	}
	dbAlias := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter DB Alias", Default: dbAlias}, &dbAlias); err != nil {
		return nil, err
	}
	collection := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter colection"}, &collection); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/eventing/config",
		Type: "eventing-config",
		Meta: map[string]string{
			"project": project,
		},
		Spec: map[string]interface{}{
			"dbAlias": dbAlias,
			"col":     collection,
			"enabled": true,
		},
	}

	return v, nil
}
