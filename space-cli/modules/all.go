package modules

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/modules/auth"
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/eventing"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/project"
	remoteservices "github.com/spaceuptech/space-cli/modules/remote-services"
	"github.com/spaceuptech/space-cli/modules/routes"
)

func GetAllProjects(c *cli.Context) error {
	projectName := c.GlobalString("project")

	obj, err := project.GetProjectConfig(projectName, "project", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("1", "project", []*model.SpecObject{obj}); err != nil {
		return err
	}

	objs, err := database.GetDbConfig(projectName, "db-config", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("2", "db-config", objs); err != nil {
		return err
	}

	objs, err = database.GetDbRule(projectName, "db-rule", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("3", "db-rule", objs); err != nil {
		return err
	}

	objs, err = database.GetDbSchema(projectName, "db-schema", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("4", "db-schema", objs); err != nil {
		return err
	}

	obj, err = filestore.GetFileStoreConfig(projectName, "filestore-config", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("5", "filestore-config", []*model.SpecObject{obj}); err != nil {
		return err
	}

	objs, err = filestore.GetFileStoreRule(projectName, "filestore-rule", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("6", "filestore-rule", objs); err != nil {
		return err
	}

	obj, err = eventing.GetEventingConfig(projectName, "eventing-config", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("7", "eventing-config", []*model.SpecObject{obj}); err != nil {
		return err
	}

	objs, err = eventing.GetEventingTrigger(projectName, "eventing-triggers", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("8", "eventing-triggers", objs); err != nil {
		return err
	}

	objs, err = eventing.GetEventingSecurityRule(projectName, "eventing-rule", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("9", "eventing-rule", objs); err != nil {
		return err
	}

	objs, err = eventing.GetEventingSchema(projectName, "eventing-schema", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("10", "eventing-schema", objs); err != nil {
		return err
	}

	objs, err = remoteservices.GetRemoteServices(projectName, "remote-services", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("12", "remote-services", objs); err != nil {
		return err
	}

	objs, err = routes.GetIngressRoutes(projectName, "ingress-routes", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("13", "-ingress-routes", objs); err != nil {
		return err
	}

	objs, err = auth.GetAuthProviders(projectName, "auth-providers", map[string]string{})
	if err != nil {
		return err
	}
	if err := createConfigFile("14", "auth-providers", objs); err != nil {
		return err
	}

	return nil
}

func createConfigFile(pos, commandName string, objs []*model.SpecObject) error {
	message := ""
	for _, val := range objs {
		data, err := yaml.Marshal(val)
		if err != nil {
			return err
		}
		message = message + string(data) + "---" + "\n"
	}
	fileName := fmt.Sprintf("%s-%s.yaml", pos, commandName)
	if err := ioutil.WriteFile(fileName, []byte(message), 0755); err != nil {
		return err
	}
	return nil
}
