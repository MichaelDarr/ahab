package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start container",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)
		internal.PrintErrFatal(internal.UpContainer(config, configPath))
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
