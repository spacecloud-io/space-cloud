package database

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteDBRules(project, dbAlias, prefix string) error {

	objs, err := GetDbRule(project, "db-rule", map[string]string{"dbAlias": "*", "col": "*"})
	if err != nil {
		return err
	}

	doesTableNameExist := false
	tableNames := []string{}
	for _, spec := range objs {
		tableNames = append(tableNames, spec.Meta["col"])
	}

	filteredTableNames := []string{}
	for _, tableName := range tableNames {
		if prefix != "" && strings.HasPrefix(strings.ToLower(tableName), strings.ToLower(prefix)) {
			filteredTableNames = append(filteredTableNames, tableName)
			doesTableNameExist = true
		}
	}

	if doesTableNameExist {
		if len(filteredTableNames) == 1 {
			prefix = filteredTableNames[0]
		} else {
			if err := input.Survey.AskOne(&survey.Select{Message: "Choose the table name: ", Options: filteredTableNames, Default: filteredTableNames[0]}, &prefix); err != nil {
				return err
			}
		}
	} else {
		if len(tableNames) == 1 {
			prefix = tableNames[0]
		} else {
			if prefix != "" {
				utils.LogInfo("Warning! No table name found for prefix provided, showing all")
			}
			if err := input.Survey.AskOne(&survey.Select{Message: "Choose the table name: ", Options: tableNames, Default: tableNames[0]}, &prefix); err != nil {
				return err
			}
		}
	}

	// Delete the db rules from the server
	url := fmt.Sprintf("/v1/config/projects/%s/database/%s/collections/%s/rules", project, dbAlias, prefix)

	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"dbAlias": dbAlias, "col": prefix}, payload); err != nil {
		return err
	}

	return nil
}
