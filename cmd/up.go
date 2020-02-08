package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start an environment",
	Run: func(cmd *cobra.Command, args []string) {
		config, _, err := internal.Config()
		internal.PrintErrFatal(err)
		internal.Docker("run", "-td", config.ImageURI)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
