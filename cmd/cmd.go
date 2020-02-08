package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var cmdCmd = &cobra.Command{
	Use:                "cmd",
	Short:              "Execute an attatched command in the container",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.DockerUp(config, configPath)
		internal.PrintErrFatal(err)

		execArgs := []string{"exec", "-it", internal.ContainerName(config, configPath)}
		execArgs = append(execArgs, args...)
		err = internal.DockerCmd(&execArgs)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(cmdCmd)
}
