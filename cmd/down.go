package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove container",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)
		internal.PrintErrFatal(internal.StopContainer(config, configPath))
		internal.PrintErrFatal(internal.RemoveContainer(config, configPath))
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
