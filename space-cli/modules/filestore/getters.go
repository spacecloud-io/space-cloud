package filestore

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

func getFileStoreConfig(project, commandName string, params map[string]string) (*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/config", project)
	// Get the spec from the server
	result := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return nil, err
	}

	// Generating the object
	meta := map[string]string{"projectId": project}
	s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/file-storage/config", commandName, meta, result)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func GetFileStoreRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/rules", project)
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["rule"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = params["ruleName"]
		array = []interface{}{obj}
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
		meta := map[string]string{"projectId": project, "id": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "name")
		delete(spec, "id")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/file-storage/rules/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
