package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var lsvCmd = &cobra.Command{
	Use:   "lsv",
	Short: "List volumes",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ListVolumes()
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(lsvCmd)
}
