package letsencrypt

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//GetLetsEncryptDomain gets encrypt domain
func GetLetsEncryptDomain(project, commandName string, params map[string]string) (*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt", project)
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, &result); err != nil {
		return nil, err
	}

	// Printing the object on the screen
	meta := map[string]string{"project": project}
	s, err := utils.CreateSpecObject("/v1/config/projects/{project}/letsencrypt", commandName, meta, result["letsEncrypt"])
	if err != nil {
		return nil, err
	}

	return s, nil
}
