package eventing

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
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
