package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

// BasicCommand constructs and returns a Docker container command which doesn not take extra options
func BasicCommand(command string, description string) cobra.Command {
	return cobra.Command{
		Use:   command,
		Short: description,
		Long: description + `

Docker Command:
  docker ` + command + ` [CONFIG_OPTIONS]
	`,
		Run: func(cmd *cobra.Command, args []string) {
			config, configPath, err := internal.ProjectConfig()
			internal.PrintErrFatal(err)

			err = internal.DockerContainerCmd(config, configPath, command, nil)
			internal.PrintErrFatal(err)
		},
	}
}

// OptionCommand constructs and returns a Docker container command which takes extra options
func OptionCommand(command string, description string) cobra.Command {
	return cobra.Command{
		Use:   command,
		Short: description,
		Run: func(cmd *cobra.Command, args []string) {
			helpRequested, err := internal.PrintDockerHelp(&args, command, description+`

Docker Command:
  docker `+command+` [OPTIONS] CONTAINER

Usage:
  `+internal.CmdName+` `+command+` [-h/--help] [OPTIONS]
`)
			internal.PrintErrFatal(err)
			if helpRequested {
				return
			}
			config, configPath, err := internal.ProjectConfig()
			internal.PrintErrFatal(err)

			err = internal.DockerContainerCmd(config, configPath, command, &args)
			internal.PrintErrFatal(err)
		},
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
	}
}

// ShellCommand constructs and returns a Docker container command to attach a shell to the active terminal
func ShellCommand(shellCommand string, description string) cobra.Command {
	return cobra.Command{
		Use:   shellCommand,
		Short: "Open a containerized " + description + " shell",
		Long: `Attach a containerized ` + description + ` shell to the active terminal.

*Warning!* the ` + description + ` shell must be installed in your image for this command to function!

Docker Command:
  docker exec -it CONTAINER ` + shellCommand,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config, configPath, err := internal.ProjectConfig()
			internal.PrintErrFatal(err)

			err = internal.UpContainer(config, configPath)
			internal.PrintErrFatal(err)

			execArgs := []string{"exec", "-it", internal.ContainerName(config, configPath), shellCommand}
			err = internal.DockerCmd(&execArgs)
			internal.PrintErrFatal(err)
		},
	}
}
