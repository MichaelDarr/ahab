package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the container",
	Long: `Create the container

Docker Command:
  docker create -t [options from ahab.json]`,
	Run: func(cmd *cobra.Command, args []string) {
		container, err := internal.GetContainer()
		internal.PrintErrFatal(err)

		err = container.Create(false)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
