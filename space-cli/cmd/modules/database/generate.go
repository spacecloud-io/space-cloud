package database

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func generateDBRule() (*model.SpecObject, error) {
	projectID := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &projectID); err != nil {
		return nil, err
	}
	ID := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Collection Name"}, &ID); err != nil {
		return nil, err
	}
	dbAlias := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter DB Alias"}, &dbAlias); err != nil {
		return nil, err
	}
	wantRealtimeEnabled := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Is Realtime enabled (Y / n) ?", Default: "n"}, &wantRealtimeEnabled); err != nil {
		return nil, err
	}

	var isRealTimeEnabled bool
	if strings.ToLower(wantRealtimeEnabled) == "y" {
		isRealTimeEnabled = true
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules",
		Type: "db-rules",
		Meta: map[string]string{
			"dbAlias": dbAlias,
			"col":     ID,
			"project": projectID,
		},
		Spec: map[string]interface{}{
			"isRealtimeEnabled": isRealTimeEnabled,
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

func generateDBConfig() (*model.SpecObject, error) {
	projectID := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &projectID); err != nil {
		return nil, err
	}

	var dbType string
	if err := input.Survey.AskOne(&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &dbType); err != nil {
		return nil, err
	}

	connDefault := ""
	switch dbType {
	case "postgres":

		connDefault = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
	case "sqlserver":

		connDefault = "Data Source=localhost,1433;Initial Catalog=master;User ID=yourID;Password=yourPassword@#;"
	case "embedded":

		connDefault = "Data.db"
	case "mongo":

		connDefault = "mongodb://localhost:27017"
	case "mysql":

		connDefault = "root:my-secret-pw@tcp(localhost:3306)/"
	default:
		return nil, fmt.Errorf("Invalid choice")
	}
	conn := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Database Connection String ", Default: connDefault}, &conn); err != nil {
		return nil, err
	}
	dbAlias := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter DB Alias", Default: dbType}, &dbAlias); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/database/{dbAlias}/config/{id}",
		Type: "db-config",
		Meta: map[string]string{
			"dbAlias": dbAlias,
			"project": projectID,
			"id":      dbAlias + "-config",
		},
		Spec: map[string]interface{}{
			"conn":      conn,
			"enabled":   true,
			"isPrimary": false,
			"type":      dbType,
		},
	}

	return v, nil
}

func generateDBSchema() (*model.SpecObject, error) {
	project := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project "}, &project); err != nil {
		return nil, err
	}
	col := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Collection "}, &col); err != nil {
		return nil, err
	}
	dbAlias := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter DB Alias"}, &dbAlias); err != nil {
		return nil, err
	}
	schema := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Schema"}, &schema); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate",
		Type: "db-schema",
		Meta: map[string]string{
			"dbAlias": dbAlias,
			"project": project,
			"col":     col,
		},
		Spec: map[string]interface{}{
			"schema": schema,
		},
	}

	return v, nil
}
