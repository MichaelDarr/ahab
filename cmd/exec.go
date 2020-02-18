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

		container, err := internal.GetContainer()
		internal.PrintErrFatal(err)
		internal.PrintErrFatal(container.Up())

		containerOpts := []string{"exec"}
		if container.Fields.User != "" {
			containerOpts = append(containerOpts, "-u", container.Fields.User)
		} else if !container.Fields.Permissions.Disable {
			containerOpts = append(containerOpts, "-u", internal.ContainerUserName)
		}
		containerOpts = append(containerOpts, container.Name())
		containerOpts = append(containerOpts, args...)
		internal.PrintErrFatal(internal.DockerCmd(&containerOpts))
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(execCmd)
}
