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
		status, err := container.Status()
		if err == nil && (status == 2 || status == 3 || status == 5) {
			internal.PrintErrFatal(container.Cmd("stop"))
			internal.PrintErrFatal(container.Cmd("rm"))
		} else if err == nil && (status == 1 || status == 6 || status == 7) {
			internal.PrintErrFatal(container.Cmd("rm"))
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
