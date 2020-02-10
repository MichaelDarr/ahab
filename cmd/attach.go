package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach local standard input, output, and error streams to container",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "attach", `Attach local standard input, output, and error streams to container

Docker Command:
  docker attach [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` attach [-h/--help] [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "attach", &args)
		internal.PrintErrFatal(err)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(attachCmd)
}
