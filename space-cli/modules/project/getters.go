package project

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetProjectConfig gets global config
func GetProjectConfig(project, commandName string, params map[string]string) (*model.SpecObject, error) {

	url := fmt.Sprintf("/v1/config/projects/%s", project)
	// Get the spec from the server
	result := new(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return nil, err
	}

	// Generating the object
	meta := map[string]string{"project": project}
	s, err := utils.CreateSpecObject("/v1/config/projects/{project}", commandName, meta, result)
	if err != nil {
		return nil, err
	}

	return s, nil
}
