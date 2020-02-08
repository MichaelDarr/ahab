package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start an environment",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		launchOpts := append([]string{"run", "-td"}, internal.LaunchOpts(config, configPath)...)
		internal.DockerCmd(&launchOpts)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
