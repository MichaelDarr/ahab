package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List containers, images, and volumes",
	Long: `List containers, images, and volumes

Docker Commands:
  docker ps -a [FORMATTING FLAGS]
  docker images [FORMATTING FLAGS]
  docker volume ls
`,
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
	lsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full Docker resource info")

	rootCmd.AddCommand(lsCmd)
}
