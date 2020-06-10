package logs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GetSubCommands is the list of commands the log module expose
func GetSubCommands() []*cobra.Command {
	var getServiceLogs = &cobra.Command{
		Use: "service-logs",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("project", cmd.Flags().Lookup("project"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('project')", nil)
			}
			err = viper.BindPFlag("service-id", cmd.Flags().Lookup("service-id"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('service-id')", nil)
			}
			err = viper.BindPFlag("task-id", cmd.Flags().Lookup("task-id"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('task-id')", nil)
			}
			err = viper.BindPFlag("replica-id", cmd.Flags().Lookup("replica-id"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('replica-id')", nil)
			}
		},
		RunE: actionGetServiceLogs,
	}
	getServiceLogs.Flags().StringP("project", "", "", "The unique id for the project")
	getServiceLogs.Flags().StringP("service-id", "", "", "The unique id for the service")
	getServiceLogs.Flags().StringP("task-id", "", "", "The unique id for the task")
	getServiceLogs.Flags().StringP("replica-id", "", "", "The unique id for the replica")

	return []*cobra.Command{getServiceLogs}
}

func actionGetServiceLogs(cmd *cobra.Command, args []string) error {
	project, check := utils.GetProjectID()
	if !check {
		_ = utils.LogError("Project not specified in flag", nil)
		return nil
	}
	serviceID := viper.GetString("service-id")
	taskID := viper.GetString("task-id")
	replicaID := viper.GetString("replica-id")

	if err := GetServiceLogs(project, serviceID, taskID, replicaID); err != nil {
		return nil
	}

	return nil
}
