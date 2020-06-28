package filestore

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

// GetFileStoreConfig gets filestore config
func GetFileStoreConfig(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/config", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		meta := map[string]string{"project": project, "id": commandName}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/file-storage/config/{id}", commandName, meta, item)
		if err != nil {
			return nil, err
		}

		objs = append(objs, s)
	}

	return objs, nil
}

// GetFileStoreRule gets filestore rule
func GetFileStoreRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/rules", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}
	var objs []*model.SpecObject
	for _, item := range payload.Result {
		spec := item.(map[string]interface{})
		meta := map[string]string{"project": project, "id": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "name")
		delete(spec, "id")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/file-storage/rules/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
