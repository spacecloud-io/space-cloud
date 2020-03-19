package remoteservices

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// GetRemoteServices gets remote services
func GetRemoteServices(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/remote-service/service", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := utils.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["service"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = params["service"]
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

		meta := map[string]string{"project": project, "id": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")
		delete(spec, "project")
		delete(spec, "version")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/remote-service/service/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}
