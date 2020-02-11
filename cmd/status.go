package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print container status",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		fmt.Println("Container:")
		internal.PrintIndentedPair("Name", internal.ContainerName(config, configPath))

		statusCode, err := internal.ContainerStatus(config, configPath)
		internal.PrintErrFatal(err)
		switch statusCode {
		case 0:
			internal.PrintIndentedPair("Status", "Not Created")
		case 1:
			internal.PrintIndentedPair("Status", "Created")
		case 2:
			internal.PrintIndentedPair("Status", "Restarting")
		case 3:
			internal.PrintIndentedPair("Status", "Running")
		case 4:
			internal.PrintIndentedPair("Status", "Removing")
		case 5:
			internal.PrintIndentedPair("Status", "Paused")
		case 6:
			internal.PrintIndentedPair("Status", "Exited")
		case 7:
			internal.PrintIndentedPair("Status", "Dead")
		default:
			internal.PrintIndentedPair("Status", "Unknown")
		}

		containerImage, err := internal.ContainerProp(config, configPath, "Config.Image")
		internal.PrintErrFatal(err)
		if containerImage != "" {
			internal.PrintIndentedPair("Image", containerImage)
		}

		containerID, err := internal.ContainerProp(config, configPath, "Id")
		internal.PrintErrFatal(err)
		if containerID != "" {
			internal.PrintIndentedPair("ID", containerID)
		}

		internal.PrintIndentedPair("Config", configPath)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
