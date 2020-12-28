package services

import (
	"fmt"
	"net/http"
	"strings"

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

	serviceID := ""
	version := ""
	doesExists := false
	doesPartialExists := false
	serviceIDs := []string{}
	for _, spec := range objs {
		if prefix["version"] != "" && !strings.HasPrefix(strings.ToLower(spec.Meta["version"]), strings.ToLower(prefix["version"])) {
			continue
		}
		serviceIDs = append(serviceIDs, fmt.Sprintf("%s-%s", spec.Meta["serviceId"], spec.Meta["version"]))
		if strings.ToLower(spec.Meta["serviceId"]) == strings.ToLower(prefix["serviceId"]) && strings.ToLower(spec.Meta["version"]) == strings.ToLower(prefix["version"]) {
			serviceID = spec.Meta["serviceId"]
			version = spec.Meta["version"]
			doesExists = true
		}
		if strings.ToLower(spec.Meta["serviceId"]) == strings.ToLower(prefix["serviceId"]) {
			doesPartialExists = true
		}
	}

	if !doesExists {
		pre := prefix["serviceId"]
		if doesPartialExists {
			pre = fmt.Sprintf("%s-%s", prefix["serviceId"], prefix["version"])
		}

		resourceID, err := filter.DeleteOptions(pre, serviceIDs)
		if err != nil {
			return err
		}

		resourceIDs := strings.Split(resourceID, "-")
		serviceID = resourceIDs[0]
		version = resourceIDs[1]
	}

	// Delete the remote service from the server
	url := fmt.Sprintf("/v1/runner/%s/services/%s/%s", project, serviceID, version)

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
	serviceID := ""
	roleID := ""
	doesExists := false
	doesPartialExists := false
	serviceIDs := []string{}
	for _, spec := range objs {
		if prefix["roleId"] != "" && !strings.HasPrefix(strings.ToLower(spec.Meta["roleId"]), strings.ToLower(prefix["roleId"])) {
			continue
		}
		serviceIDs = append(serviceIDs, fmt.Sprintf("%s-%s", spec.Meta["serviceId"], spec.Meta["roleId"]))
		if strings.ToLower(spec.Meta["serviceId"]) == strings.ToLower(prefix["serviceId"]) && strings.ToLower(spec.Meta["roleId"]) == strings.ToLower(prefix["roleId"]) {
			serviceID = spec.Meta["serviceId"]
			roleID = spec.Meta["roleId"]
			doesExists = true
		}
		if strings.ToLower(spec.Meta["serviceId"]) == strings.ToLower(prefix["serviceId"]) {
			doesPartialExists = true
		}
	}

	if !doesExists {
		pre := prefix["serviceId"]
		if doesPartialExists {
			pre = fmt.Sprintf("%s-%s", prefix["serviceId"], prefix["roleId"])
		}

		resourceID, err := filter.DeleteOptions(pre, serviceIDs)
		if err != nil {
			return err
		}

		resourceIDs := strings.Split(resourceID, "-")
		serviceID = resourceIDs[0]
		roleID = resourceIDs[1]
	}

	// Delete the remote service from the server
	url := fmt.Sprintf("/v1/runner/%s/service-roles/%s/%s", project, serviceID, roleID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
