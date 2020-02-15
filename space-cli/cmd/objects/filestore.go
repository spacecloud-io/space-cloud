package objects

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

func GetFileStoreConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/config", project)

	// Get the spec from the server
	result := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project}
	s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/file-storage/config", c.Command.Name, meta, result)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func GetFileStoreRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/rules", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["ruleName"] = c.Args()[0]
	}
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["rule"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = c.Args()[0]
		array = []interface{}{obj}
	}
	if value, p := result["rules"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}

	for _, item := range array {
		spec := item.(map[string]interface{})
		meta := map[string]string{"projectId": project, "id": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "name")
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/file-storage/rules/{id}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}
