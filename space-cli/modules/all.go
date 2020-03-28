package modules

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/modules/auth"
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/eventing"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/ingress"
	"github.com/spaceuptech/space-cli/modules/letsencrypt"
	"github.com/spaceuptech/space-cli/modules/project"
	remoteservices "github.com/spaceuptech/space-cli/modules/remote-services"
	"github.com/spaceuptech/space-cli/modules/services"
)

func getSubCommands() []*cobra.Command {

	var getProjects = &cobra.Command{
		Use:   "all",
		Short: "Gets entire project config",
		RunE:  getAllProjects,
	}

	return []*cobra.Command{getProjects}
}

func getAllProjects(cmd *cobra.Command, args []string) error {
	projectName := viper.GetString("project")

	if len(args) == 0 {
		return fmt.Errorf("Directory not specified as an arguement to store config files")
	}
	dir := args[0]
	// create directory if directory doesn't exists
	_ = os.MkdirAll(dir, os.ModePerm)

	_ = os.Chdir(dir)

	obj, _ := project.GetProjectConfig(projectName, "project", map[string]string{})

	_ = createConfigFile("1", "project", obj)

	objs, _ := database.GetDbConfig(projectName, "db-config", map[string]string{})

	_ = createConfigFile("2", "db-config", objs)

	objs, _ = database.GetDbRule(projectName, "db-rules", map[string]string{})

	_ = createConfigFile("3", "db-rules", objs)

	objs, _ = database.GetDbSchema(projectName, "db-schema", map[string]string{})

	_ = createConfigFile("4", "db-schema", objs)

	obj, _ = filestore.GetFileStoreConfig(projectName, "filestore-config", map[string]string{})

	_ = createConfigFile("5", "filestore-config", obj)

	objs, _ = filestore.GetFileStoreRule(projectName, "filestore-rule", map[string]string{})

	_ = createConfigFile("6", "filestore-rule", objs)

	obj, _ = eventing.GetEventingConfig(projectName, "eventing-config", map[string]string{})

	_ = createConfigFile("7", "eventing-config", obj)

	objs, _ = eventing.GetEventingTrigger(projectName, "eventing-triggers", map[string]string{})

	_ = createConfigFile("8", "eventing-triggers", objs)

	objs, _ = eventing.GetEventingSecurityRule(projectName, "eventing-rule", map[string]string{})

	_ = createConfigFile("9", "eventing-rule", objs)

	objs, _ = eventing.GetEventingSchema(projectName, "eventing-schema", map[string]string{})

	_ = createConfigFile("10", "eventing-schema", objs)

	objs, _ = remoteservices.GetRemoteServices(projectName, "remote-services", map[string]string{})

	_ = createConfigFile("11", "remote-services", objs)

	objs, _ = services.GetServices(projectName, "services", map[string]string{})

	_ = createConfigFile("12", "services", objs)

	objs, _ = services.GetServicesRoutes(projectName, "services-routes", map[string]string{})

	_ = createConfigFile("13", "services-routes", objs)

	// objs, _ = services.GetServicesSecrets(projectName, "services-secrets", map[string]string{})
	// if _ != nil {
	// 	return _
	// }
	// _ = createConfigFile("14", "services-secrets", objs); _ != nil {
	// 	return _
	// }

	objs, _ = ingress.GetIngressRoutes(projectName, "ingress-routes", map[string]string{})

	_ = createConfigFile("15", "-ingress-routes", objs)

	objs, _ = auth.GetAuthProviders(projectName, "auth-providers", map[string]string{})

	_ = createConfigFile("16", "auth-providers", objs)

	obj, _ = letsencrypt.GetLetsEncryptDomain(projectName, "letsencrypt", map[string]string{})

	_ = createConfigFile("17", "letsencrypt", obj)

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
