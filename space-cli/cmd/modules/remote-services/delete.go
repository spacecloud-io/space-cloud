package remoteservices

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteRemoteService(project, prefix string) error {

	objs, err := GetRemoteServices(project, "remote-service", map[string]string{})
	if err != nil {
		return err
	}

	serviceIDs := []string{}
	for _, spec := range objs {
		serviceIDs = append(serviceIDs, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, serviceIDs)
	if err != nil {
		return err
	}

	// Delete the remote service from the server
	url := fmt.Sprintf("/v1/config/projects/%s/remote-service/service/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
