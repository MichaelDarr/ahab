package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove an environment",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)
		status, err := internal.ContainerStatus(config, configPath)
		internal.PrintErrFatal(err)
		println(status)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
