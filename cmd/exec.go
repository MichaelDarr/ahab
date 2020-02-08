package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var execCmd = &cobra.Command{
	Use:                "exec",
	Short:              "Execute a command in the container",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.DockerUp(config, configPath)
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "exec", &args)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
