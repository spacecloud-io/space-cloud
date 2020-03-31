package main

import (
	"fmt"
	"os"
	"plugin"

	"github.com/spf13/cobra"
)

type Cmd interface {
	Getcommand() *cobra.Command
}

func getplugin() (*cobra.Command, error) {
	mod := fmt.Sprintf("./cmd/cmd_%s.so", LatestVersion)
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}
	commands, err := plug.Lookup("Cmd")
	if err != nil {
		return nil, err
		os.Exit(1)
	}

	var cmd Cmd
	cmd, ok := commands.(Cmd)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		return nil, nil
		os.Exit(1)
	}

	rootCmd := cmd.Getcommand()
	return rootCmd, nil

}
