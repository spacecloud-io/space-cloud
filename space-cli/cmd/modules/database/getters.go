package database

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

// GetDbRule gets database rule
func GetDbRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/rules", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		obj := item.(map[string]interface{})
		for key, value := range obj {
			str := strings.Split(key, "-")
			if str[1] == "event_logs" || str[1] == "invocation_logs" {
				continue
			}
			meta := map[string]string{"project": project, "col": str[1], "dbAlias": str[0]}

			delete(obj, "schema")

			// Generating the object
			s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules", commandName, meta, value)
			if err != nil {
				return nil, err
			}
			objs = append(objs, s)
		}
	}
	return objs, nil
}

// GetDbConfig gets database config
func GetDbConfig(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/config", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		spec := item.(map[string]interface{})
		for key, value := range spec {
			configID := fmt.Sprintf("%s-config", key)
			meta := map[string]string{"project": project, "dbAlias": key, "id": configID}

			// Delete the unwanted keys from spec
			delete(spec, "id")

			// Generating the object
			s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{dbAlias}/config/{id}", commandName, meta, value)
			if err != nil {
				return nil, err
			}
			objs = append(objs, s)
		}
	}

	return objs, nil
}

// GetDbSchema gets database schema
func GetDbSchema(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/schema/mutate", project)

	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, item := range payload.Result {
		obj := item.(map[string]interface{})
		for key, value := range obj {
			str := strings.Split(key, "-")
			if str[1] == "event_logs" || str[1] == "invocation_logs" {
				continue
			}
			meta := map[string]string{"project": project, "col": str[1], "dbAlias": str[0]}

			delete(obj, "isRealtimeEnabled")
			delete(obj, "rules")

			// Generating the object
			s, err := utils.CreateSpecObject("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate", commandName, meta, value)
			if err != nil {
				return nil, err
			}
			objs = append(objs, s)
		}
	}
	return objs, nil
}
