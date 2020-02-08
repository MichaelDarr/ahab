package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var zshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Open a containerized zsh shell",
	Long: `Attach a containerized zsh shell to the active console.

*Warning!* zsh must be installed in your image for this command to function!`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.UpContainer(config, configPath)
		internal.PrintErrFatal(err)

		execArgs := []string{"exec", "-it", internal.ContainerName(config, configPath), "zsh"}
		err = internal.DockerCmd(&execArgs)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(zshCmd)
}
