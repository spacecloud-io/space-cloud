package objects

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

func GetAuthProviders(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/user-management", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["provider"] = c.Args()[0]
	}
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["provider"]; p {
		obj := value.(map[string]interface{})
		obj["provider"] = c.Args()[0]
		array = []interface{}{obj}
	}
	if value, p := result["providers"]; p {
		obj := value.(map[string]interface{})
		for provider, value := range obj {
			o := value.(map[string]interface{})
			o["provider"] = provider
			array = append(array, o)
		}
	}

	for _, item := range array {
		spec := item.(map[string]interface{})
		meta := map[string]string{"projectId": project, "provider": spec["provider"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "provider")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/user-management/{provider}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}
