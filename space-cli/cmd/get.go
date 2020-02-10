package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

func GetGlobalConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/config", project)

	// Get the spec from the server
	result := new(interface{})
	if err := get(http.MethodGet, url, map[string]string{}, result); err != nil {
		return err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project}
	s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/config", c.Command.Name, meta, result)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

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
	if err := get(http.MethodGet, url, params, &result); err != nil {
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
	if err := get(http.MethodGet, url, params, &result); err != nil {
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
	if err := get(http.MethodGet, url, params, &result); err != nil {
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
	if err := get(http.MethodGet, url, map[string]string{}, vPtr); err != nil {
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
	if err := get(http.MethodGet, url, params, &result); err != nil {
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
	if err := get(http.MethodGet, url, params, &result); err != nil {
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

func GetFileStoreConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/file-storage/config", project)

	// Get the spec from the server
	result := new(interface{})
	if err := get(http.MethodGet, url, map[string]string{}, result); err != nil {
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
	if err := get(http.MethodGet, url, params, &result); err != nil {
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

func GetDbRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/rules", project)

	params := map[string]string{}
	switch len(c.Args()) {
	case 1:
		params["dbType"] = c.Args()[0]
	case 2:
		params["dbType"] = c.Args()[0]
		params["col"] = c.Args()[1]
	}

	// Get the spec from the server
	result := make(map[string]interface{})
	if err := get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["rule"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
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
		meta := map[string]string{"projectId": project, "id": spec["id"].(string), "dbType": c.Args()[0]}

		// Delete the unwanted keys from spec
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("v1/config/projects/{projectId}/database/{dbType}/collections/{id}/rules", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}

func GetDbConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/database/config", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["dbType"] = c.Args()[0]
		//params["col"] = c.Args()[1]
	}

	// Get the spec from the server
	result := new(interface{})
	if err := get(http.MethodGet, url, params, result); err != nil {
		return err
	}

	// Printing the object on the screen
	meta := map[string]string{"projectId": project, "dbType": c.Args()[0]}
	s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/database/{dbType}/config", c.Command.Name, meta, result)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func GetDbSchema(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/database/collections/modify-schema", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["dbType"] = c.Args()[0]
	}
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := get(http.MethodGet, url, params, &result); err != nil {
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
		meta := map[string]string{"projectId": project, "dbType": c.Args()[0]}

		// Delete the unwanted keys from spec
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{projectId}/file-storage/rules/{dbType}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}

func GetLetsEncryptDomain(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/letsencrypt", project)

	// Get the spec from the server
	result := new(interface{})
	if err := get(http.MethodGet, url, map[string]string{}, result); err != nil {
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

func GetRoutes(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	url := fmt.Sprintf("/v1/config/projects/%s/routing/route", project)

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["routesId"] = c.Args()[0]
	}
	// Get the spec from the server
	result := make(map[string]interface{})
	if err := get(http.MethodGet, url, params, &result); err != nil {
		return err
	}

	var array []interface{}
	if value, p := result["route"]; p {
		obj := value.(map[string]interface{})
		obj["id"] = c.Args()[0]
		array = []interface{}{obj}
	}
	if value, p := result["routes"]; p {
		obj := value.(map[string]interface{})
		for rule, value := range obj {
			o := value.(map[string]interface{})
			o["id"] = rule
			array = append(array, o)
		}
	}

	for _, item := range array {
		spec := item.(map[string]interface{})
		meta := map[string]string{"projectId": project, "routeId": spec["id"].(string)}

		// Delete the unwanted keys from spec
		delete(spec, "id")

		// Printing the object on the screen
		s, err := utils.GetYamlObject("/v1/config/projects/{project}/routing/{routeId}", c.Command.Name, meta, spec)
		if err != nil {
			return err
		}
		fmt.Print(s)
		fmt.Println("---")
	}
	return nil
}

func get(method, url string, params map[string]string, vPtr interface{}) error {
	account, err := getSelectedAccount()
	if err != nil {
		return err
	}
	login, err := login(account)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%s%s", account.ServerUrl, url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	if login.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", login.Token))
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("data:", string(data))
	if err := json.Unmarshal(data, vPtr); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("received invalid status code (%d)", resp.StatusCode)
	}

	return nil
}
