package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

// Execute is the entrypoint to the dcfg CLI
func Execute() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:     "dcon",
	Short:   "dcon is a Docker configuration tool",
	Long:    "Configure, launch, and work in Dockerized environments",
	Version: internal.Version,
}

// verbose is used as a flag for various commands
var verbose bool

// Docker commands that don't take options
var diffCmd = BasicCommand("diff", "Inspect changes to files or directories on container filesystem")
var pauseCmd = BasicCommand("pause", "Pause all processes within container")
var unpauseCmd = BasicCommand("unpause", "Unpause all processes within container")

// Docker commands that take options
var attachCmd = OptionCommand("attach", "Attach local standard input, output, and error streams to container")
var commitCmd = OptionCommand("commit", "Create a new image from container's changes")
var exportCmd = OptionCommand("export", "Export container’s filesystem as a tar archive")
var killCmd = OptionCommand("kill", "Kill container")
var rmCmd = OptionCommand("rm", "Remove container")
var startCmd = OptionCommand("start", "Start stopped container")
var stopCmd = OptionCommand("stop", "Stop running container")

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
	rootCmd.AddCommand(&pauseCmd)
	rootCmd.AddCommand(&rmCmd)
	rootCmd.AddCommand(&shCmd)
	rootCmd.AddCommand(&startCmd)
	rootCmd.AddCommand(&stopCmd)
	rootCmd.AddCommand(&unpauseCmd)
	rootCmd.AddCommand(&zshCmd)
}
