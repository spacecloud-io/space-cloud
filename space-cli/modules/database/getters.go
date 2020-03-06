package database

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

func getDbRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/rules", project)
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["rule"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}
	if value, p := result["rules"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}

	var objs []*model.SpecObject
	for _, item := range array {
		spec := item.(map[string]interface{})
		meta := map[string]string{"projectId": project, "id": spec["id"].(string), "dbAlias": params["dbAlias"]}

		// Delete the unwanted keys from spec
		delete(spec, "id")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/database/{dbAlias}/collections/{id}/rules", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

func getDbConfig(project, commandName string, params map[string]string) (*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/config", project)
	// Get the spec from the server
	result := new(interface{})
	if err := cmd.Get(http.MethodGet, url, params, result); err != nil {
		return nil, err
	}

	// Generating the object
	meta := map[string]string{"projectId": project, "dbAlias": params["dbAlias"]}
	s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/database/{dbAlias}/config", commandName, meta, result)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getDbSchema(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/modify-schema", project)
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["schema"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = params["dbAlias"]
		array = []interface{}{obj}
	}
	if value, p := result["schemas"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}

	var objs []*model.SpecObject
	for _, item := range array {
		spec := item.(map[string]interface{})
		meta := map[string]string{"projectId": project, "dbAlias": params["dbAlias"], "col": spec["col"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "col")
		delete(spec, "id")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/database/{dbAlias}/collections/{col}/modify-schema", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
