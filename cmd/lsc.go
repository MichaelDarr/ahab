package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var lscCmd = &cobra.Command{
	Use:   "lsc",
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListContainers(verbose)
		internal.PrintErrFatal(err)
	},
}

func init() {
	lscCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "View full container info")

	rootCmd.AddCommand(lscCmd)
}
