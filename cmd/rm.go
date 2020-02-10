package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove container",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "rm", `Remove container

Docker Command:
  docker rm CONTAINER [ARG...]
				
Usage:
  `+internal.CmdName+` rm [-h/--help] [ARG...]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "rm", &args)
		internal.PrintErrFatal(err)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
