package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var lsiCmd = &cobra.Command{
	Use:   "lsi",
	Short: "List images",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListImages(verbose)
		internal.PrintErrFatal(err)
	},
}

func init() {
	lsiCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full image info")

	rootCmd.AddCommand(lsiCmd)
}
