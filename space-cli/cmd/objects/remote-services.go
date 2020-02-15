package objects

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

func GetRemoteServices(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/services", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["service"] = c.Args()[0]
	}

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["service"]; p {
		array = []interface{}{value}
	}
	if value, p := result["services"]; p {
		array = value.([]interface{})
	}

	for _, item := range array {
		spec := item.(map[string]interface{})

		meta := map[string]string{"projectId": project, "id": spec["id"].(string), "version": spec["version"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")
		delete(spec, "project")
		delete(spec, "version")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/services/{id}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Println(s)
		fmt.Println("---")
	}

	return nil
}
