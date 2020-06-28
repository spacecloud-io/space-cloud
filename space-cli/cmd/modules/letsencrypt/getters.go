package letsencrypt

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

// GetLetsEncryptDomain gets encrypt domain
func GetLetsEncryptDomain(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt/config", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		meta := map[string]string{"project": project, "id": "letsencrypt"}
		delete(item.(map[string]interface{}), "id")
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/letsencrypt/config/{id}", commandName, meta, item)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}

	return objs, nil
}
