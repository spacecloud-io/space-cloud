package objects

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
)

func GetLetsEncryptDomain(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt", project)

	// Get the spec from the server
	result := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project}
	s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/letsencrypt", c.Command.Name, meta, result)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}
