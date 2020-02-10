package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause all processes within container",
	Long: `Pause all processes within container

Docker Command:
  docker pause CONTAINER
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "pause", nil)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
