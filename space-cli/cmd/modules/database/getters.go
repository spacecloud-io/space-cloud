package database

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

// GetDbRule gets database rule
func GetDbRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/rules", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		obj := item.(map[string]interface{})
		table := obj["table"].(string)
		dbAlias := obj["dbAlias"].(string)
		if table == "event_logs" || table == "invocation_logs" {
			continue
		}
		meta := map[string]string{"project": project, "col": table, "dbAlias": dbAlias}

		delete(obj, "col")
		delete(obj, "dbAlias")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules", commandName, meta, obj)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

// GetDbConfig gets database config
func GetDbConfig(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/config", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		obj := item.(map[string]interface{})
		dbAlias := obj["dbAlias"].(string)
		configID := fmt.Sprintf("%s-config", dbAlias)
		meta := map[string]string{"project": project, "dbAlias": dbAlias, "id": configID}

		// Delete the unwanted keys from spec
		delete(obj, "id")
		delete(obj, "dbAlias")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{dbAlias}/config/{id}", commandName, meta, obj)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)

	}

	return objs, nil
}

// GetDbSchema gets database schema
func GetDbSchema(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/schema/mutate", project)

	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		obj := item.(map[string]interface{})
		tableName := obj["col"]
		dbAlias := obj["dbAlias"]
		if tableName == "event_logs" || tableName == "invocation_logs" || tableName == "default" {
			continue
		}
		meta := map[string]string{"project": project, "col": tableName.(string), "dbAlias": dbAlias.(string)}

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate", commandName, meta, map[string]interface{}{"schema": obj["schema"]})
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

// GetDbPreparedQuery gets database prepared query
func GetDbPreparedQuery(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/prepared-queries", project)

	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		obj := item.(map[string]interface{})
		meta := map[string]string{"project": project, "db": obj["db"].(string), "id": obj["id"].(string)}
		delete(obj, "db")
		delete(obj, "id")
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{db}/prepared-queries/{id}", commandName, meta, obj)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
