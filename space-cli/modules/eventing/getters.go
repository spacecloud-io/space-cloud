package eventing

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

func getEventingTrigger(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/triggers", project)

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
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/eventing/triggers/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

func getEventingConfig(project, commandName string, params map[string]string) (*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/config", project)
	// Get the spec from the server
	vPtr := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, vPtr); err != nil {
		return nil, err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project}
	s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/eventing/config", commandName, meta, vPtr)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func getEventingSchema(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/schema", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["schema"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = params["type"]
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
		meta := map[string]string{"projectId": project, "id": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/eventing/schema/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

func getEventingSecurityRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/rules", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["securityRule"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = params["type"]
		array = []interface{}{obj}
	}
	if value, p := result["securityRules"]; p {
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
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/eventing/rules/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
