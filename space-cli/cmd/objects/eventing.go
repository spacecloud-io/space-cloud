package objects

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

func GetEventingTrigger(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/event-triggers/rules", project)

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
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/event-triggers/rules/{id}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}

func GetEventingConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/event-triggers/config", project)

	// Get the spec from the server
	vPtr := new(interface{})
	if err := cmd.Get(http.MethodGet, url, map[string]string{}, vPtr); err != nil {
		return err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project}
	s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/event-triggers/config", c.Command.Name, meta, vPtr)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func GetEventingSchema(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/schema", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["type"] = c.Args()[0]
	}
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["schema"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = c.Args()[0]
		array = []interface{}{obj}
	}
	if value, p := result["schemas"]; p {
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
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/eventing/schema/{id}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}

func GetEventingSecurityRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/eventing/rules", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["type"] = c.Args()[0]
	}
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := cmd.Get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["securityRule"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = c.Args()[0]
		array = []interface{}{obj}
	}
	if value, p := result["securityRules"]; p {
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
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/eventing/rules/{id}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}
