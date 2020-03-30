package main

import (
	"fmt"
	"os"
	"plugin"

	"github.com/spf13/cobra"
)

func getplugin() (*cobra.Command, error) {
	mod := fmt.Sprintf("%s/cmd_%s.so", getSpaceCliDirectory(), latestVersion)
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = plug.Lookup("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// //var greeter Greeter
	// greeter, ok := symGreeter.(Greeter)
	// if !ok {
	// 	fmt.Println("unexpected type from module symbol")
	// 	os.Exit(1)
	// }
	return nil, nil
}
