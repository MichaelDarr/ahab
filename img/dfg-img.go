package main

import (
	"fmt"
	"os"
)

func argPresent(searchArg string) bool {
	for _, curArg := range os.Args {
		if curArg == searchArg {
			return true
		}
	}
	return false
}

func main() {

	if len(os.Args) < 2 || argPresent("-h") || argPresent("--help") {
		fmt.Print(`usage: dfg-img [OPTIONS] [IMAGE]
Create a dfg-style Dockerfile for image IMAGE

Options
	-f, --family        Distribution family
		debian (apt-get)
		more coming soon!
	-h, --help          Display this help menu

report any bugs at https://github.com/MichaelDarr/docker-config/issues
`)
		os.Exit(0)
	}

	dockerfileHeader := "FROM " + os.Args[1]

	fmt.Println(dockerfileHeader)
}
