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
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		err = internal.CreateContainer(config, configPath, false)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
