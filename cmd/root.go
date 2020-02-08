package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var rootCmd = &cobra.Command{
	Use:     "dcfg",
	Short:   "dcfg is a Docker configuration tool",
	Long:    "A tool to configure and run Dockerized environments with ease.",
	Version: internal.Version,
}

// Execute is the entrypoint to the dcfg CLI
func Execute() {
	rootCmd.Execute()
}
