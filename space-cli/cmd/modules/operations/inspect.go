package operations

import (
	"fmt"

	"github.com/ghodss/yaml"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Inspect prints values file used for setting up cluster
func Inspect(clusterID string) error {
	if clusterID == "" {
		charList, err := utils.HelmList(model.HelmSpaceCloudNamespace)
		if err != nil {
			return err
		}
		if len(charList) < 1 {
			utils.LogInfo("space cloud cluster not found, setup a new cluster using the setup command")
			return nil
		}
		clusterID = charList[0].Name
	}

	chartInfo, err := utils.HelmGet(clusterID)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(chartInfo.Config)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
