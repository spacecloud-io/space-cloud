package logs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GetSubCommands is the list of commands the log module expose
func GetSubCommands() []*cobra.Command {
	var getServiceLogs = &cobra.Command{
		Use:     "logs [replica-id]",
		Example: "1) space-cli logs service1--v1 --project myproject --follow\n2) space-cli logs service1--v1 --project myproject --follow --task-id greeting",
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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				project, check := utils.GetProjectID()
				if !check {
					utils.LogDebug("Project not specified in flag", nil)
					return nil, cobra.ShellCompDirectiveDefault
				}
				replicaIDs, err := getServiceStatus(project, "service-status", map[string]string{})
				if err != nil {
					utils.LogDebug("Unable to get service status", map[string]interface{}{"error": err})
					return nil, cobra.ShellCompDirectiveDefault
				}
				return replicaIDs, cobra.ShellCompDirectiveDefault
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
	}

	getServiceLogs.Flags().StringP("task-id", "", "", "The unique id for the task")
	getServiceLogs.Flags().BoolP("follow", "", false, "Follow log output")
	return []*cobra.Command{getServiceLogs}
}

func actionGetServiceLogs(cmd *cobra.Command, args []string) error {
	project, check := utils.GetProjectID()
	if !check {
		_ = utils.LogError("Project not specified in flag", nil)
		return nil
	}
	if len(args) == 0 {
		_ = utils.LogError("Replica name not provide", nil)
		return nil
	}
	replicaID := args[0]

	if err := GetServiceLogs(project, viper.GetString("task-id"), replicaID, viper.GetBool("follow")); err != nil {
		return err
	}
	return nil
}
