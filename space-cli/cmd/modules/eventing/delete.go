package eventing

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteEventingConfig(project string) error {

	// Delete the filestore config from the server
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/config/%s", project, "eventing-config")

	if err := transport.Client.MakeHTTPRequest(http.MethodPost, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteEventingTriggers(project string, prefix string) error {

	objs, err := GetEventingTrigger(project, "eventing-trigger", map[string]string{})
	if err != nil {
		return err
	}

	triggerIDs := []string{}
	for _, spec := range objs {
		triggerIDs = append(triggerIDs, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, triggerIDs)
	if err != nil {
		return err
	}

	// Delete the eventing trigger from the server
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/triggers/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteEventingSchemas(project string, prefix string) error {

	objs, err := GetEventingSchema(project, "eventing-schema", map[string]string{})
	if err != nil {
		return err
	}

	schemaIDs := []string{}
	for _, spec := range objs {
		schemaIDs = append(schemaIDs, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, schemaIDs)
	if err != nil {
		return err
	}

	// Delete the eventing schema from the server
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/schema/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteEventingRules(project string, prefix string) error {

	objs, err := GetEventingSecurityRule(project, "eventing-rule", map[string]string{})
	if err != nil {
		return err
	}

	ruleIDs := []string{}
	for _, spec := range objs {
		ruleIDs = append(ruleIDs, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, ruleIDs)
	if err != nil {
		return err
	}

	// Delete the eventing schema from the server
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/rules/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
