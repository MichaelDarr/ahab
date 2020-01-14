package main

import (
	"fmt"
	"log"

	"github.com/MichaelDarr/docker-config/dfg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dfg",
	Short: "dfg is a Docker configuration tool",
	Long:  "Easily configure Docker environments with dfg!",
	Run: func(cmd *cobra.Command, args []string) {
		dockerConfig, err := dfg.LoadDockerConfig("dfg.json")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dockerConfig.Image)
	},
}

func main() {
	rootCmd.Execute()
}
