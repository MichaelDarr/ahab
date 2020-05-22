package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
	ahab "github.com/MichaelDarr/ahab/pkg"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the container",
	Long: `Create the container

Docker Command:
  docker create -t [options from ahab.json]`,
	Run: func(cmd *cobra.Command, args []string) {
		container, err := internal.GetContainer()
		ahab.PrintErrFatal(err)

		err = container.Create(false)
		ahab.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
