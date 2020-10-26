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

	resourceID, err := filter.DeleteOptions(prefix, tableNames, doesTableNameExist)
	if err != nil {
		return err
	}

	// Delete the db rules from the server
	url := fmt.Sprintf("/v1/config/projects/%s/database/%s/collections/%s/rules", project, dbAlias, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"dbAlias": dbAlias, "col": resourceID}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteDBConfigs(project, prefix string) error {

	objs, err := GetDbConfig(project, "db-config", map[string]string{"dbAlias": "*"})
	if err != nil {
		return err
	}

	doesAliasExist := false
	aliasNames := []string{}
	for _, spec := range objs {
		aliasNames = append(aliasNames, spec.Meta["dbAlias"])
	}

	resourceID, err := filter.DeleteOptions(prefix, aliasNames, doesAliasExist)
	if err != nil {
		return err
	}

	// Delete the db config from the server
	url := fmt.Sprintf("/v1/config/projects/%s/database/%s/config/%s", project, resourceID, "database-config")

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"dbAlias": resourceID}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteDBPreparedQuery(project, dbAlias, prefix string) error {

	objs, err := GetDbPreparedQuery(project, "db-prepared-query", map[string]string{"dbAlias": dbAlias, "id": "*"})
	if err != nil {
		return err
	}

	doesPreparedQueryExist := false
	preparedQueries := []string{}
	for _, spec := range objs {
		preparedQueries = append(preparedQueries, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, preparedQueries, doesPreparedQueryExist)
	if err != nil {
		return err
	}

	// Delete the db prepared query from the server
	url := fmt.Sprintf("/v1/config/projects/%s/database/%s/prepared-queries/%s", project, dbAlias, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"dbAlias": dbAlias, "id": resourceID}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteDBSchemas(project, dbAlias, prefix string) error {

	objs, err := GetDbSchema(project, "db-schema", map[string]string{"dbAlias": dbAlias, "col": "*"})
	if err != nil {
		return err
	}

	doesTableExist := false
	tableNames := []string{}
	for _, spec := range objs {
		tableNames = append(tableNames, spec.Meta["col"])
	}

	resourceID, err := filter.DeleteOptions(prefix, tableNames, doesTableExist)
	if err != nil {
		return err
	}

	// Delete the db prepared query from the server
	url := fmt.Sprintf("/v1/config/projects/%s/database/%s/collections/%s", project, dbAlias, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"dbAlias": dbAlias, "col": resourceID}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
