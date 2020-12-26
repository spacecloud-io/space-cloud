package services

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteSecret(project, prefix string) error {

	objs, err := GetServicesSecrets(project, "secret", map[string]string{})
	if err != nil {
		return err
	}

	secretIDs := []string{}
	for _, spec := range objs {
		secretIDs = append(secretIDs, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, secretIDs)
	if err != nil {
		return err
	}

	// Delete the remote service from the server
	url := fmt.Sprintf("/v1/runner/%s/secrets/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteService(project string, prefix map[string]string) error {
	objs, err := GetServices(project, "service", map[string]string{})
	if err != nil {
		return err
	}

	serviceIDs := []string{}
	for _, spec := range objs {
		serviceIDs = append(serviceIDs, spec.Meta["serviceId"])
	}

	resourceID, err := filter.DeleteOptions(prefix["serviceId"], serviceIDs)
	if err != nil {
		return err
	}

	versions := []string{}
	for _, spec := range objs {
		versions = append(versions, spec.Meta["version"])
	}

	version, err := filter.DeleteOptions(prefix["version"], versions)
	if err != nil {
		return err
	}

	// Delete the remote service from the server
	url := fmt.Sprintf("/v1/runner/%s/secrets/%s/%s", project, resourceID, version)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

// func deleteServiceRoute(project, prefix string) error {

// }

func deleteServiceRole(project string, prefix map[string]string) error {

	objs, err := GetServicesRole(project, "service-role", map[string]string{})
	if err != nil {
		return err
	}

	serviceIDs := []string{}
	for _, spec := range objs {
		serviceIDs = append(serviceIDs, spec.Meta["serviceId"])
	}

	resourceID, err := filter.DeleteOptions(prefix["serviceID"], serviceIDs)
	if err != nil {
		return err
	}

	roleIDs := []string{}
	for _, spec := range objs {
		roleIDs = append(roleIDs, spec.Meta["roleId"])
	}

	roleID, err := filter.DeleteOptions(prefix["roleID"], roleIDs)
	if err != nil {
		return err
	}

	// Delete the remote service from the server
	url := fmt.Sprintf("/v1/runner/%s/secrets/%s/%s", project, resourceID, roleID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
