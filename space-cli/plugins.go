package main

import (
	"fmt"
	"plugin"

	"github.com/spf13/cobra"
)

type GetRootCommand interface {
	Rootcommand() *cobra.Command
}

func getplugin(latestVersion string) (*cobra.Command, error) {
	mod := fmt.Sprintf("./cmd/cmd_%s.so", latestVersion)
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}
	commands, err := plug.Lookup("Cmd")
	if err != nil {
		return nil, err
	}

	var getRootCommand GetRootCommand
	getRootCommand, ok := commands.(GetRootCommand)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		return nil, nil
	}

	rootCmd := getRootCommand.Rootcommand()
	return rootCmd, nil

}
