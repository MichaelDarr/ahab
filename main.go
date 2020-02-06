package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dfg",
	Short: "dfg is a Docker configuration tool",
	Long:  "Easily configure Docker environments with dfg!",
	Run: func(cmd *cobra.Command, args []string) {
		err := InitConfig()
		PrintErrFatal(err)
	},
}

func main() {
	rootCmd.Execute()
}
