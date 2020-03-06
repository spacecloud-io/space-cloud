package globalconfig

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

func getGlobalConfig(project, commandName string, params map[string]string) (*model.SpecObject, error) {

	url := fmt.Sprintf("/v1/config/projects/%s/config", project)
	// Get the spec from the server
	result := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return nil, err
	}

	// Generating the object
	meta := map[string]string{"projectId": project}
	s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/config", commandName, meta, result)
	if err != nil {
		return nil, err
	}

	return s, nil
}
