package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start stopped container",
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "start", `Start stopped container

Docker Command:
  docker start [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` start [-h/--help] [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "start", &args)
		internal.PrintErrFatal(err)
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(startCmd)
}
