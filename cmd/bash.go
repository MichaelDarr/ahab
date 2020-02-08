package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Open a containerized bash shell",
	Long: `Attach a containerized bash shell to the active console.

*Warning!* bash must be installed in your image for this command to function!`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.UpContainer(config, configPath)
		internal.PrintErrFatal(err)

		execArgs := []string{"exec", "-it", internal.ContainerName(config, configPath), "bash"}
		err = internal.DockerCmd(&execArgs)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(bashCmd)
}
