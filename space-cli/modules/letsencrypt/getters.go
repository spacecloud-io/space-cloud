package letsencrypt

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// GetLetsEncryptDomain gets encrypt domain
func GetLetsEncryptDomain(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt/config", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := utils.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		meta := map[string]string{"project": project, "id": "letsencrypt-config"}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/letsencrypt/config/{id}", commandName, meta, item)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}

	return objs, nil
}
