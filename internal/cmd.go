package internal

import (
	"os"
	"os/exec"
)

// Docker runs a Docker command
func Docker(opts *[]string) (err error) {
	cmd := exec.Command("docker", *opts...)
	err = usrCmd(cmd)
	return
}

// usrCmd attatches a command to the user's terminal, prints, and runs it
func usrCmd(cmd *exec.Cmd) (err error) {
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	PrintCmd(cmd)
	err = cmd.Run()
	return
}
