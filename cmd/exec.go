package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a command in a container",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "exec", `Run a command in a container

Docker Command:
  docker exec CONTAINER COMMAND [ARG...]

Usage:
  `+internal.CmdName+` exec [-h/--help] COMMAND [ARG...]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.UpContainer(config, configPath)
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "exec", &args)
		internal.PrintErrFatal(err)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(execCmd)
}
