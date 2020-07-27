package logs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
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
			err = viper.BindPFlag("task-id", cmd.Flags().Lookup("task-id"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('task-id')", nil)
			}
			err = viper.BindPFlag("follow", cmd.Flags().Lookup("follow"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('follow')", nil)
			}
		},
		RunE: actionGetServiceLogs,
	}
	getServiceLogs.Flags().StringP("project", "", "", "The unique id for the project")
	getServiceLogs.Flags().StringP("task-id", "", "", "The unique id for the task")
	getServiceLogs.Flags().BoolP("follow", "", false, "Follow log output")

	getServiceLogs.Flags().StringP("replica-id", "", "", "The unique id for the replica")
	return []*cobra.Command{getServiceLogs}
}

func actionGetServiceLogs(cmd *cobra.Command, args []string) error {
	project, check := utils.GetProjectID()
	if !check {
		_ = utils.LogError("Project not specified in flag", nil)
		return nil
	}
	taskID := viper.GetString("task-id")
	if taskID == "" {
		taskID = project
	}

	if err := GetServiceLogs(project, taskID, viper.GetString("replica-id"), viper.GetBool("follow")); err != nil {
		return nil
	}

	return nil
}
