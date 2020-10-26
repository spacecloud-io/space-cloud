package auth

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteAuthProvider(project, prefix string) error {

	objs, err := GetAuthProviders(project, "auth-provider", map[string]string{"id": "*"})
	if err != nil {
		return err
	}

	doesProviderExist := false
	providers := []string{}
	for _, spec := range objs {
		providers = append(providers, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, providers, doesProviderExist)
	if err != nil {
		return err
	}

	// Delete the provider from the server
	url := fmt.Sprintf("/v1/config/projects/%s/user-management/provider/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"id": resourceID}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
