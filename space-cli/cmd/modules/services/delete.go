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
