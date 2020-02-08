package internal

import "os/exec"

// Docker runs a Docker command
func Docker() {
	cmd := exec.Command("docker", "--help")

	cmd.Run()
}
