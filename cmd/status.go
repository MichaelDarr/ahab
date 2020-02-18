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
		container, err := internal.GetContainer()
		internal.PrintErrFatal(err)

		fmt.Println("Container:")
		internal.PrintIndentedPair("Name", container.Name())

		statusCode, err := container.Status()
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

		containerImage, err := container.Prop("Config.Image")
		internal.PrintErrFatal(err)
		if containerImage != "" {
			internal.PrintIndentedPair("Image", containerImage)
		}

		containerID, err := container.Prop("Id")
		internal.PrintErrFatal(err)
		if containerID != "" {
			internal.PrintIndentedPair("ID", containerID)
		}

		internal.PrintIndentedPair("Config", container.FilePath)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
