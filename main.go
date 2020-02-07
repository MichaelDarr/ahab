package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "dcfg",
	Short:   "docker-config is a Docker configuration tool",
	Long:    "Easily configure Docker environments with the dcfg CLI tool!",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		err := InitConfig()
		PrintErrFatal(err)

		containerName, err := ContainerName()
		PrintErrFatal(err)
		println(containerName)
	},
}

func main() {
	rootCmd.Execute()
}
