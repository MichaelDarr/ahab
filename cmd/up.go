package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
	ahab "github.com/MichaelDarr/ahab/pkg"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start container",
	Run: func(cmd *cobra.Command, args []string) {
		container, err := internal.GetContainer()
		ahab.PrintErrFatal(err)
		ahab.PrintErrFatal(container.Up())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
