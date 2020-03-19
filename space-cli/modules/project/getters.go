package project

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// GetProjectConfig gets global config
func GetProjectConfig(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {

	url := fmt.Sprintf("/v1/config/projects/%s", project)
	// Get the spec from the server
	result := make([]interface{}, 0)
	if err := utils.Get(http.MethodGet, url, params, &result); err != nil {
		return nil, err
	}

	// Generating the object
	specObjArr := make([]*model.SpecObject, len(result))
	for index, value := range result {
		projectObj := value.(map[string]interface{})
		meta := map[string]string{"project": projectObj["id"].(string)}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}", commandName, meta, projectObj)
		if err != nil {
			return nil, err
		}
		specObjArr[index] = s
	}
	return specObjArr, nil
}
