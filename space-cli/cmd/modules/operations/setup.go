package operations

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Setup initializes development environment
func Setup(setValuesFlag, valuesYamlFile, chartLocation, version string) error {
	_ = utils.CreateDirIfNotExist(utils.GetSpaceCloudDirectory())

	valuesFileObj, err := utils.ExtractValuesObj(setValuesFlag, valuesYamlFile)
	if err != nil {
		return err
	}

	// override the clusterId
	newClusterID := ""
	clusterID, ok := valuesFileObj["clusterId"]
	if !ok {
		newClusterID = "space-cloud-" + uuid.New().String()
		valuesFileObj["clusterId"] = newClusterID
	} else {
		value, ok := clusterID.(string)
		if !ok {
			return fmt.Errorf("clusterId should be of type string got (%v)", reflect.TypeOf(clusterID))
		}
		newClusterID = value + "-" + uuid.New().String()
		valuesFileObj["clusterId"] = newClusterID
	}

	_, ok = valuesFileObj["version"]
	if !ok {
		// set the version
		valuesFileObj["version"] = model.Version
	}

	helmChart, err := utils.HelmInstall(newClusterID, chartLocation, utils.GetHelmChartDownloadURL(model.HelmSpaceCloudChartDownloadURL, version), model.HelmSpaceCloudNamespace, valuesFileObj)
	if err != nil {
		return err
	}

	selectedAccount := model.Account{
		ID:        newClusterID,
		UserName:  helmChart.Values["admin"].(map[string]interface{})["username"].(string),
		Key:       helmChart.Values["admin"].(map[string]interface{})["password"].(string),
		ServerURL: "http://localhost:4122",
	}

	fmt.Println()
	utils.LogInfo(fmt.Sprintf("Space Cloud (cluster id: \"%s\") has been successfully setup! 👍", selectedAccount.ID))
	utils.LogInfo(fmt.Sprintf("You can visit mission control at %s/mission-control 💻", selectedAccount.ServerURL))
	utils.LogInfo("Note: The url is only valid if you have done port forwarding using the commnad as per the docs at https://docs.spaceuptech.com/install/kubernetes/minikube/")
	utils.LogInfo("Command => kubectl port-forward -n istio-system deployments/istio-ingressgateway 4122:8080")
	utils.LogInfo("If you have done forwarding on other port, use the login command to change the url")
	utils.LogInfo(fmt.Sprintf("Your login credentials: [username: \"%s\"; key: \"%s\"] 🤫", selectedAccount.UserName, selectedAccount.Key))
	return nil
}
