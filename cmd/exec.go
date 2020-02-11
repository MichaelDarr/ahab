package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a command in a container",
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "exec", `Run a command in a container

Docker Command:
  docker exec CONTAINER COMMAND [ARG...]

Usage:
  ahab exec [-h/--help] COMMAND [ARG...]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.UpContainer(config, configPath)
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "exec", &args)
		internal.PrintErrFatal(err)
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(execCmd)
}
