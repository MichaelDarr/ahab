package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause all processes within container",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "pause", `Pause all processes within container

Docker Command:
  docker pause CONTAINER [ARG...]
		
Usage:
  `+internal.CmdName+` pause [-h/--help] [ARG...]
		`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "pause", &args)
		internal.PrintErrFatal(err)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
