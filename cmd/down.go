package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove container",
	Long: `Stop and remove container

Docker Command:
  varies based on current container state`,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.StopContainer(config, configPath)
		internal.PrintErrFatal(err)

		err = internal.RemoveContainer(config, configPath)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
