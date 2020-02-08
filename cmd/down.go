package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove the container",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.DockerStop(config, configPath)
		internal.PrintErrFatal(err)

		err = internal.DockerRemove(config, configPath)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
