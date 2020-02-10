package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

// TODO: support argument to tag the commit
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the container’s filesystem as a tar archive",
	Run: func(cmd *cobra.Command, args []string) {
		helpRequested, err := internal.PrintDockerHelp(&args, "export", `Export the container’s filesystem as a tar archive

Docker Command:
  docker export [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` export [-h/--help] [OPTIONS]
`)
		internal.PrintErrFatal(err)
		if helpRequested {
			return
		}
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.DockerContainerCmd(config, configPath, "export", &args)
		internal.PrintErrFatal(err)
	},
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
