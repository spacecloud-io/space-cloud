package operations

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Commands is the list of commands the operations module exposes
func Commands() []*cobra.Command {
	clusterNameAutoComplete := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			utils.LogDebug("Unable to initialize docker client ", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		connArr, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud"), filters.Arg("label", "service=gateway"))})
		if err != nil {
			utils.LogDebug("Unable to list space cloud containers ", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		accountIDs := []string{}
		for _, v := range connArr {
			arr := strings.Split(strings.Split(v.Names[0], "--")[0], "-")
			if len(arr) != 4 {
				// default gateway container
				continue
			}
			accountIDs = append(accountIDs, arr[2])
		}
		return accountIDs, cobra.ShellCompDirectiveDefault
	}

	var setup = &cobra.Command{
		Use:   "setup",
		Short: "setup development environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("local-chart-dir", cmd.Flags().Lookup("local-chart-dir"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('local-chart-dir')", nil)
			}
			if err := viper.BindPFlag("file", cmd.Flags().Lookup("file")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", nil)
			}
			if err := viper.BindPFlag("set", cmd.Flags().Lookup("set")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('set')", nil)
			}
			if err := viper.BindPFlag("version", cmd.Flags().Lookup("version")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('version')", nil)
			}
			if err := viper.BindPFlag("get-defaults", cmd.Flags().Lookup("get-defaults")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('get-defaults')", nil)
			}
		},
		RunE: actionSetup,
	}

	setup.Flags().StringP("version", "v", "", "Space cloud version to use for setup, default to space cli version")
	err := viper.BindEnv("version", "VERSION")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('version') to environment variables", nil)
	}

	setup.Flags().BoolP("get-defaults", "", false, "Prints the default values of cluster config yaml file")
	err = viper.BindEnv("get-defaults", "GET_DEFAULT")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('get-defaults') to environment variables", nil)
	}

	setup.Flags().StringP("local-chart-dir", "c", "", "Path to the space cloud helm chart directory")
	err = viper.BindEnv("local-chart-dir", "LOCAL_CHART_DIR")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('local-chart-dir') to environment variables", nil)
	}

	setup.Flags().StringP("file", "f", "", "Path to the cluster config yaml file")
	err = viper.BindEnv("file", "FILE")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('file' to environment variables", nil)
	}

	setup.Flags().StringP("set", "", "", "Set root string values of chart in format foo1=bar1,foo2=bar2")
	err = viper.BindEnv("`set`", "SET")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('`SET`' to environment variables", nil)
	}

	var update = &cobra.Command{
		Use:   "update",
		Short: "updates the existing space cloud cluster",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("local-chart-dir", cmd.Flags().Lookup("local-chart-dir"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('local-chart-dir')", nil)
			}
			if err := viper.BindPFlag("file", cmd.Flags().Lookup("file")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", nil)
			}
			if err := viper.BindPFlag("set", cmd.Flags().Lookup("set")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('set')", nil)
			}
			if err := viper.BindPFlag("version", cmd.Flags().Lookup("version")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('version')", nil)
			}
		},
		RunE: actionUpdate,
	}

	update.Flags().StringP("local-chart-dir", "c", "", "Path to the space cloud helm chart directory")
	err = viper.BindEnv("local-chart-dir", "LOCAL_CHART_DIR")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('local-chart-dir') to environment variables", nil)
	}

	update.Flags().StringP("version", "v", "", "Space cloud version to use for setup, default to space cli version")
	err = viper.BindEnv("version", "VERSION")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('version') to environment variables", nil)
	}

	update.Flags().StringP("file", "f", "", "Path to the config yaml file")
	err = viper.BindEnv("file", "FILE")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('file' to environment variables", nil)
	}

	update.Flags().StringP("set", "", "", "Set root string values of chart in format foo1=bar1,foo2=bar2")
	err = viper.BindEnv("`set`", "SET")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('`SET`' to environment variables", nil)
	}

	var destroy = &cobra.Command{
		Use:   "destroy",
		Short: "Remove the space cloud cluster from kubernetes",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE: actionDestroy,
	}

	var list = &cobra.Command{
		Use:    "list",
		Short:  "List space-cloud clusters",
		PreRun: nil,
		RunE:   actionList,
	}

	var inspect = &cobra.Command{
		Use:    "inspect",
		Short:  "View applied config file for space cloud cluster",
		PreRun: nil,
		RunE:   actionInspect,
	}

	var apply = &cobra.Command{
		Use:   "apply",
		Short: "Applies a config file or directory",
		RunE:  actionApply,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("delay", cmd.Flags().Lookup("delay")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('delay')", err)
			}
			if err := viper.BindPFlag("force", cmd.Flags().Lookup("force")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('force')", err)
			}
			if err := viper.BindPFlag("file", cmd.Flags().Lookup("file")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", err)
			}
			if err := viper.BindPFlag("retry", cmd.Flags().Lookup("retry")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('retry')", err)
			}
		},
	}
	apply.Flags().DurationP("delay", "", time.Duration(0), "Adds a delay between 2 subsequent request made by space cli to space cloud")
	apply.Flags().BoolP("force", "", false, "Doesn't show warning prompts if some risky changes are made to the config")
	apply.Flags().StringP("file", "f", "", "Path to the resource yaml file or directory")
	apply.Flags().IntP("retry", "r", 1, "Number of retries in case of failure")
	err = viper.BindEnv("file", "FILE")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('file') to environment variables", nil)
	}

	var start = &cobra.Command{
		Use:   "start",
		Short: "Resumes the space-cloud docker environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE: actionStart,
	}
	start.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	if err := start.RegisterFlagCompletionFunc("cluster-name", clusterNameAutoComplete); err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('project')", nil)
	}

	var stop = &cobra.Command{
		Use:   "stop",
		Short: "Stops the space-cloud docker environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE: actionStop,
	}
	stop.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	if err := stop.RegisterFlagCompletionFunc("cluster-name", clusterNameAutoComplete); err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('project')", nil)
	}
	return []*cobra.Command{setup, list, update, inspect, destroy, apply, start, stop}

}

