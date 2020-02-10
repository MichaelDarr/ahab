package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause all processes within container",
	Long: `Unpause all processes within container

Docker Command:
  docker unpause CONTAINER
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "unpause", nil)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(unpauseCmd)
}
