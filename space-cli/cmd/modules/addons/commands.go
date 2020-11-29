package addons

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Commands is the list of commands the addon module exposes
func Commands() []*cobra.Command {
	var addCmd = &cobra.Command{
		Use:           "add",
		Short:         "Add a add-on to the environment",
		SilenceErrors: true,
	}

	var addDatabaseCmd = &cobra.Command{
		Use:   "database",
		Short: "Add a database",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("local-chart-dir", cmd.Flags().Lookup("local-chart-dir"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('local-chart-dir')", nil)
			}
			if err := viper.BindPFlag("values", cmd.Flags().Lookup("values")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('values')", nil)
			}
			err = viper.BindPFlag("name", cmd.Flags().Lookup("name"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('name')", nil)
			}
			if err := viper.BindPFlag("set", cmd.Flags().Lookup("set")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('set')", nil)
			}

		},
		RunE:      ActionAddDatabase,
		ValidArgs: []string{"mysql", "postgres", "sqlserver", "mongo"},
	}

	addDatabaseCmd.Flags().StringP("name", "", "", "provide the name for the database")

	addDatabaseCmd.Flags().StringP("local-chart-dir", "c", "", "Path to the space cloud helm chart directory")
	err := viper.BindEnv("local-chart-dir", "LOCAL_CHART_DIR")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('local-chart-dir') to environment variables", nil)
	}

	addDatabaseCmd.Flags().StringP("values", "f", "", "Path to the config yaml file")
	err = viper.BindEnv("values", "VALUES")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('values' to environment variables", nil)
	}

	addDatabaseCmd.Flags().StringP("set", "", "", "Set root string values of chart in format foo1=bar1,foo2=bar2")
	err = viper.BindEnv("`set`", "SET")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('`SET`' to environment variables", nil)
	}

	var removeCmd = &cobra.Command{
		Use:           "remove",
		Short:         "Remove a add-on from the environment",
		SilenceErrors: true,
	}

	var removeDatabaseCmd = &cobra.Command{
		Use:    "database",
		Short:  "Remove a database",
		PreRun: func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				ctx := context.Background()
				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					utils.LogDebug("Unable to initialize docker client ", nil)
					return nil, cobra.ShellCompDirectiveDefault
				}
				clusterName := cmd.Flag("cluster-name").Value.String()
				conArr, err := utils.GetContainers(ctx, cli, clusterName, model.DbContainers)
				if err != nil {
					utils.LogDebug("Unable to list database containers ", nil)
					return nil, cobra.ShellCompDirectiveDefault
				}
				dbAlias := make([]string, 0)
				for _, container := range conArr {
					value, ok := container.Labels["name"]
					if !ok {
						continue
					}
					dbAlias = append(dbAlias, value)
				}
				return dbAlias, cobra.ShellCompDirectiveDefault
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		RunE: ActionRemoveDatabase,
	}

	addCmd.AddCommand(addDatabaseCmd)
	removeCmd.AddCommand(removeDatabaseCmd)

	return []*cobra.Command{addCmd, removeCmd}
}

// ActionAddDatabase adds a database add on
func ActionAddDatabase(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return utils.LogError("Database type not provided as an argument", nil)
	}
	dbType := args[0]
	switch dbType {
	case "postgres", "mysql", "sqlserver", "mongo":
	default:
		return fmt.Errorf("unkown database (%s) provided as argument", dbType)
	}

	name := viper.GetString("name")
	if name == "" {
		utils.LogInfo(fmt.Sprintf("--name flag not provided using the name (%s) for database", name))
		name = dbType
	}

	chartDir := viper.GetString("local-chart-dir")
	valuesYamlFile := viper.GetString("values")
	setValue := viper.GetString("set")
	return addDatabase(name, dbType, setValue, valuesYamlFile, chartDir)
}

// ActionRemoveDatabase removes a database add on
func ActionRemoveDatabase(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return utils.LogError("Database name not provided as an argument", nil)
	}
	return removeDatabase(args[0])
}
