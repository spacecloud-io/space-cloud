package filter

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
)

// DeleteOptions filters the resource id based on prefix to delete
func DeleteOptions(prefix string, resources []string) (string, error) {
	doesResourceExist := false
	filteredResources := []string{}
	for _, resource := range resources {
		if prefix != "" && strings.HasPrefix(strings.ToLower(resource), strings.ToLower(prefix)) {
			filteredResources = append(filteredResources, resource)
			doesResourceExist = true
		}
	}

	if doesResourceExist {
		if err := input.Survey.AskOne(&survey.Select{Message: "Choose the resource ID: ", Options: filteredResources, Default: filteredResources[0]}, &prefix); err != nil {
			return "", err
		}
	} else {
		if prefix != "" {
			return "", utils.LogError(fmt.Sprintf("Warning! No resource found for prefix-(%s)", prefix), nil)
		}
		if err := input.Survey.AskOne(&survey.Select{Message: "Choose the resource ID: ", Options: resources, Default: resources[0]}, &prefix); err != nil {
			return "", err
		}

	}

	return prefix, nil
}
