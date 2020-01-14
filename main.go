package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// DfgConfig contains all possible docker config fields
type DfgConfig struct {
	Image string `json:"image"`
}

var rootCmd = &cobra.Command{
	Use:   "dfg",
	Short: "dfg is a Docker configuration tool",
	Long:  "Easily configure Docker environments with dfg!",
	Run: func(cmd *cobra.Command, args []string) {
		dockerConfig, err := LoadDockerConfig("dfg.json")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dockerConfig.Image)
	},
}

func main() {
	rootCmd.Execute()
}
