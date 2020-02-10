package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Execute an attatched command",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// The way we execute the exec command, Docker interprets the --help flag as part of the
		// command 'exec' is trying to run. Here, we trigger it manually. This prevents the user
		// from intentionally passing help flags into exec, which may be an issue down the road.
		for _, arg := range args {
			if arg == "-h" || arg == "--help" {
				helpArgs := []string{"exec", "--help"}
				err := internal.DockerCmd(&helpArgs)
				internal.PrintErrFatal(err)
				return
			}
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
