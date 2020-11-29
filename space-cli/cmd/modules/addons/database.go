package addons

import (
	"fmt"
	"regexp"

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
	return err
}

func removeDatabase(dbType string) error {
	if err := utils.HelmUninstall(dbType); err != nil {
		return err
	}
	utils.LogInfo(fmt.Sprintf("Removed database (%s) from kubernetes", dbType))
	return nil
}
