package project

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteProject(project string) error {
	// Delete the project config from the server
	url := fmt.Sprintf("/v1/config/projects/%s", project)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
