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

func deleteService(project string, args map[string]string) error {
	objs, err := GetServices(project, "service", map[string]string{})
	if err != nil {
		return err
	}

	serviceID := ""
	version := ""
	isExactMatch := false
	serviceIDs := []string{}
	for _, spec := range objs {
		// allow only those services that match version prefix provided by the user
		if args["version"] != "" && !strings.HasPrefix(spec.Meta["version"], args["version"]) {
			continue
		}
		serviceIDs = append(serviceIDs, fmt.Sprintf("%s::%s", spec.Meta["serviceId"], spec.Meta["version"]))
		if strings.EqualFold(spec.Meta["serviceId"], args["serviceId"]) && strings.EqualFold(spec.Meta["version"], args["version"]) {
			serviceID = spec.Meta["serviceId"]
			version = spec.Meta["version"]
			isExactMatch = true
			break
		}
	}

	if !isExactMatch {
		resourceID, err := filter.DeleteOptions(args["serviceId"], serviceIDs)
		if err != nil {
			return err
		}

		if resourceID == "" {
			return nil
		}

		resourceIDs := strings.Split(resourceID, "::")
		serviceID = resourceIDs[0]
		version = resourceIDs[1]
	}

	// Remove the deployed service from the space cloud
	url := fmt.Sprintf("/v1/runner/%s/services/%s/%s", project, serviceID, version)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteServiceRole(project string, args map[string]string) error {

	objs, err := GetServicesRole(project, "service-role", map[string]string{})
	if err != nil {
		return err
	}
	serviceID := ""
	roleID := ""
	isExactMatch := false
	serviceIDs := []string{}
	for _, spec := range objs {
		// allow only those services that match version prefix provided by the user
		if args["roleId"] != "" && !strings.HasPrefix(spec.Meta["roleId"], args["roleId"]) {
			continue
		}
		serviceIDs = append(serviceIDs, fmt.Sprintf("%s::%s", spec.Meta["serviceId"], spec.Meta["roleId"]))
		if strings.EqualFold(spec.Meta["serviceId"], args["serviceId"]) && strings.EqualFold(spec.Meta["roleId"], args["roleId"]) {
			serviceID = spec.Meta["serviceId"]
			roleID = spec.Meta["roleId"]
			isExactMatch = true
			break
		}
	}

	if !isExactMatch {
		resourceID, err := filter.DeleteOptions(args["serviceId"], serviceIDs)
		if err != nil {
			return err
		}

		if resourceID == "" {
			return nil
		}

		resourceIDs := strings.Split(resourceID, "::")
		serviceID = resourceIDs[0]
		roleID = resourceIDs[1]
	}

	// Remove the deployed service-role	 from the space cloud
	url := fmt.Sprintf("/v1/runner/%s/service-roles/%s/%s", project, serviceID, roleID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
