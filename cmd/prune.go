package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MichaelDarr/ahab/internal"
	ahab "github.com/MichaelDarr/ahab/pkg"
)

// volumes is used as a flag by the prune command
var volumes bool

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove unused Docker assets",
	Long: `Remove unused Docker assets

Docker Command:
  docker system prune -a [--volumes]`,
	Run: func(cmd *cobra.Command, args []string) {
		pruneArgs := []string{"system", "prune", "-a"}
		if volumes {
			pruneArgs = append(pruneArgs, "--volumes")
		}
		ahab.PrintErrFatal(internal.DockerCmd(&pruneArgs))
	},
}

func init() {
	pruneCmd.Flags().BoolVar(&volumes, "volumes", false, "Also prune docker volumes")
	rootCmd.AddCommand(pruneCmd)
}
