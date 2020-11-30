package addons

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

const helmDatabaseNamespace = "db"

func addDatabase(chartReleaseName, dbType, setValuesFlag, valuesYamlFile, chartLocation string) error {
	valuesFileObj, err := utils.ExtractValuesObj(setValuesFlag, valuesYamlFile)
	if err != nil {
		return err
	}

	// The regex stratifies kubernetes resource name specification
	var validID = regexp.MustCompile(`[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`)
	if !validID.MatchString(chartReleaseName) {
		return fmt.Errorf(`invalid name for database: (%s): a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character (e.g. 'example.com', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*'`, chartReleaseName)
	}

	valuesFileObj["name"] = chartReleaseName

	downloadURL := ""
	switch dbType {
	case "postgres":
		downloadURL = model.HelmPostgresChartDownloadURL
	case "mysql":
		downloadURL = model.HelmMysqlChartDownloadURL
	case "sqlserver":
		downloadURL = model.HelmSQLServerCloudChartDownloadURL
	case "mongo":
		downloadURL = model.HelmMongoChartDownloadURL
	default:
		return fmt.Errorf("unkown database (%s) provided as argument", chartReleaseName)
	}

	_, err = utils.HelmInstall(chartReleaseName, chartLocation, downloadURL, helmDatabaseNamespace, valuesFileObj)

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}

	kubeClint, err := actionConfig.KubernetesClientSet()
	w, err := kubeClint.CoreV1().Pods(helmDatabaseNamespace).Watch(context.Background(), metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", chartReleaseName)})
	if err != nil {
		return nil
	}

	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	defer w.Stop()

	utils.LogInfo("Waiting for database to start...")
	for {
		select {
		case p := <-w.ResultChan():
			pod, ok := p.Object.(*v1.Pod)
			if !ok {
				continue
			}
			if pod.Status.Phase == v1.PodRunning {
				utils.LogInfo(fmt.Sprintf("Database is up & running, use this domain name (%s.%s.svc.cluster.local) in mission control for adding database to space cloud", chartReleaseName, helmDatabaseNamespace))
				return nil
			} else {
				utils.LogInfo(fmt.Sprintf("Database pod status (%s)", pod.Status.Phase))
			}
		case <-ticker.C:
			utils.LogInfo("Database has been provisioned on kubernetes cluster")
			utils.LogInfo("But it is taking to much time to start, check if you have enough resource on the cluster or")
			return nil
		}
	}
}

func removeDatabase(dbType string) error {
	if err := utils.HelmUninstall(dbType); err != nil {
		return err
	}
	utils.LogInfo(fmt.Sprintf("Removed database (%s) from kubernetes", dbType))
	return nil
}
