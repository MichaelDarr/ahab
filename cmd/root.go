package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:     "dcon",
	Short:   "dcon is a Docker configuration tool",
	Long:    "Configure, launch, and work in Dockerized environments",
	Version: internal.Version,
}

// Execute is the entrypoint to the dcfg CLI
func Execute() {
	rootCmd.Execute()
}
