package eventing

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/transport"
)

// GetEventingTrigger gets eventing trigger
func GetEventingTrigger(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/triggers", project)

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
		delete(spec, "id")
		delete(spec, "name")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/eventing/triggers/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

// GetEventingConfig gets eventing config
func GetEventingConfig(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/config", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, params, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	for _, value := range payload.Result {
		meta := map[string]string{"project": project, "id": commandName}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/eventing/config/{id}", commandName, meta, value)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}

	return objs, nil
}

// GetEventingSchema gets eventing schema
func GetEventingSchema(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/schema", project)

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
		delete(spec, "id")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/eventing/schema/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}

// GetEventingSecurityRule gets eventing security rules
func GetEventingSecurityRule(project, commandName string, params map[string]string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/rules", project)

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
		delete(spec, "id")

		// Generating the object
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/eventing/rules/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
