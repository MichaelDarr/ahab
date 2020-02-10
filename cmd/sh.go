package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var shCmd = &cobra.Command{
	Use:   "sh",
	Short: "Open a containerized bourne shell",
	Long: `Attach a containerized bourne shell to the active console.

*Warning!* the bourne shell must be installed in your image for this command to function!

Docker Command:
  docker exec -it CONTAINER sh`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.UpContainer(config, configPath)
		internal.PrintErrFatal(err)

		execArgs := []string{"exec", "-it", internal.ContainerName(config, configPath), "sh"}
		err = internal.DockerCmd(&execArgs)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(shCmd)
}
