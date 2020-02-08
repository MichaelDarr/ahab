package cmd

import (
	"github.com/spf13/cobra"
)

// Version is the build-time dcfg version variable
var Version string

// RootCmd is the Cobra root command
var RootCmd = &cobra.Command{
	Use:     "dcfg",
	Short:   "dcfg is a Docker configuration tool",
	Long:    "A tool to configure and run Dockerized environments with ease.",
	Version: Version,
}

// Execute is the entrypoint to the dcfg CLI
func Execute() {
	RootCmd.Execute()
}
