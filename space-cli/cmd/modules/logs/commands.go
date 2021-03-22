package logs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/services"
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
			err = viper.BindPFlag("since", cmd.Flags().Lookup("since"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('since')", nil)
			}
			err = viper.BindPFlag("since-time", cmd.Flags().Lookup("since-time"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('since-time')", nil)
			}
			err = viper.BindPFlag("tail", cmd.Flags().Lookup("tail"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('tail')", nil)
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
	getServiceLogs.Flags().StringP("since", "", "", "Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of\nsince-time / since may be used.")
	getServiceLogs.Flags().StringP("since-time", "", "", "Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time /\nsince may be used.")
	getServiceLogs.Flags().StringP("tail", "", "", "Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise")

	if err := getServiceLogs.RegisterFlagCompletionFunc("task-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		specObjects, err := services.GetServices(project, "", map[string]string{})
		if err != nil {
			utils.LogDebug("Unable to get services from server", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		tasksArr := make([]string, 0)
		for _, object := range specObjects {
			obj, ok := object.Spec.(map[string]interface{})
			if !ok {
				continue
			}
			tasks, ok := obj["tasks"].([]interface{})
			if !ok {
				continue
			}
			for _, task := range tasks {
				taskID, ok := task.(map[string]interface{})["id"]
				if !ok {
					continue
				}
				tasksArr = append(tasksArr, taskID.(string))
			}
		}
		return tasksArr, cobra.ShellCompDirectiveDefault
	}); err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('project')", nil)
	}
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
