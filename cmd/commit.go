package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

// TODO: support argument to tag the commit
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new image from the container's changes",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "commit", `Create a new image from the container's changes

Docker Command:
  docker commit [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` commit [-h/--help] [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		containerExists, err := internal.ContainerExists(config, configPath)
		internal.PrintErrFatal(err)
		if !containerExists {
			internal.PrintErrStr("Container is not created, cannot commit changes.")
		} else {
			err = internal.DockerContainerCmd(config, configPath, "commit", &args)
			internal.PrintErrFatal(err)
		}
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
