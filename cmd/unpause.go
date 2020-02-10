package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause all processes within container",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "unpause", `Unpause all processes within container

Docker Command:
  docker unpause CONTAINER [ARG...]
				
Usage:
  `+internal.CmdName+` unpause [-h/--help] [ARG...]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "unpause", &args)
		internal.PrintErrFatal(err)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(unpauseCmd)
}
