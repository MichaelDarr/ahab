package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill the container",
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "kill", `Kill the container

Docker Command:
  docker kill [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` kill [-h/--help] [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "kill", &args)
		internal.PrintErrFatal(err)
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(killCmd)
}