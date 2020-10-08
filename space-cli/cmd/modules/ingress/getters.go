package ingress

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

// GetIngressRoutes gets ingress routes
func GetIngressRoutes(project, commandName string, params map[string]string, filters []string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/routing/ingress", project)
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
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/routing/ingress/{id}", commandName, meta, spec)
		if err != nil {
			return nil, err
		}

		if len(filters) > 0 {
			for _, filter := range filters {
				arr := strings.Split(filter, "=")
				if len(arr) != 2 {
					continue
				}
				filterKey := arr[0]
				filterValue := arr[1]
				switch filterKey {
				case "url":
					value, ok := spec["source"].(map[string]interface{})
					if !ok {
						continue
					}
					if filterValue == value["url"].(string) {
						objs = append(objs, s)
					}
				case "service":
					targets, ok := spec["targets"].([]interface{})
					if !ok {
						continue
					}
					hostName := fmt.Sprintf("%s.%s.svc.cluster.local", filterValue, project)
					for _, target := range targets {
						targetObj, ok := target.(map[string]interface{})
						if !ok {
							continue
						}
						if hostName == targetObj["host"] {
							objs = append(objs, s)
						}
					}
				case "target-host":
					targets, ok := spec["targets"].([]interface{})
					if !ok {
						continue
					}
					for _, target := range targets {
						targetObj, ok := target.(map[string]interface{})
						if !ok {
							continue
						}
						if filterValue == targetObj["host"] {
							objs = append(objs, s)
						}
					}
				case "request-host":
					value, ok := spec["source"].(map[string]interface{})
					if !ok {
						continue
					}
					for _, requestHost := range value["hosts"].([]interface{}) {
						if filterValue == requestHost.(string) {
							objs = append(objs, s)
						}
					}

				}
			}
			continue
		}
		objs = append(objs, s)
	}
	return objs, nil
}

// GetIngressGlobal gets ingress global
func GetIngressGlobal(project, commandName string) ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/routing/ingress/global", project)
	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.Get(http.MethodGet, url, map[string]string{}, payload); err != nil {
		return nil, err
	}
	var objs []*model.SpecObject
	for _, item := range payload.Result {
		if item == nil {
			continue
		}
		spec := item.(map[string]interface{})
		meta := map[string]string{"project": project}
		s, err := utils.CreateSpecObject("/v1/config/projects/{project}/routing/ingress/global", commandName, meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}
	return objs, nil
}
