package operations

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Destroy cleans the environment which has been setup. It removes the containers, secrets & host file
func Destroy() error {

	charList, err := utils.HelmList(model.HelmSpaceCloudNamespace)
	if err != nil {
		return err
	}
	if len(charList) < 1 {
		utils.LogInfo("space cloud cluster not found, setup a new cluster using the setup command")
		return nil
	}

	isOk := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Space cloud cluster with id (%s) will be destoryed, Do you want to continue", charList[0].Name),
	}
	if err := survey.AskOne(prompt, &isOk); err != nil {
		return err
	}
	if !isOk {
		return nil
	}

	// Delete all projects
	objects, err := project.GetProjectConfig("*", "projects", nil)
	if err != nil {
		return err
	}
	for _, object := range objects {
		projectID, ok := object.Meta["project"]
		if !ok {
			continue
		}
		if err := deleteProject(projectID); err != nil {
			return err
		}
	}

	if err := utils.HelmUninstall(charList[0].Name); err != nil {
		return err
	}

	if err := utils.RemoveAccount(charList[0].Name); err != nil {
		return err
	}
	utils.LogInfo("Space cloud cluster has been destroyed successfully 😢")
	return nil
}

func deleteProject(projectID string) error {
	// delete project from kubernetes
	if err := project.DeleteProject(projectID); err != nil {
		return err
	}

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}

	kubeClint, err := actionConfig.KubernetesClientSet()
	if err != nil {
		return err
	}

	// wait for project to get deleted in kubernetes
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	maxCount := 18 // wait for 3 minutes, 10 seconds * 18 = 180 seconds
	counter := 0
	utils.LogInfo(fmt.Sprintf("Waiting for project (%s) to get deleted, this might take up to 3 minutes", projectID))

	for range ticker.C {
		namespacesList, err := kubeClint.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=space-cloud"})
		if err != nil {
			return err
		}
		doesExists := false
		var ns v1.Namespace
		for _, namespace := range namespacesList.Items {
			if namespace.Name == projectID {
				doesExists = true
				ns = namespace
				break
			}
		}
		if !doesExists {
			utils.LogInfo(fmt.Sprintf("Successfully deleted project (%s)", projectID))
			return nil
		}

		counter++
		if counter == maxCount {
			utils.LogInfo(fmt.Sprintf("Deleting project (%s) is taking to much time, skipping project (%s)", projectID, projectID))
			return nil
		}

		utils.LogInfo(fmt.Sprintf("Project deletion status (%s)", ns.Status.Phase))
	}
	return nil
}
