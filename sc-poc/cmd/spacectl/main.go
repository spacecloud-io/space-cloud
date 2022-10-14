package main

import (
	"fmt"
	"os"

	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands"

	// Import all apps, modules & managers
	_ "github.com/spacecloud-io/space-cloud/managers"
	_ "github.com/spacecloud-io/space-cloud/modules"

	_ "github.com/spacecloud-io/space-cloud/pkg/client/clientset/versioned"
)

func main() {
	if err := commands.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
