package project

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

// GetProjectConfig gets global config
func GetProjectConfig(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	if project == "" {
		project = "*" // for getting all projects
		value, ok := params["id"]
		if ok {
			project = value
		}
	}
	url := fmt.Sprintf("/v1/config/projects/%s", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		projectObj := item.(map[string]interface{})
		meta := map[string]string{"project": projectObj["id"].(string)}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}", commandName, meta, projectObj)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
