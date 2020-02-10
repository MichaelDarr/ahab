package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List containers, images, and volumes",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListContainers(verbose)
		internal.PrintErrFatal(err)
		err = internal.ListImages(verbose)
		internal.PrintErrFatal(err)
		err = internal.ListVolumes()
		internal.PrintErrFatal(err)
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full info")

	rootCmd.AddCommand(lsCmd)
}
