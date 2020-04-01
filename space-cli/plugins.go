package main

import (
	"fmt"
	"plugin"

	"github.com/spf13/cobra"
)

type GetRootCommand interface {
	Rootcommand() *cobra.Command
}

func getplugin(latestversion string) (*cobra.Command, error) {
	mod := fmt.Sprintf("%s/cmd_%s.so", getSpaceCLIDirectory(), latestversion)
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}
	commands, err := plug.Lookup("GetRootCommand")
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