func actionUpdate(cmd *cobra.Command, args []string) error {
	chartDir := viper.GetString("local-chart-dir")
	valuesYamlFile := viper.GetString("file")
	setValue := viper.GetString("set")
	version := viper.GetString("version")
	if version == "" {
		version = model.Version
	}
	return Update(setValue, valuesYamlFile, chartDir, version)
}

func actionSetup(cmd *cobra.Command, args []string) error {
	chartDir := viper.GetString("local-chart-dir")
	valuesYamlFile := viper.GetString("file")
	setValue := viper.GetString("set")
	version := viper.GetString("version")
	if version == "" {
		version = model.Version
	}
	isGetDefaults := viper.GetBool("get-defaults")

	return Setup(setValue, valuesYamlFile, chartDir, version, isGetDefaults)
}

func actionInspect(cmd *cobra.Command, args []string) error {
	clusterID := ""
	if len(args) == 1 {
		clusterID = args[0]
	}
	return Inspect(clusterID)
}

func actionList(cmd *cobra.Command, args []string) error {
	return List()
}

func actionDestroy(cmd *cobra.Command, args []string) error {
	return Destroy()
}

func actionApply(cmd *cobra.Command, args []string) error {
	delay := viper.GetDuration("delay")
	isForce := viper.GetBool("force")
	retry := viper.GetInt("retry")
	var dirName string
	file := viper.GetString("file")
	if file == "" {
		if len(args) > 0 {
			dirName = args[0]
		}
		if dirName == "" {
			return fmt.Errorf("provide the path for spec file or directory using -f flag")
		}
	} else {
		dirName = file
	}
	return Apply(dirName, isForce, delay, retry)
}

func actionStart(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return DockerStart(clusterName)
}

func actionStop(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return DockerStop(clusterName)
}
