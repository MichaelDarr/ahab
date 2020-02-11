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
  docker exec [-u ahab] CONTAINER COMMAND [ARG...]

Usage:
  ahab exec [-h/--help] COMMAND [ARG...]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)
		internal.PrintErrFatal(internal.UpContainer(config, configPath))

		containerOpts := []string{"exec"}
		if !config.ManualPermissions {
			containerOpts = append(containerOpts, "-u", internal.ContainerUserName)
		}
		containerOpts = append(containerOpts, internal.ContainerName(config, configPath))
		containerOpts = append(containerOpts, args...)
		internal.PrintErrFatal(internal.DockerCmd(&containerOpts))
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(execCmd)
}
