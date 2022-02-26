package main

import (
	"fmt"
	"os"

	"github.com/spacecloud-io/space-cloud/cmd/space-cli/commands"

	// Import all apps, modules & managers
	_ "github.com/spacecloud-io/space-cloud/managers"
	_ "github.com/spacecloud-io/space-cloud/modules"
)

func main() {
	if err := commands.GetRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
