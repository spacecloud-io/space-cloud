package encrypt

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//GetLetsEncrypt gets encrypt domain
func GetLetsEncryptDomain(project, commandName string, params map[string]string) (*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt", project)
	// Get the spec from the server
	result := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return nil, err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project}
	s, err := utils.CreateSpecObject("/v1/config/projects/{projectId}/letsencrypt", commandName, meta, result)
	if err != nil {
		return nil, err
	}

	return s, nil
}
