package main

import (
	"github.com/spf13/cobra"
)

// Build-time variables
var version string // dcfg version

// Cobra entrypoint
var rootCmd = &cobra.Command{
	Use:     "dcfg",
	Short:   "docker-config is a Docker configuration tool",
	Long:    "Easily configure Docker environments with the dcfg CLI tool!",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := Config()
		PrintErrFatal(err)

		_, err = LaunchOptions(config, configPath)
		PrintErrFatal(err)
	},
}

// Program entrypoint
func main() {
	rootCmd.Execute()
}
