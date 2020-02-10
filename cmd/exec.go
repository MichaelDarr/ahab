package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a command in a detatched command",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			if arg == "-h" || arg == "--help" {
				err := internal.PrintDockerHelp("exec", `Usage:
  `+internal.CmdName+` exec COMMAND [ARG...]

Docker Command:
  docker exec CONTAINER COMMAND [ARG...]
`)
				internal.PrintErrFatal(err)
				return
			}
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
