package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Execute an attatched command in a container",
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "exec", `Execute an attatched command in a container

Docker Command:
  docker exec -it CONTAINER COMMAND [OPTIONS]

Usage:
  ahab cmd [-h/--help] COMMAND [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}

		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)
		internal.PrintErrFatal(internal.UpContainer(config, configPath))

		containerOpts := []string{"exec", "-it"}
		if config.User != "" {
			containerOpts = append(containerOpts, "-u", config.User)
		} else if !config.Permissions.Disable {
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
	rootCmd.AddCommand(cmdCmd)
}
