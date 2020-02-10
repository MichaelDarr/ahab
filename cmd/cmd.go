package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Execute an attatched command in a container",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "exec", `Execute an attatched command in a container

Docker Command:
  docker exec -it CONTAINER COMMAND [ARG...]
			  
Usage:
  `+internal.CmdName+` cmd [-h/--help] COMMAND [ARG...]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)

		err = internal.UpContainer(config, configPath)
		internal.PrintErrFatal(err)

		execArgs := []string{"exec", "-it", internal.ContainerName(config, configPath)}
		execArgs = append(execArgs, args...)
		err = internal.DockerCmd(&execArgs)
		internal.PrintErrFatal(err)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(cmdCmd)
}
