package operations

import (
	"fmt"
	"reflect"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

const spaceCloudChart = "space-cloud"

// Setup initializes development environment
func Setup(setValuesFlag, valuesYamlFile, chartLocation string) error {
	valuesFileObj, err := utils.ExtractValuesObj(setValuesFlag, valuesYamlFile)
	if err != nil {
		return err
	}

	// override the clusterId
	clusterID, ok := valuesFileObj["clusterId"]
	if !ok {
		valuesFileObj["clusterId"] = "default" + ksuid.New().String()
	} else {
		value, ok := clusterID.(string)
		if !ok {
			return fmt.Errorf("clusterId should be of type string got (%v)", reflect.TypeOf(clusterID))
		}
		valuesFileObj["clusterId"] = value + ksuid.New().String()
	}

	return utils.HelmInstall(spaceCloudChart, chartLocation, model.HelmSpaceCloudChartDownloadURL, "", valuesFileObj)
}
