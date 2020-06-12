package main

import (
	"os"

	"github.com/spaceuptech/space-cli/cmd"
)

func main() {
	// rootCmd, err := getModule()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	//
	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	if err := cmd.GetRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
}
