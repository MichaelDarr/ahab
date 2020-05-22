package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
	ahab "github.com/MichaelDarr/ahab/pkg"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print container status",
	Run: func(cmd *cobra.Command, args []string) {
		container, err := internal.GetContainer()
		ahab.PrintErrFatal(err)

		fmt.Println("Container:")
		ahab.PrintIndentedPair("Name", container.Name())

		statusCode, err := container.Status()
		ahab.PrintErrFatal(err)
		switch statusCode {
		case 0:
			ahab.PrintIndentedPair("Status", "Not Created")
		case 1:
			ahab.PrintIndentedPair("Status", "Created")
		case 2:
			ahab.PrintIndentedPair("Status", "Restarting")
		case 3:
			ahab.PrintIndentedPair("Status", "Running")
		case 4:
			ahab.PrintIndentedPair("Status", "Removing")
		case 5:
			ahab.PrintIndentedPair("Status", "Paused")
		case 6:
			ahab.PrintIndentedPair("Status", "Exited")
		case 7:
			ahab.PrintIndentedPair("Status", "Dead")
		default:
			ahab.PrintIndentedPair("Status", "Unknown")
		}

		containerImage, err := container.Prop("Config.Image")
		ahab.PrintErrFatal(err)
		if containerImage != "" {
			ahab.PrintIndentedPair("Image", containerImage)
		}

		containerID, err := container.Prop("Id")
		ahab.PrintErrFatal(err)
		if containerID != "" {
			ahab.PrintIndentedPair("ID", containerID)
		}

		ahab.PrintIndentedPair("Config", container.FilePath)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
