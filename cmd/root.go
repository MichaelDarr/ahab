package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

// rootCmd is the base command which all others are added to
var rootCmd = &cobra.Command{
	Use:     "ahab",
	Short:   "ahab is a Docker configuration tool",
	Long:    "Configure, launch, and work in Dockerized environments with ahab",
	Version: internal.Version,
}

// Execute is the entrypoint to the CLI
func Execute() {
	rootCmd.Execute()
}

// verbose is used as a flag for various commands
var verbose bool

// Docker commands that don't take options
var diffCmd = BasicCommand("diff", "Inspect changes to files or directories on container filesystem")
var pauseCmd = BasicCommand("pause", "Pause all processes within container")
var portCmd = BasicCommand("port", "List port mappings for the container")
var topCmd = BasicCommand("top", "Display the running processes of the container")
var unpauseCmd = BasicCommand("unpause", "Unpause all processes within container")
var waitCmd = BasicCommand("wait", "Block until the container stops, then print its exit code")

// Docker commands that take options
var attachCmd = OptionCommand("attach", "Attach local standard input, output, and error streams to container")
var commitCmd = OptionCommand("commit", "Create a new image from container's changes")
var exportCmd = OptionCommand("export", "Export container’s filesystem as a tar archive")
var killCmd = OptionCommand("kill", "Kill container")
var logsCmd = OptionCommand("logs", "Fetch the container logs")
var restartCmd = OptionCommand("restart", "Restart container")
var rmCmd = OptionCommand("rm", "Remove container")
var startCmd = OptionCommand("start", "Start stopped container")
var statsCmd = OptionCommand("stats", "Display a live stream of container resource usage statistics")
var stopCmd = OptionCommand("stop", "Stop running container")
var updateCmd = OptionCommand("update", "Update configuration of the container")

// Shell attachment commands
var bashCmd = ShellCommand("bash", "bash")
var shCmd = ShellCommand("sh", "bourne")
var zshCmd = ShellCommand("zsh", "z")

// init adds all the generic subcommands
func init() {
	rootCmd.AddCommand(&attachCmd)
	rootCmd.AddCommand(&bashCmd)
	rootCmd.AddCommand(&commitCmd)
	rootCmd.AddCommand(&diffCmd)
	rootCmd.AddCommand(&exportCmd)
	rootCmd.AddCommand(&killCmd)
	rootCmd.AddCommand(&logsCmd)
	rootCmd.AddCommand(&pauseCmd)
	rootCmd.AddCommand(&portCmd)
	rootCmd.AddCommand(&restartCmd)
	rootCmd.AddCommand(&rmCmd)
	rootCmd.AddCommand(&shCmd)
	rootCmd.AddCommand(&startCmd)
	rootCmd.AddCommand(&statsCmd)
	rootCmd.AddCommand(&stopCmd)
	rootCmd.AddCommand(&topCmd)
	rootCmd.AddCommand(&unpauseCmd)
	rootCmd.AddCommand(&updateCmd)
	rootCmd.AddCommand(&waitCmd)
	rootCmd.AddCommand(&zshCmd)
}

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
  ahab `+command+` [-h/--help] [OPTIONS]
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
