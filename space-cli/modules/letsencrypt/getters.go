package letsencrypt

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// GetLetsEncryptDomain gets encrypt domain
func GetLetsEncryptDomain(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt/config", project)
	// Get the spec from the server
	result := make([]interface{}, 0)
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, &result); err != nil {
		return nil, err
	}

	// Printing the object on the screen
	objs := []*model.SpecObject{}
	for _, value := range result {
		meta := map[string]string{"project": project, "id": commandName}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/letsencrypt/config/{id}", commandName, meta, value)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}

	return objs, nil
}
