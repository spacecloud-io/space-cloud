package pkg

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	// Add all sub commands
	cmd.AddCommand(newCommandInitialize())
	cmd.AddCommand(newCommandApply())
	cmd.AddCommand(newCommandGet())

	return cmd
}
