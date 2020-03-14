package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//GetServicesRoutes gets services routes
func GetServicesRoutes(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/runner/%s/service-routes", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	b, _ := json.MarshalIndent(result, "", " ")
	fmt.Println("result: ", b)

	var array []interface{}
	if value, p := result["service"]; p {
		array = []interface{}{value}
	}
	if value, p := result["services"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}
	var services []*model.SpecObject
	for _, item := range array {
		spec := item.(map[string]interface{})

		meta := map[string]string{"projectId": project, "id": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")
		delete(spec, "project")
		delete(spec, "version")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/runner/{project}/service-routes/{serviceId}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}

//GetServicesSecrets gets services secrets
func GetServicesSecrets(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/runner/%s/secrets", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["secret"]; p {
		array = []interface{}{value}
	}
	if value, p := result["secrets"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}
	var services []*model.SpecObject
	for _, item := range array {
		spec := item.(map[string]interface{})

		meta := map[string]string{"project": project, "name": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")
		delete(spec, "name")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("//v1/runner/{project}/secrets/{name}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}

//GetServices gets services
func GetServices(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/runner/%s/services", project)
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["service"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}
	if value, p := result["services"]; p {
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
		str := strings.Split(spec["id"].(string), "-")
		meta := map[string]string{"projectId": project, "version": str[1], "serviceId": str[0]}

		// Delete the unwanted keys from spec
		delete(spec, "id")
		delete(spec, "name")
		delete(spec, "version")
		delete(spec, "projectId")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/runner/{project}/services/{serviceId}/{version}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
