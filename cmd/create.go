package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/docker-config/internal"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the container",
	Long: `Create the container

Docker Command:
  docker create [CONFIG_OPTIONS]
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := internal.ProjectConfig()
		internal.PrintErrFatal(err)

		launchOpts, err := internal.LaunchOpts(config, configPath)
		internal.PrintErrFatal(err)

		launchOpts = append([]string{"create", "-t"}, launchOpts...)
		err = internal.DockerCmd(&launchOpts)
		internal.PrintErrFatal(err)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
