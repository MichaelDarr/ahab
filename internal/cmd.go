package internal

import (
	"os"
	"os/exec"
)

// Docker runs a Docker command
func Docker(arg ...string) (err error) {
	cmd := exec.Command("docker", arg...)
	err = usrCmd(cmd)
	return
}

// usrCmd attatches a command to the user's terminal and runs it
func usrCmd(cmd *exec.Cmd) (err error) {
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err = cmd.Run()
	return
}
