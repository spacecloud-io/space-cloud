package filestore

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteFileStoreConfig(project string) error {

	// Delete the filestore config from the server
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/config/%s", project, "filestore-config")

	if err := transport.Client.MakeHTTPRequest(http.MethodPost, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}

func deleteFileStoreRule(project, prefix string) error {

	objs, err := GetFileStoreRule(project, "filestore-rule", map[string]string{})
	if err != nil {
		return err
	}

	doesRuleExist := false
	ruleIDs := []string{}
	for _, spec := range objs {
		ruleIDs = append(ruleIDs, spec.Meta["id"])
	}

	resourceID, err := filter.DeleteOptions(prefix, ruleIDs, doesRuleExist)
	if err != nil {
		return err
	}

	// Delete the filestore rule from the server
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/rules/%s", project, resourceID)

	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, new(model.Response)); err != nil {
		return err
	}

	return nil
}
