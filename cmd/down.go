package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove container",
	Run: func(cmd *cobra.Command, args []string) {
		container, err := internal.GetContainer()
		internal.PrintErrFatal(err)
		internal.PrintErrFatal(container.Down())
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
