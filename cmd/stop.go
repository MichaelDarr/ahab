package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop container",
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "stop", `Stop container

Docker Command:
  docker stop [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` stop [-h/--help] [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "stop", &args)
		internal.PrintErrFatal(err)
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
