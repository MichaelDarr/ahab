package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove container",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
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
