package remoteservices

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

func getRemoteServices(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/services", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["service"]; p {
		array = []interface{}{value}
	}
	if value, p := result["services"]; p {
		array = value.([]interface{})
	}
	var services []*model.SpecObject
	for _, item := range array {
		spec := item.(map[string]interface{})

		meta := map[string]string{"projectId": project, "id": spec["id"].(string), "version": spec["version"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")
		delete(spec, "project")
		delete(spec, "version")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/services/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}
