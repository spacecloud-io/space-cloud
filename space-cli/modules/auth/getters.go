package auth

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//GetAuthProviders gets auth providers
func GetAuthProviders(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/user-management/provider", project)

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := utils.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	var array []interface{}
	if value, p := result["provider"]; p {
		obj := value.(map[string]interface{})
		obj["provider"] = params["provider"]
		array = []interface{}{obj}
	}
	if value, p := result["providers"]; p {
		obj := value.(map[string]interface{})
		for provider, value := range obj {
			o := value.(map[string]interface{})
			o["provider"] = provider
			array = append(array, o)
		}
	}
	var objs []*model.SpecObject
	for _, item := range array {
		spec := item.(map[string]interface{})
		meta := map[string]string{"project": project, "provider": spec["provider"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "provider")

		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/user-management/provider/{provider}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
