package letsencrypt

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteLetsencryptDomains(project string) error {
	// Delete the letsencrpyt domains from the server
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt/config/%s", project, "letsencrypt")

	if err := transport.Client.MakeHTTPRequest(http.MethodPost, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
