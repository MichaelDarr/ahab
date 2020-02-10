package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Inspect changes to files or directories on the container’s filesystem",
	Long: `Inspect changes to files or directories on the container’s filesystem

Docker Command:
  docker diff CONTAINER
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "diff", nil)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
