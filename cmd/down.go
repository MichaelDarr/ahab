package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove an environment",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.Config()
		internal.PrintErrFatal(err)
		status, err := internal.ContainerStatus(config, configPath)
		internal.PrintErrFatal(err)

		if status == 0 {
			internal.PrintWarning(1, "Container "+internal.ContainerName(config, configPath)+" is already down")
			return
		} else if status == 4 {
			internal.PrintErrStr("Container " + internal.ContainerName(config, configPath) + " is being removed")
			return
		} else if status == 2 || status == 3 || status == 5 {
			internal.DockerContainerCmd(config, configPath, "stop")
		}

		internal.DockerContainerCmd(config, configPath, "rm")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
