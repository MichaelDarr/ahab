package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
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

var lscCmd = &cobra.Command{
	Use:   "lsc",
	Short: "List containers",
	Long: `List containers

Docker Command:
  docker ps -a [FORMATTING FLAGS]
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListContainers(verbose)
		internal.PrintErrFatal(err)
	},
}

var lsiCmd = &cobra.Command{
	Use:   "lsi",
	Short: "List images",
	Long: `List images

Docker Command:
  docker images [FORMATTING FLAGS]
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListImages(verbose)
		internal.PrintErrFatal(err)
	},
}

var lsvCmd = &cobra.Command{
	Use:   "lsv",
	Short: "List volumes",
	Long: `List volumes

Docker Command:
  docker volume ls
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListVolumes()
		internal.PrintErrFatal(err)
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full Docker resource info")
	lscCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full container info")
	lsiCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full image info")

	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(lscCmd)
	rootCmd.AddCommand(lsiCmd)
	rootCmd.AddCommand(lsvCmd)
}
