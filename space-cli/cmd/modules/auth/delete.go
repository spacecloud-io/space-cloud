package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
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

	filteredProviders := []string{}
	for _, provider := range providers {
		if prefix != "" && strings.HasPrefix(strings.ToLower(provider), strings.ToLower(prefix)) {
			filteredProviders = append(filteredProviders, provider)
			doesProviderExist = true
		}
	}

	if doesProviderExist {
		if len(filteredProviders) == 1 {
			prefix = filteredProviders[0]
		} else {
			if err := input.Survey.AskOne(&survey.Select{Message: "Choose the account ID: ", Options: filteredProviders, Default: filteredProviders[0]}, &prefix); err != nil {
				return err
			}
		}
	} else {
		if len(providers) == 1 {
			prefix = providers[0]
		} else {
			if prefix != "" {
				utils.LogInfo("Warning! No provider found for prefix provided, showing all")
			}
			if err := input.Survey.AskOne(&survey.Select{Message: "Choose the account ID: ", Options: providers, Default: providers[0]}, &prefix); err != nil {
				return err
			}
		}
	}

	// Delete the provider from the server
	url := fmt.Sprintf("/v1/config/projects/%s/user-management/provider/%s", project, prefix)

	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{"id": prefix}, payload); err != nil {
		return err
	}

	return nil
}
