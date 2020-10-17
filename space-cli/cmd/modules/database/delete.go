package database

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteDBRules(project, dbAlias, prefix string) error {

	objs, err := GetDbRule(project, "db-rule", map[string]string{"dbAlias": dbAlias, "col": "*"})
	if err != nil {
		return err
	}

	doesTableNameExist := false
	tableNames := []string{}
	for _, spec := range objs {
		tableNames = append(tableNames, spec.Meta["col"])
	}

	prefix, err = filter.DeleteOptions(project, prefix, tableNames, doesTableNameExist)
	if err != nil {
		return err
	}

	// Delete the db rules from the server
	url := fmt.Sprintf("/v1/config/projects/%s/database/%s/collections/%s/rules", project, dbAlias, prefix)

	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"dbAlias": dbAlias, "col": prefix}, payload); err != nil {
		return err
	}

	return nil
}
