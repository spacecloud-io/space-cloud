package operations

import (
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// List initializes development environment
func List() error {
	chartList, err := utils.HelmList(model.HelmSpaceCloudNamespace)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CLUSTER ID", "NAMESPACE", "UPDATED", "STATUS", "CHART", "APP VERSION"})

	table.SetBorder(true)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("|")
	table.SetAutoWrapText(false)
	for _, release := range chartList {
		table.Append([]string{release.Name, release.Namespace, release.Info.LastDeployed.String(), release.Info.Status.String(), release.Chart.Name(), release.Chart.AppVersion()})
	}
	table.Render()
	return nil
}
